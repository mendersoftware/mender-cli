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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

const (
	argTokenShowRaw  = "raw"
	argTokenShowJSON = "json"
)

// jwtDateClaims lists the standard JWT registered claim names that hold a
// NumericDate (seconds since the Unix epoch). In the default decoded view,
// each of these is rendered as an RFC3339 timestamp with a relative-time
// hint.
var jwtDateClaims = map[string]bool{
	"exp":       true,
	"iat":       true,
	"nbf":       true,
	"auth_time": true,
}

// jwtClaimLabels maps well-known JWT / OIDC / mender claim names to friendly
// labels used by the default decoded view. Unknown keys are rendered with
// their raw name.
var jwtClaimLabels = map[string]string{
	// RFC 7519 registered claims
	"iss": "Issuer",
	"sub": "Subject",
	"aud": "Audience",
	"exp": "Expiration",
	"nbf": "Not Before",
	"iat": "Issued At",
	"jti": "JWT ID",
	// Common OIDC / OAuth2 claims
	"auth_time": "Auth Time",
	"azp":       "Authorized Party",
	"scp":       "Scope",
	"scope":     "Scope",
	"email":     "Email",
	"name":      "Name",
	// JOSE header
	"alg": "Algorithm",
	"typ": "Type",
	"kid": "Key ID",
	"cty": "Content Type",
	// mender-specific
	"mender.tenant": "Tenant",
	"mender.plan":   "Plan",
	"mender.trial":  "Trial",
	"mender.user":   "User",
	"mender.addons": "Add-ons",
}

var tokenShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the locally stored token (JWT header + payload, decoded).",
	Long: "Print the locally stored authentication token. By default the JWT " +
		"header and payload are rendered as a human-readable table with " +
		"friendly claim labels and RFC3339 timestamps for the standard date " +
		"claims (exp, iat, nbf, auth_time); the signature segment is " +
		"intentionally omitted. Use --json to emit the decoded header and " +
		"payload as JSON. Use --raw to print the token verbatim (useful for " +
		"piping to clipboard tools or Authorization headers).",
	Run: func(c *cobra.Command, args []string) {
		raw, err := c.Flags().GetBool(argTokenShowRaw)
		CheckErr(err)
		asJSON, err := c.Flags().GetBool(argTokenShowJSON)
		CheckErr(err)
		if raw && asJSON {
			CheckErr(fmt.Errorf("--raw and --json are mutually exclusive"))
		}
		CheckErr(runTokenShow(c, raw, asJSON, os.Stdout))
	},
}

func init() {
	tokenShowCmd.Flags().Bool(argTokenShowRaw, false,
		"print the unparsed token verbatim instead of the decoded JWT")
	tokenShowCmd.Flags().Bool(argTokenShowJSON, false,
		"print the decoded JWT header and payload as JSON")
}

func runTokenShow(cmd *cobra.Command, raw, asJSON bool, out io.Writer) error {
	path, err := resolveTokenPath(cmd)
	if err != nil {
		return err
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(
				"no token stored; run `mender-cli login` or " +
					"`mender-cli token set` first")
		}
		return fmt.Errorf("failed to read token: %w", err)
	}
	token := strings.TrimSpace(string(contents))

	if raw {
		fmt.Fprintln(out, token)
		return nil
	}

	header, payload, err := decodeJWT(token)
	if err != nil {
		return fmt.Errorf("stored token is not a JWT: %w "+
			"(use --raw to print it verbatim)", err)
	}

	if asJSON {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]any{
			"header":  header,
			"payload": payload,
		})
	}

	return renderDecodedJWT(out, header, payload)
}

// renderDecodedJWT prints the JWT header and payload as a friendly aligned
// table: one section per part, friendly labels for known claim names, RFC3339
// timestamps with a relative-time hint for standard date claims, compact JSON
// for nested arrays / objects.
func renderDecodedJWT(out io.Writer, header, payload map[string]any) error {
	now := time.Now()
	if _, err := fmt.Fprintln(out, "Header:"); err != nil {
		return err
	}
	if err := writeClaimTable(out, header, now); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "\nPayload:"); err != nil {
		return err
	}
	return writeClaimTable(out, payload, now)
}

func writeClaimTable(out io.Writer, claims map[string]any, now time.Time) error {
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return claimLabel(keys[i]) < claimLabel(keys[j])
	})
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		label := claimLabel(k)
		value := formatClaimValue(k, claims[k], now)
		if _, err := fmt.Fprintf(tw, "  %s:\t%s\n", label, value); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func claimLabel(key string) string {
	if l, ok := jwtClaimLabels[key]; ok {
		return l
	}
	return key
}

// formatClaimValue renders a single claim value. Known date claims are shown
// as RFC3339 UTC with a relative-time hint. Strings / bools / numbers are
// rendered verbatim. Arrays and objects fall back to compact JSON so they
// stay on one line.
func formatClaimValue(key string, v any, now time.Time) string {
	if jwtDateClaims[key] {
		if sec, ok := toUnixSeconds(v); ok {
			return formatTimestamp(sec, now)
		}
	}
	switch x := v.(type) {
	case string:
		return x
	case bool:
		return fmt.Sprintf("%t", x)
	case float64:
		// JSON numbers come in as float64; print without a trailing .0 when
		// the value is integral.
		if x == float64(int64(x)) {
			return fmt.Sprintf("%d", int64(x))
		}
		return fmt.Sprintf("%g", x)
	case json.Number:
		return x.String()
	case nil:
		return "null"
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

func formatTimestamp(sec int64, now time.Time) string {
	t := time.Unix(sec, 0).UTC()
	diff := t.Sub(now)
	rel := "in " + humanDuration(diff)
	if diff < 0 {
		rel = humanDuration(-diff) + " ago"
	}
	return fmt.Sprintf("%s (%s)", t.Format(time.RFC3339), rel)
}

func toUnixSeconds(v any) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

// resolveTokenPath returns the token path that subcommands should read/write,
// honoring the persistent --token override or falling back to the default
// platform-specific location.
func resolveTokenPath(cmd *cobra.Command) (string, error) {
	tokenPath, err := cmd.Flags().GetString(argRootToken)
	if err != nil {
		return "", err
	}
	if tokenPath != "" {
		return tokenPath, nil
	}
	return getDefaultAuthTokenPath()
}
