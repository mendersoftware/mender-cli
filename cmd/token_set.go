// Copyright 2026 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package cmd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/useradm"
	"github.com/mendersoftware/mender-cli/log"
)

var tokenSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Persist an authentication token (e.g. a Personal Access Token).",
	Long: "Persist an authentication token in the local store used by all other " +
		"mender-cli commands. The token is read from stdin when piped, or via a " +
		"masked interactive prompt otherwise. After saving, the token is " +
		"validated against the configured server; a failed validation only " +
		"produces a warning — the token is always saved.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewSetTokenCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

// SetTokenCmd implements `mender-cli token set`.
type SetTokenCmd struct {
	server     string
	skipVerify bool
	tokenPath  string
	stdin      io.Reader
	stdinIsTTY bool
	prompt     func() ([]byte, error)
	// verifier verifies the token against the configured server. Injectable
	// for testing. When nil, a real useradm.Client is used.
	verifier func(server string, skipVerify bool, token string) error
}

func NewSetTokenCmd(cmd *cobra.Command, _ []string) (*SetTokenCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("no server configured")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	tokenPath, err := cmd.Flags().GetString(argRootToken)
	if err != nil {
		return nil, err
	}
	if tokenPath == "" {
		tokenPath, err = getDefaultAuthTokenPath()
		if err != nil {
			return nil, err
		}
	}

	return &SetTokenCmd{
		server:     server,
		skipVerify: skipVerify,
		tokenPath:  tokenPath,
		stdin:      os.Stdin,
		stdinIsTTY: isStdinTTY(),
		prompt:     gopass.GetPasswdMasked,
	}, nil
}

func (c *SetTokenCmd) Run() error {
	token, err := c.readToken()
	if err != nil {
		return err
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("empty token; refusing to save")
	}

	if err := writeAuthToken(c.tokenPath, []byte(token)); err != nil {
		return err
	}
	log.Verb("saved token to: " + c.tokenPath)
	log.Info("token saved")

	c.verify(token)
	return nil
}

func (c *SetTokenCmd) readToken() (string, error) {
	if !c.stdinIsTTY {
		b, err := io.ReadAll(c.stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read token from stdin: %w", err)
		}
		return string(b), nil
	}
	fmt.Fprint(os.Stderr, "Personal Access Token: ")
	b, err := c.prompt()
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}
	return string(b), nil
}

func (c *SetTokenCmd) verify(token string) {
	verify := c.verifier
	if verify == nil {
		verify = func(server string, skipVerify bool, token string) error {
			return useradm.NewClient(server, skipVerify).Verify(token)
		}
	}
	if err := verify(c.server, c.skipVerify, token); err != nil {
		var ve *useradm.VerifyError
		if errors.As(err, &ve) {
			log.Errf("warning: server rejected the token (%s); saved anyway", err)
		} else {
			log.Errf("warning: could not verify token against %s (%s); saved anyway",
				c.server, err)
		}
		return
	}
	log.Infof("token verified against %s", c.server)

	// On success, surface the expiry from the JWT exp claim if present.
	_, payload, err := decodeJWT(token)
	if err != nil {
		// Opaque (non-JWT) PATs are valid; nothing to surface.
		log.Verbf("token is not a JWT, skipping expiry parsing: %v", err)
		return
	}
	exp, ok := jwtExpiry(payload)
	if !ok {
		return
	}
	now := time.Now()
	if exp.Before(now) {
		log.Errf("warning: token expired %s ago (at %s)",
			humanDuration(now.Sub(exp)), exp.Format(time.RFC3339))
		return
	}
	log.Infof("token expires in %s (at %s)",
		humanDuration(exp.Sub(now)), exp.Format(time.RFC3339))
}

// isStdinTTY reports whether os.Stdin is connected to a character device
// (terminal). Returns true if stat fails, treating the unknown case as
// "interactive" to avoid silently consuming pipe data that isn't there.
func isStdinTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return true
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// humanDuration formats d as a coarse human-readable string ("27 days",
// "3 hours", "45 minutes"). Sub-minute durations round to "less than a minute".
func humanDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}
	const day = 24 * time.Hour
	switch {
	case d >= day:
		days := int(math.Round(d.Hours() / 24))
		return pluralize(days, "day")
	case d >= time.Hour:
		hours := int(math.Round(d.Hours()))
		return pluralize(hours, "hour")
	default:
		mins := int(math.Round(d.Minutes()))
		return pluralize(mins, "minute")
	}
}

func pluralize(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
