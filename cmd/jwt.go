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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// decodeJWT parses a compact-serialized JWT and returns its decoded header and
// payload as generic maps. The signature segment is intentionally not returned
// or validated — mender-cli is an API client, not the auth server.
func decodeJWT(token string) (header, payload map[string]any, err error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf(
			"not a JWT: expected 3 dot-separated segments, got %d", len(parts))
	}

	header, err = decodeJWTSegment(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid JWT header: %w", err)
	}
	payload, err = decodeJWTSegment(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid JWT payload: %w", err)
	}
	return header, payload, nil
}

func decodeJWTSegment(seg string) (map[string]any, error) {
	raw, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		// Some encoders include padding; tolerate that.
		raw, err = base64.URLEncoding.DecodeString(seg)
		if err != nil {
			return nil, fmt.Errorf("base64 decode failed: %w", err)
		}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}
	return out, nil
}

// jwtExpiry returns the value of the "exp" claim as a time, plus true if a
// numeric exp claim was present.
func jwtExpiry(payload map[string]any) (time.Time, bool) {
	v, ok := payload["exp"]
	if !ok {
		return time.Time{}, false
	}
	switch n := v.(type) {
	case float64:
		return time.Unix(int64(n), 0), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return time.Time{}, false
		}
		return time.Unix(i, 0), true
	default:
		return time.Time{}, false
	}
}
