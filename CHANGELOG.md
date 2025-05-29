---
## 1.12.0-build4.1 - 2025-05-29


### Bug Fixes


- *(tests)* Use mongosh instead of mongo command
 ([03ade99](https://github.com/mendersoftware/mender-cli/commit/03ade99fa8c16d7e1e35a219d6f18c3921aef652))  by @kjaskiewiczz


  Legacy "mongo" shell has been removed in MongoDB 6.0

- Check file format when uploading artifacts
([MEN-7860](https://northerntech.atlassian.net/browse/MEN-7860)) ([238e860](https://github.com/mendersoftware/mender-cli/commit/238e8608c2cbebc0aa8733dd73dc016feda03356))  by @bahaa-ghazal




### Features


- Add support for pagination to `devices list` command
([MEN-7794](https://northerntech.atlassian.net/browse/MEN-7794)) ([7162a7c](https://github.com/mendersoftware/mender-cli/commit/7162a7cb606903ede6e84a4ef1b7f15946eaabb6))  by @alfrunes


  Two new flags `--per-page` and `--page` is added to list devices beyond
  the first page of results.
- Use consistent structured logging (slog) to stderr
([MEN-8304](https://northerntech.atlassian.net/browse/MEN-8304)) ([a1c75d9](https://github.com/mendersoftware/mender-cli/commit/a1c75d9b3123c000f61841abe8df47b77021c98b))  by @alfrunes
- Use paginated endpoint to list artifacts
([MEN-8302](https://northerntech.atlassian.net/browse/MEN-8302)) ([3a59a66](https://github.com/mendersoftware/mender-cli/commit/3a59a6679fda21f4b141eb5e24e83da6b972863a))  by @danielskinstad
  - **BREAKING**: `artifacts list` no longer returns all artifacts by
default. It now uses a paginated API endpoint, so only one page of results
is shown at a time, and it can be modified with `--page` and `--per-page`


  Use the new paginated API endpoint to list artifacts and add support
  for pagination to the `artifacts  list` command.




### Security


- Bump golang.org/x/crypto from 0.16.0 to 0.17.0
 ([9bb6ed7](https://github.com/mendersoftware/mender-cli/commit/9bb6ed7149f571989ec3d009f8aff8e345d56e04))  by @dependabot[bot]


  Bumps [golang.org/x/crypto](https://github.com/golang/crypto) from 0.16.0 to 0.17.0.
  - [Commits](https://github.com/golang/crypto/compare/v0.16.0...v0.17.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/crypto
    dependency-type: indirect
  ...
- Bump the golang-dependencies group with 2 updates
 ([780b00b](https://github.com/mendersoftware/mender-cli/commit/780b00b1aaf29e64bbd0d3e18428f756d683a24f))  by @dependabot[bot]


  Bumps the golang-dependencies group with 2 updates: [github.com/google/uuid](https://github.com/google/uuid) and [github.com/spf13/viper](https://github.com/spf13/viper).
  
  
  Updates `github.com/google/uuid` from 1.4.0 to 1.5.0
  - [Release notes](https://github.com/google/uuid/releases)
  - [Changelog](https://github.com/google/uuid/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/google/uuid/compare/v1.4.0...v1.5.0)
  
  Updates `github.com/spf13/viper` from 1.18.0 to 1.18.2
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.18.0...v1.18.2)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/google/uuid
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-patch
    dependency-group: golang-dependencies
  ...
- Bump the golang-dependencies group with 4 updates
 ([a07412e](https://github.com/mendersoftware/mender-cli/commit/a07412ec27749047c6f818b3874350f01f914ce9))  by @dependabot[bot]


  Bumps the golang-dependencies group with 4 updates: [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb), [github.com/google/uuid](https://github.com/google/uuid), [golang.org/x/sys](https://github.com/golang/sys) and [golang.org/x/term](https://github.com/golang/term).
  
  
  Updates `github.com/cheggaaa/pb/v3` from 3.1.4 to 3.1.5
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.1.4...v3.1.5)
  
  Updates `github.com/google/uuid` from 1.5.0 to 1.6.0
  - [Release notes](https://github.com/google/uuid/releases)
  - [Changelog](https://github.com/google/uuid/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/google/uuid/compare/v1.5.0...v1.6.0)
  
  Updates `golang.org/x/sys` from 0.15.0 to 0.16.0
  - [Commits](https://github.com/golang/sys/compare/v0.15.0...v0.16.0)
  
  Updates `golang.org/x/term` from 0.15.0 to 0.16.0
  - [Commits](https://github.com/golang/term/compare/v0.15.0...v0.16.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/cheggaaa/pb/v3
    dependency-type: direct:production
    update-type: version-update:semver-patch
    dependency-group: golang-dependencies
  - dependency-name: github.com/google/uuid
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/sys
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...
- Bump google.golang.org/protobuf from 1.31.0 to 1.33.0
 ([0c648c7](https://github.com/mendersoftware/mender-cli/commit/0c648c71dc0649d653a29c8bea542ab4d8a41e0a))  by @dependabot[bot]


  Bumps google.golang.org/protobuf from 1.31.0 to 1.33.0.
  
  ---
  updated-dependencies:
  - dependency-name: google.golang.org/protobuf
    dependency-type: indirect
  ...
- Bump golang.org/x/net from 0.19.0 to 0.23.0
 ([a80befd](https://github.com/mendersoftware/mender-cli/commit/a80befdc4f927bcf52050a488141d96fe2c371eb))  by @dependabot[bot]


  Bumps [golang.org/x/net](https://github.com/golang/net) from 0.19.0 to 0.23.0.
  - [Commits](https://github.com/golang/net/compare/v0.19.0...v0.23.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/net
    dependency-type: indirect
  ...
- Bump the golang-dependencies group with 5 updates
 ([1d437a0](https://github.com/mendersoftware/mender-cli/commit/1d437a053cb5de9ccda7c6778ec664007a7a75a5))  by @dependabot[bot]


  Bumps the golang-dependencies group with 5 updates:
  
  | Package | From | To |
  | --- | --- | --- |
  | [github.com/gorilla/websocket](https://github.com/gorilla/websocket) | `1.5.1` | `1.5.3` |
  | [github.com/spf13/cobra](https://github.com/spf13/cobra) | `1.8.0` | `1.8.1` |
  | [github.com/spf13/viper](https://github.com/spf13/viper) | `1.18.2` | `1.19.0` |
  | [golang.org/x/sys](https://github.com/golang/sys) | `0.20.0` | `0.21.0` |
  | [golang.org/x/term](https://github.com/golang/term) | `0.20.0` | `0.21.0` |
  
  
  Updates `github.com/gorilla/websocket` from 1.5.1 to 1.5.3
  - [Release notes](https://github.com/gorilla/websocket/releases)
  - [Commits](https://github.com/gorilla/websocket/compare/v1.5.1...v1.5.3)
  
  Updates `github.com/spf13/cobra` from 1.8.0 to 1.8.1
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Commits](https://github.com/spf13/cobra/compare/v1.8.0...v1.8.1)
  
  Updates `github.com/spf13/viper` from 1.18.2 to 1.19.0
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.18.2...v1.19.0)
  
  Updates `golang.org/x/sys` from 0.20.0 to 0.21.0
  - [Commits](https://github.com/golang/sys/compare/v0.20.0...v0.21.0)
  
  Updates `golang.org/x/term` from 0.20.0 to 0.21.0
  - [Commits](https://github.com/golang/term/compare/v0.20.0...v0.21.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/gorilla/websocket
    dependency-type: direct:production
    update-type: version-update:semver-patch
    dependency-group: golang-dependencies
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-patch
    dependency-group: golang-dependencies
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/sys
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...
- Bump the golang-dependencies group across 1 directory with 2 updates
 ([dc6a734](https://github.com/mendersoftware/mender-cli/commit/dc6a734447233e74b226bdc6f904a595b73fb83c))  by @dependabot[bot]


  Bumps the golang-dependencies group with 2 updates in the / directory: [golang.org/x/sys](https://github.com/golang/sys) and [golang.org/x/term](https://github.com/golang/term).
  
  
  Updates `golang.org/x/sys` from 0.22.0 to 0.25.0
  - [Commits](https://github.com/golang/sys/compare/v0.22.0...v0.25.0)
  
  Updates `golang.org/x/term` from 0.22.0 to 0.24.0
  - [Commits](https://github.com/golang/term/compare/v0.22.0...v0.24.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/sys
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...
- Bump golang.org/x/crypto from 0.21.0 to 0.31.0
 ([c3511ef](https://github.com/mendersoftware/mender-cli/commit/c3511ef24d5eace32ae9afdf33e4066f7033935a))  by @dependabot[bot]


  Bumps [golang.org/x/crypto](https://github.com/golang/crypto) from 0.21.0 to 0.31.0.
  - [Commits](https://github.com/golang/crypto/compare/v0.21.0...v0.31.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/crypto
    dependency-type: indirect
  ...
- Bump golang.org/x/crypto from 0.31.0 to 0.35.0
 ([81a7613](https://github.com/mendersoftware/mender-cli/commit/81a761393a7fdeb4550e305bbab714241ca44813))  by @dependabot[bot]


  Bumps [golang.org/x/crypto](https://github.com/golang/crypto) from 0.31.0 to 0.35.0.
  - [Commits](https://github.com/golang/crypto/compare/v0.31.0...v0.35.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/crypto
    dependency-version: 0.35.0
    dependency-type: indirect
  ...
- Bump busybox from 1.36.1 to 1.37.0
 ([92c4f52](https://github.com/mendersoftware/mender-cli/commit/92c4f5269c9f7a4e0415da6ec1fea8d20a2b2963))  by @dependabot[bot]


  Bumps busybox from 1.36.1 to 1.37.0.
  
  ---
  updated-dependencies:
  - dependency-name: busybox
    dependency-version: 1.37.0
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump the golang-dependencies group with 3 updates
 ([c33c148](https://github.com/mendersoftware/mender-cli/commit/c33c148ea51b4c192456fa875030ff5a84e615b2))  by @dependabot[bot]


  Bumps the golang-dependencies group with 3 updates: [github.com/spf13/viper](https://github.com/spf13/viper), [golang.org/x/sys](https://github.com/golang/sys) and [golang.org/x/term](https://github.com/golang/term).
  
  
  Updates `github.com/spf13/viper` from 1.19.0 to 1.20.1
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.19.0...v1.20.1)
  
  Updates `golang.org/x/sys` from 0.31.0 to 0.32.0
  - [Commits](https://github.com/golang/sys/compare/v0.31.0...v0.32.0)
  
  Updates `golang.org/x/term` from 0.30.0 to 0.31.0
  - [Commits](https://github.com/golang/term/compare/v0.30.0...v0.31.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-version: 1.20.1
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/sys
    dependency-version: 0.32.0
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/term
    dependency-version: 0.31.0
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...





## 1.12.0-build4 - 2023-12-08


### Bug Fixes


- *(terminal)* Keep stdout clean
 ([80f7d47](https://github.com/mendersoftware/mender-cli/commit/80f7d47b7aacb60968491d8fd4317a1eeda283d0))  by @fellerts


  When connecting to a device through the remote terminal, the CLI prints the
  following on stdout: 'Connecting to the device <uuid>...' This should rather be
  printed to stderr to make scripting easier.

- *(terminal)* Use term's built-in IsTerminal
 ([0c3a24c](https://github.com/mendersoftware/mender-cli/commit/0c3a24c3e4cc9039c8d3a7c76bb484abcc12c2db))  by @fellerts


  Solves the following issue:
  $ mender-cli terminal foo >/dev/null
  FAILURE: Unable to get the terminal size: inappropriate ioctl for device





### Features


- Adding artifact downloads to CLI
 ([7350f39](https://github.com/mendersoftware/mender-cli/commit/7350f391fdd1473f6ac82e84617500ead804bec8))  by @stevenleadbeater


  Added a new command to download artifacts through the API by artifact ID. The client will first lookup the artifact details to get the download link, this contains the download URL.
- Post metadata with the direct upload with skip verify
([MEN-6696](https://northerntech.atlassian.net/browse/MEN-6696)) ([4dd47e9](https://github.com/mendersoftware/mender-cli/commit/4dd47e959f96f78678128aedab46dda617ad5e36))  by @merlin-northern
- Add meaningful error message for conflict artifact
 ([e71aa55](https://github.com/mendersoftware/mender-cli/commit/e71aa55266f0075868eb43f4c7d7e574b5c7efe3))  by @MuchoLucho


  Proposed by the community from
  https://hub.mender.io/t/error-message-not-clear-enough-in-mender-cli-artifacts-upload




### Security


- Bump golang from 1.20.5-alpine3.17 to 1.20.6-alpine3.17
 ([5a08d0e](https://github.com/mendersoftware/mender-cli/commit/5a08d0e83562bbf974f1c75d661d12bbb5d97d8b))  by @dependabot[bot]


  Bumps golang from 1.20.5-alpine3.17 to 1.20.6-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump github.com/cheggaaa/pb/v3 from 3.1.2 to 3.1.4
 ([055b225](https://github.com/mendersoftware/mender-cli/commit/055b2251a409d473a68fd2a8784ad6aceb1c84db))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.1.2 to 3.1.4.
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.1.2...v3.1.4)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/cheggaaa/pb/v3
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang.org/x/term from 0.9.0 to 0.10.0
 ([1ee437c](https://github.com/mendersoftware/mender-cli/commit/1ee437c57905053bd22bd46f7c50427a5b91e84a))  by @dependabot[bot]


  Bumps [golang.org/x/term](https://github.com/golang/term) from 0.9.0 to 0.10.0.
  - [Commits](https://github.com/golang/term/compare/v0.9.0...v0.10.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang from 1.20.6-alpine3.17 to 1.21.0-alpine3.17
 ([6354eb5](https://github.com/mendersoftware/mender-cli/commit/6354eb5ac7699a07742b656b5294cd129f5c45de))  by @dependabot[bot]


  Bumps golang from 1.20.6-alpine3.17 to 1.21.0-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/google/uuid from 1.3.0 to 1.3.1
 ([00fe49c](https://github.com/mendersoftware/mender-cli/commit/00fe49c81c48f493cda998bcb095b171c285f31e))  by @dependabot[bot]


  Bumps [github.com/google/uuid](https://github.com/google/uuid) from 1.3.0 to 1.3.1.
  - [Release notes](https://github.com/google/uuid/releases)
  - [Changelog](https://github.com/google/uuid/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/google/uuid/compare/v1.3.0...v1.3.1)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/google/uuid
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang.org/x/sys from 0.10.0 to 0.12.0
 ([011adfc](https://github.com/mendersoftware/mender-cli/commit/011adfc7aa4c2957d111dd0ee51a117982eb294c))  by @dependabot[bot]


  Bumps [golang.org/x/sys](https://github.com/golang/sys) from 0.10.0 to 0.12.0.
  - [Commits](https://github.com/golang/sys/compare/v0.10.0...v0.12.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/sys
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang.org/x/term from 0.10.0 to 0.12.0
 ([7e459a2](https://github.com/mendersoftware/mender-cli/commit/7e459a27aca21fc1fc05524d002db06c5f2eb284))  by @dependabot[bot]


  Bumps [golang.org/x/term](https://github.com/golang/term) from 0.10.0 to 0.12.0.
  - [Commits](https://github.com/golang/term/compare/v0.10.0...v0.12.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang from 1.21.0-alpine3.17 to 1.21.1-alpine3.17
 ([cb6cd3f](https://github.com/mendersoftware/mender-cli/commit/cb6cd3fd5af8a0219bc09028075dee9528c37a40))  by @dependabot[bot]


  Bumps golang from 1.21.0-alpine3.17 to 1.21.1-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang.org/x/net from 0.10.0 to 0.17.0
 ([c9e830b](https://github.com/mendersoftware/mender-cli/commit/c9e830b0c679959ed41f3279310f26b03627f7b4))  by @dependabot[bot]


  Bumps [golang.org/x/net](https://github.com/golang/net) from 0.10.0 to 0.17.0.
  - [Commits](https://github.com/golang/net/compare/v0.10.0...v0.17.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/net
    dependency-type: indirect
  ...
- Bump golang from 1.21.1-alpine3.17 to 1.21.3-alpine3.17
 ([73437f9](https://github.com/mendersoftware/mender-cli/commit/73437f9bcad49ab9e8bd55543f5f5ed19b2db688))  by @dependabot[bot]


  Bumps golang from 1.21.1-alpine3.17 to 1.21.3-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump github.com/google/uuid from 1.3.1 to 1.4.0
 ([eb6dc6e](https://github.com/mendersoftware/mender-cli/commit/eb6dc6efaba1fc3033d023c09aef81855c24e06a))  by @dependabot[bot]


  Bumps [github.com/google/uuid](https://github.com/google/uuid) from 1.3.1 to 1.4.0.
  - [Release notes](https://github.com/google/uuid/releases)
  - [Changelog](https://github.com/google/uuid/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/google/uuid/compare/v1.3.1...v1.4.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/google/uuid
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/spf13/viper from 1.16.0 to 1.17.0
 ([c8b9686](https://github.com/mendersoftware/mender-cli/commit/c8b96863690e3314076809f6939a7f3d2c691a9b))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.16.0 to 1.17.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.16.0...v1.17.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump the golang-dependencies group with 4 updates
 ([9be26e6](https://github.com/mendersoftware/mender-cli/commit/9be26e64903130d20a616d0c2160609ee5024882))  by @dependabot[bot]


  Bumps the golang-dependencies group with 4 updates: [github.com/gorilla/websocket](https://github.com/gorilla/websocket), [github.com/spf13/cobra](https://github.com/spf13/cobra), [golang.org/x/sys](https://github.com/golang/sys) and [golang.org/x/term](https://github.com/golang/term).
  
  Updates `github.com/gorilla/websocket` from 1.5.0 to 1.5.1
  - [Release notes](https://github.com/gorilla/websocket/releases)
  - [Commits](https://github.com/gorilla/websocket/compare/v1.5.0...v1.5.1)
  
  Updates `github.com/spf13/cobra` from 1.7.0 to 1.8.0
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Commits](https://github.com/spf13/cobra/compare/v1.7.0...v1.8.0)
  
  Updates `golang.org/x/sys` from 0.13.0 to 0.15.0
  - [Commits](https://github.com/golang/sys/compare/v0.13.0...v0.15.0)
  
  Updates `golang.org/x/term` from 0.13.0 to 0.15.0
  - [Commits](https://github.com/golang/term/compare/v0.13.0...v0.15.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/gorilla/websocket
    dependency-type: direct:production
    update-type: version-update:semver-patch
    dependency-group: golang-dependencies
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/sys
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...
- Bump golang from 1.21.3-alpine3.17 to 1.21.4-alpine3.17
 ([4ac5349](https://github.com/mendersoftware/mender-cli/commit/4ac5349eadd835ee260a09f5e8da186a83cca71b))  by @dependabot[bot]


  Bumps golang from 1.21.3-alpine3.17 to 1.21.4-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang from 1.21.4-alpine3.17 to 1.21.5-alpine3.17
 ([1fb94c2](https://github.com/mendersoftware/mender-cli/commit/1fb94c26f7f713e6b089e60439f246bcb7c1e984))  by @dependabot[bot]


  Bumps golang from 1.21.4-alpine3.17 to 1.21.5-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump the golang-dependencies group with 1 update
 ([18f02bb](https://github.com/mendersoftware/mender-cli/commit/18f02bb000ab1fad3eec81fb80807697a703885a))  by @dependabot[bot]


  Bumps the golang-dependencies group with 1 update: [github.com/spf13/viper](https://github.com/spf13/viper).
  
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.17.0...v1.18.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
    dependency-group: golang-dependencies
  ...





## 1.11.0-build9 - 2023-07-04


### Bug Fixes


- Hide errors when help flag is present
([MEN-6357](https://northerntech.atlassian.net/browse/MEN-6357)) ([f989fa1](https://github.com/mendersoftware/mender-cli/commit/f989fa1d19b1ff7963bf9ddc43f171575ddcd970))  by @svenschwermer
- Allow to use --token to specify API token
([MEN-6357](https://northerntech.atlassian.net/browse/MEN-6357)) ([589971b](https://github.com/mendersoftware/mender-cli/commit/589971be7401baee1d2afc38e363151288b785fc))  by @svenschwermer


  Previously, this flag had been ignored for all commands but the login
  command.
- List devices in raw mode not in stdout
 ([bd593f7](https://github.com/mendersoftware/mender-cli/commit/bd593f79181fdc7df715bb802b02fb957e046b6d))  by @estape11


  Printing the list of devices with flag -r was not directed to
  stdout but stderr.




### Features


- Direct upload.
([MEN-6338](https://northerntech.atlassian.net/browse/MEN-6338)) ([eadae02](https://github.com/mendersoftware/mender-cli/commit/eadae02dd62a6c56580eb0be3a46526c1e39e148))  by @merlin-northern




### Security


- Bump github.com/spf13/viper from 1.14.0 to 1.15.0
 ([1a566cc](https://github.com/mendersoftware/mender-cli/commit/1a566cce790ec230c8c533bfc771205149718b7a))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.14.0 to 1.15.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.14.0...v1.15.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang.org/x/net from 0.4.0 to 0.7.0
 ([2be207c](https://github.com/mendersoftware/mender-cli/commit/2be207c6a62c604a365b089744809cf6a6b3a7f2))  by @dependabot[bot]


  Bumps [golang.org/x/net](https://github.com/golang/net) from 0.4.0 to 0.7.0.
  - [Release notes](https://github.com/golang/net/releases)
  - [Commits](https://github.com/golang/net/compare/v0.4.0...v0.7.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/net
    dependency-type: indirect
  ...
- Bump github.com/cheggaaa/pb/v3 from 3.1.0 to 3.1.2
 ([fa685d0](https://github.com/mendersoftware/mender-cli/commit/fa685d0849a0141fa76ae9d5f5a00a0c0b41b00c))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.1.0 to 3.1.2.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.1.0...v3.1.2)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/cheggaaa/pb/v3
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang from 1.20.1-alpine3.17 to 1.20.4-alpine3.17
 ([2ac99e8](https://github.com/mendersoftware/mender-cli/commit/2ac99e8932c385222b7052943bbcfac4adfb2d5a))  by @dependabot[bot]


  Bumps golang from 1.20.1-alpine3.17 to 1.20.4-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump github.com/spf13/cobra from 1.6.1 to 1.7.0
 ([becafe4](https://github.com/mendersoftware/mender-cli/commit/becafe47f0e8975e725621bf350d5316bd08d119))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.6.1 to 1.7.0.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Commits](https://github.com/spf13/cobra/compare/v1.6.1...v1.7.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang.org/x/term from 0.6.0 to 0.9.0
 ([92911c9](https://github.com/mendersoftware/mender-cli/commit/92911c90d0e81db1a4c097bf5bd26c6ce8dc43ac))  by @dependabot[bot]


  Bumps [golang.org/x/term](https://github.com/golang/term) from 0.6.0 to 0.9.0.
  - [Commits](https://github.com/golang/term/compare/v0.6.0...v0.9.0)
  
  ---
  updated-dependencies:
  - dependency-name: golang.org/x/term
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/spf13/viper from 1.15.0 to 1.16.0
 ([248519e](https://github.com/mendersoftware/mender-cli/commit/248519e62ae86eae0f16d6fa765bf9054f20a010))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.15.0 to 1.16.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.15.0...v1.16.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang from 1.20.4-alpine3.17 to 1.20.5-alpine3.17
 ([72cf993](https://github.com/mendersoftware/mender-cli/commit/72cf993087c0444d7ffc83e05a3f78498e3d3b82))  by @dependabot[bot]


  Bumps golang from 1.20.4-alpine3.17 to 1.20.5-alpine3.17.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...





## 1.10.0-build7 - 2022-12-19


### Bug Fixes


- *(portforward)* Allow to get the server from mender-clirc configuration file
 ([a0d6c05](https://github.com/mendersoftware/mender-cli/commit/a0d6c05c7cbfa765e8cb74c6d74ab9522d951db3))  by @clementperon


  At the moment portforward doesn't take into account the mender clirc
  configuration file.
  
  Duplicate what's is done in mender terminal to get the proper configuration.





### Security


- Bump github.com/spf13/cobra from 1.5.0 to 1.6.1
 ([b37d190](https://github.com/mendersoftware/mender-cli/commit/b37d1900d58427bd70d263c5f028a92c4ff23bb9))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.5.0 to 1.6.1.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Commits](https://github.com/spf13/cobra/compare/v1.5.0...v1.6.1)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/spf13/viper from 1.13.0 to 1.14.0
 ([0b273c2](https://github.com/mendersoftware/mender-cli/commit/0b273c2f92c22d056ddc2ef8f3b2ac1e592ae478))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.13.0 to 1.14.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.13.0...v1.14.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...





## 1.9.0-build8 - 2022-09-13


### Bug Fixes


- Stop the port-forward command on errors when reading the websocket
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([6231746](https://github.com/mendersoftware/mender-cli/commit/6231746cc6a211002dbb0c089bc053718c5c230b)) 


  There is no way to recover from failures when reading the websocket, as
  gorilla/websocket will panic with `repeated read on failed websocket
  connection`. Exit the command instead.
- Automatically handle reconnections in port-forward
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([77072de](https://github.com/mendersoftware/mender-cli/commit/77072de43c2898160127d77df8e99e8e7a481ac4)) 


  In case of errors, instead of stopping the process and exiting,
  automatically reconnect to the device and continue port-forwarding the
  connections. This improves the user experience in case of temporary
  errors when port-forwarding an HTTP server, where each request is a new
  connection. Additionally, ignore errors on connection close.
- Use a mutex lock per connection instead of a global one
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([33ed7ae](https://github.com/mendersoftware/mender-cli/commit/33ed7aebeec428cecdd4291810126301ff437f20)) 




### Features


- Add option to specify a JWT or personal access token
([MEN-5660](https://northerntech.atlassian.net/browse/MEN-5660)) ([fa0134b](https://github.com/mendersoftware/mender-cli/commit/fa0134b4a56e0c46e9e7b6313de43713de8c1cf7)) 




### Security


- Bump github.com/cheggaaa/pb/v3 from 3.0.8 to 3.1.0
 ([cd5eadf](https://github.com/mendersoftware/mender-cli/commit/cd5eadfcf7751bfec0adc1248fef9f39ce44bd47))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.8 to 3.1.0.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.8...v3.1.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/cheggaaa/pb/v3
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/spf13/cobra from 1.4.0 to 1.5.0
 ([2571435](https://github.com/mendersoftware/mender-cli/commit/2571435f67c5da459654cb3c0b73896895a0e403))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.4.0 to 1.5.0.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Commits](https://github.com/spf13/cobra/compare/v1.4.0...v1.5.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump github.com/spf13/viper from 1.10.1 to 1.13.0
 ([1ba7e35](https://github.com/mendersoftware/mender-cli/commit/1ba7e358c3dc9abf2c93dedef7066651470015f5))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.10.1 to 1.13.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.10.1...v1.13.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...





## 1.8.0-build9 - 2022-05-02


### Changelog


- All: Bump golang from 1.16.2-alpine3.12 to 1.16.3-alpine3.12
 ([603286e](https://github.com/mendersoftware/mender-cli/commit/603286e5a8036f4b6c608dec4067518cd9329ba9))  by @dependabot[bot]


  Bumps golang from 1.16.2-alpine3.12 to 1.16.3-alpine3.12.
- All: Bump github.com/cheggaaa/pb/v3 from 3.0.7 to 3.0.8
 ([ac5e395](https://github.com/mendersoftware/mender-cli/commit/ac5e395c9b36a0e8f714c3943d923f07d7ee705c))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.7 to 3.0.8.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.7...v3.0.8)
- All: Bump golang from 1.16.3-alpine3.12 to 1.16.4-alpine3.12
 ([ca48cfe](https://github.com/mendersoftware/mender-cli/commit/ca48cfee4fb0d0a1c32471ef61edc3d4dab20984))  by @dependabot[bot]


  Bumps golang from 1.16.3-alpine3.12 to 1.16.4-alpine3.12.
- All: Bump golang from 1.16.4-alpine3.12 to 1.16.5-alpine3.12
 ([a381e47](https://github.com/mendersoftware/mender-cli/commit/a381e47fd3ab5c6ea1e0d2273f6b64de95cd6dbf))  by @dependabot[bot]


  Bumps golang from 1.16.4-alpine3.12 to 1.16.5-alpine3.12.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- All: Bump github.com/spf13/viper from 1.7.1 to 1.8.0
 ([d108b2f](https://github.com/mendersoftware/mender-cli/commit/d108b2f85cfb13c601e87cc4941e3ba350cf7870))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.7.1 to 1.8.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.7.1...v1.8.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/viper from 1.8.0 to 1.8.1
 ([58c1e07](https://github.com/mendersoftware/mender-cli/commit/58c1e07f71659bcd6d91f80d0cc85541e83edc37))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.8.0 to 1.8.1.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.8.0...v1.8.1)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- All: Bump github.com/spf13/cobra from 1.1.3 to 1.2.1
 ([4811012](https://github.com/mendersoftware/mender-cli/commit/4811012fcac95f2d453a27a83205d33ece838e56))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.1.3 to 1.2.1.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/spf13/cobra/compare/v1.1.3...v1.2.1)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/google/uuid from 1.2.0 to 1.3.0
 ([e798ef3](https://github.com/mendersoftware/mender-cli/commit/e798ef30dd58f8aeab8f1511f95367398f9481e1))  by @dependabot[bot]


  Bumps [github.com/google/uuid](https://github.com/google/uuid) from 1.2.0 to 1.3.0.
  - [Release notes](https://github.com/google/uuid/releases)
  - [Commits](https://github.com/google/uuid/compare/v1.2.0...v1.3.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/google/uuid
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/viper from 1.8.1 to 1.9.0
 ([77774a7](https://github.com/mendersoftware/mender-cli/commit/77774a7da52af535d765cbb217fead9feaa63085))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.8.1 to 1.9.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.8.1...v1.9.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/viper from 1.9.0 to 1.10.0
 ([825eccc](https://github.com/mendersoftware/mender-cli/commit/825eccc4aeaac90c28af25ca261be19c9e700435))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.9.0 to 1.10.0.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.9.0...v1.10.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/cobra from 1.2.1 to 1.3.0
 ([a64e63d](https://github.com/mendersoftware/mender-cli/commit/a64e63d11336e6a4979e248e15d242c9b8cf9618))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.2.1 to 1.3.0.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/spf13/cobra/compare/v1.2.1...v1.3.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/viper from 1.10.0 to 1.10.1
 ([d4966ff](https://github.com/mendersoftware/mender-cli/commit/d4966ff77d536c6f5cd3ec9e255dc9b5083adb5e))  by @dependabot[bot]


  Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.10.0 to 1.10.1.
  - [Release notes](https://github.com/spf13/viper/releases)
  - [Commits](https://github.com/spf13/viper/compare/v1.10.0...v1.10.1)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/viper
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- All: Bump github.com/gorilla/websocket from 1.4.2 to 1.5.0
 ([2afb517](https://github.com/mendersoftware/mender-cli/commit/2afb5178d7696b6d9f023544c5da273274d1c802))  by @dependabot[bot]


  Bumps [github.com/gorilla/websocket](https://github.com/gorilla/websocket) from 1.4.2 to 1.5.0.
  - [Release notes](https://github.com/gorilla/websocket/releases)
  - [Commits](https://github.com/gorilla/websocket/compare/v1.4.2...v1.5.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/gorilla/websocket
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- All: Bump github.com/spf13/cobra from 1.3.0 to 1.4.0
 ([3cd4d4a](https://github.com/mendersoftware/mender-cli/commit/3cd4d4a565c6c077eb450f9a8fbe810837e6baf8))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.3.0 to 1.4.0.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/spf13/cobra/compare/v1.3.0...v1.4.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/spf13/cobra
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...





## 1.7.0-build8 - 2021-03-23


### Changelog


- All: Bump golang from 1.15.6-alpine3.12 to 1.15.8-alpine3.12
 ([9b8fc05](https://github.com/mendersoftware/mender-cli/commit/9b8fc05a33174e1890031841c4b6c74cea016831))  by @dependabot[bot]


  Bumps golang from 1.15.6-alpine3.12 to 1.15.8-alpine3.12.
- All: Bump github.com/cheggaaa/pb/v3 from 3.0.5 to 3.0.6
 ([519b70e](https://github.com/mendersoftware/mender-cli/commit/519b70e10ab32dbc37e2c6a830d2407f404d364d))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.5 to 3.0.6.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.5...v3.0.6)
- All: Bump github.com/spf13/cobra from 1.1.1 to 1.1.3
 ([4672d09](https://github.com/mendersoftware/mender-cli/commit/4672d093b506d0443a4f22ba3443f879e2b43d3f))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.1.1 to 1.1.3.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/spf13/cobra/compare/v1.1.1...v1.1.3)
- All: Bump golang from 1.15.8-alpine3.12 to 1.16.0-alpine3.12
 ([3e960f3](https://github.com/mendersoftware/mender-cli/commit/3e960f35b9274eabbb4621c7aba707f3e67b1a35))  by @dependabot[bot]


  Bumps golang from 1.15.8-alpine3.12 to 1.16.0-alpine3.12.
- All: Bump golang from 1.16.0-alpine3.12 to 1.16.2-alpine3.12
 ([f2e4d5b](https://github.com/mendersoftware/mender-cli/commit/f2e4d5ba37ba7569d3648b56a29ad3a15d83e6b2))  by @dependabot[bot]


  Bumps golang from 1.16.0-alpine3.12 to 1.16.2-alpine3.12.
- All: Bump github.com/cheggaaa/pb/v3 from 3.0.6 to 3.0.7
 ([a58e5c4](https://github.com/mendersoftware/mender-cli/commit/a58e5c4f294ed8d5c1d106f7200c100eb803d63d))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.6 to 3.0.7.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.6...v3.0.7)




### Fix


- Respect the --server flag from config everywhere
 ([6a814cd](https://github.com/mendersoftware/mender-cli/commit/6a814cd889b037f04c250b332d8360f53c7739d2))  by @oleorhagen


  Previously, the configuration --server flag was only bound to the configuration
  file value in the login command.
  
  By moving the viper configuration to the root command, and fetching the value
  from viper everywhere, the flag is now properly handled everywhere.




### MEN-4318


- Add --record and --playback to terminal command
([MEN-4318](https://northerntech.atlassian.net/browse/MEN-4318)) ([9b6d546](https://github.com/mendersoftware/mender-cli/commit/9b6d546561ab7d535c90e5722887e8f6600a44db))  by @lluiscampos


  Save first a header with some meta-data, and then encode the stream of
  bytes from stdout. Only stdout is recorded.
  
  Made the positional argument for `terminal` command optional, as now for
  playback we don't need to connect to any device.




### MEN-4323


- Add filetransfer upload and download support
([MEN-4323](https://northerntech.atlassian.net/browse/MEN-4323)) ([6f19310](https://github.com/mendersoftware/mender-cli/commit/6f19310418b3633fde3077ad40d2d18d19e40e20))  by @oleorhagen




### Makefile


- Fix port sed'ing
 ([25740ad](https://github.com/mendersoftware/mender-cli/commit/25740ad801674ce4de40c7a11c49a841cffcc5b6))  by @mchalski


  must remove also the 80:80 mapping.
  otherwise leaves an invalid yaml in the .testing compose file:
  
  ```
  mender-api-gateway:
     - 80:80
  ```
  
  this is the direct cause of recent pipeline problems





## 1.6.0-build2 - 2021-01-14


### Changelog


- All: Bump golang from 1.14-alpine3.12 to 1.15.1-alpine3.12
 ([4d35ed3](https://github.com/mendersoftware/mender-cli/commit/4d35ed31f4bfc10377b7ad1eb42bbf4dd6ba74fb))  by @dependabot[bot]


  Bumps golang from 1.14-alpine3.12 to 1.15.1-alpine3.12.
- All: Bump golang from 1.15.1-alpine3.12 to 1.15.2-alpine3.12
 ([277851f](https://github.com/mendersoftware/mender-cli/commit/277851f85e512e1bb2de5ae3e57c6002ed78d588))  by @dependabot[bot]


  Bumps golang from 1.15.1-alpine3.12 to 1.15.2-alpine3.12.
- All: Bump github.com/cheggaaa/pb/v3 from 3.0.4 to 3.0.5
 ([8a30c47](https://github.com/mendersoftware/mender-cli/commit/8a30c475e5cca24941899c4921ca79cb5a6ab326))  by @dependabot[bot]


  Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.4 to 3.0.5.
  - [Release notes](https://github.com/cheggaaa/pb/releases)
  - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.4...v3.0.5)
- All: Bump golang from 1.15.2-alpine3.12 to 1.15.3-alpine3.12
 ([4e601e7](https://github.com/mendersoftware/mender-cli/commit/4e601e7a49b79034eda854aaa0c677f2da7fd7b9))  by @dependabot[bot]


  Bumps golang from 1.15.2-alpine3.12 to 1.15.3-alpine3.12.
- All: Bump github.com/spf13/cobra from 1.0.0 to 1.1.1
 ([a5a1640](https://github.com/mendersoftware/mender-cli/commit/a5a164057120603c0ee1c420f96d27a7b4c1c073))  by @dependabot[bot]


  Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.0.0 to 1.1.1.
  - [Release notes](https://github.com/spf13/cobra/releases)
  - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
  - [Commits](https://github.com/spf13/cobra/compare/v1.0.0...v1.1.1)
- All: Bump golang from 1.15.3-alpine3.12 to 1.15.4-alpine3.12
 ([fcc21a6](https://github.com/mendersoftware/mender-cli/commit/fcc21a69783cbd0ac7488f6a3ae684944d8a4294))  by @dependabot[bot]


  Bumps golang from 1.15.3-alpine3.12 to 1.15.4-alpine3.12.
- All: Bump golang from 1.15.4-alpine3.12 to 1.15.5-alpine3.12
 ([d29b8af](https://github.com/mendersoftware/mender-cli/commit/d29b8af810358846a36d12ae1f7648b65f35720c))  by @dependabot[bot]


  Bumps golang from 1.15.4-alpine3.12 to 1.15.5-alpine3.12.
- All: Bump golang from 1.15.5-alpine3.12 to 1.15.6-alpine3.12
 ([f4f209f](https://github.com/mendersoftware/mender-cli/commit/f4f209ff135d8ed5c9e7186e36cdbf409174d1c7))  by @dependabot[bot]


  Bumps golang from 1.15.5-alpine3.12 to 1.15.6-alpine3.12.




### Fix


- Gocyclo refactoring and switch to go mod
 ([e1c0e07](https://github.com/mendersoftware/mender-cli/commit/e1c0e072e5b9231310f4b9fc7617ab8a4f944267))  by @tranchitella




### MEN-4305


- Send shell pong messages in response to shell ping messages
([MEN-4305](https://northerntech.atlassian.net/browse/MEN-4305)) ([662cd74](https://github.com/mendersoftware/mender-cli/commit/662cd748dccb3c2983dbce248e9e16374a2b3727)) 




### QA-238


- Check_commits: move from unit tests to template
([QA-238](https://northerntech.atlassian.net/browse/QA-238)) ([0ef51e9](https://github.com/mendersoftware/mender-cli/commit/0ef51e944e1144b610c2aa5a09a6f724a0d18d12))  by @tranchitella




### Gitlab


- Cleanup test:format
 ([0247d69](https://github.com/mendersoftware/mender-cli/commit/0247d691c28c4eba44e700471f74601c1dae08a5))  by @mchalski


  it mixes up tests with static checks againg - untangle them.
  
  becomes our usual test:static, and uses new targets to get
  deps and run.
- Move fast tests to front
 ([c189281](https://github.com/mendersoftware/mender-cli/commit/c189281303440e06156e486b74eaf741d6c98f03))  by @mchalski


  this is not the case with standard mendertesting templates,
  but since we're customizing this one - let's make our life easier.
- Cleanup acceptance test stages
 ([d2af13f](https://github.com/mendersoftware/mender-cli/commit/d2af13f4e0c8aec45922b71fc28ebc85e4835d75))  by @mchalski


  first off, they aren't 'fast'. rename stages and jobs
  to reflect what they actually do.
- Cleanup test:check
 ([1912569](https://github.com/mendersoftware/mender-cli/commit/191256942275365646ec5de22f9a8aad9b34e3f5))  by @mchalski


  the templated unittest stage does that and more.
- Move the global before_script
 ([3fb204d](https://github.com/mendersoftware/mender-cli/commit/3fb204d15654927a5460f6f3454101dc887e280b))  by @mchalski


  it's specific to the compile phase.
- Cosmetic name changes
 ([5d43fe5](https://github.com/mendersoftware/mender-cli/commit/5d43fe581959277a7d8576dc0c2b7f7a334ecd4c))  by @mchalski


  correspond closer to the Makefile.





## 1.5.0-build4 - 2020-08-27


### Add


- Bash auto-completion functionality
 ([0c87a96](https://github.com/mendersoftware/mender-cli/commit/0c87a96f32f58daf3e1014251117e537b719a25b))  by @oleorhagen


  This adds the possibility to have mender-cli auto-complete commands in Bash,
  just like any other mature CLI tool. The change is little, and for someone using
  this tool a fair bit, auto-completion should help make the tool more accessible,
  and self-discoverable, as the functionality is automatically shown on [Tab].
- Zsh auto-completion support
 ([0ae72a6](https://github.com/mendersoftware/mender-cli/commit/0ae72a60e685394cbede7c1755e42a530c980d3c))  by @oleorhagen


  This mirrors the support for auto-completion in Bash
- Make the server flag default to hosted Mender
 ([8d9f652](https://github.com/mendersoftware/mender-cli/commit/8d9f6520082e3fa04100341215cb365a9f537b0e))  by @oleorhagen


  This makes more sense I think, as we generally want users to default to hosted.
  Having to manually enter it every time is also annoying.
- Configuration file functionality
 ([ca1078c](https://github.com/mendersoftware/mender-cli/commit/ca1078c0d7b6bc0d810d50cf253e070c092c921d))  by @oleorhagen


  This adds the possibility to add the username and password to a configuration
  file, in which the 'mender-cli' tool will look if no password or username is set
  on the CLI. The configuration file is expected to be JSON.
  
  The configuration file can be located in one of:
  
  * /etc/mender-cli
  * $HOME
  * . (directory where binary is run from)
  
  and must be named like:
  
  ```console
  .mender-clirc.json
  ```
  
  This helps usage, in that now, in order to login, a user with a configuration
  file can do:
  
  ```console
  $ mender-cli login
  ```
  
  as opposed to:
  ```console
  $ mender-cli --username foo --password bar --server bar.com
  ```
  
  The parameters which are configurable from the config file are:
  
  * username
  * password
  * server





## 1.2.0b1-build4 - 2019-04-10


### Travis


- Bump go to 1.10
 ([d0a362a](https://github.com/mendersoftware/mender-cli/commit/d0a362ad48ad24f11cb7c955121d8ee14d76da12))  by @mchalski


  quotes are manadatory:
  https://github.com/travis-ci/travis-ci/issues/9247




### Vendor


- Update
 ([25f3ce0](https://github.com/mendersoftware/mender-cli/commit/25f3ce097dc940ddb34d48dec9489d3d7c66a61f))  by @mchalski





## 1.1.0b1 - 2018-10-09


### Acceptance


- Travis and docker/compose setup
 ([7243c18](https://github.com/mendersoftware/mender-cli/commit/7243c187ebf57e8d041611d458088dda05e686d9)) 


  - Dockerfile for acceptance tests (also contains the cli binary via mount)
  - test entrypoint
  - single-tenant acceptance tests compose
  
  note that the whole 'integration' repo is mounted - the testing container
  will need access to compose files for running commands.




### Cmd


- Fix 'login' error message typo
 ([17afe14](https://github.com/mendersoftware/mender-cli/commit/17afe14c29f0f70f4b687306b629b753c1c96a3f)) 




### Travis


- Build stages + backend integration tests
 ([78d9efe](https://github.com/mendersoftware/mender-cli/commit/78d9efe166792b272761c201bd816376d021baeb)) 


  actually, this is mostly about build stages.
  
  the refactor is not pretty because build stages vs env matrix seems
  like a half-done feature. most notably - defining a global matrix
  with os/arch combinations doesn't help - jobs inherit only the first
  entry in the matrix. it's also not possible to define matrices inside jobs.
  
  the only way is to copy and paste a job N times, with different env vars.
  for most jobs this is ok - esp. with a helper bash script. unfortunately
  with deploy stage, one has to repeat large blocks of yml, as it's impossible
  to push provider defs into a helper.
  
  there is a placeholder for backend testing jobs but the cli is not
  yet a part of the integration setup (no dockerfile, compose, etc.).
  first decide how we want to include/test it, then actually add a
  section here.





## 1.0.0b1 - 2018-04-20


### Client


- Deployments client with artifact upload method
 ([1af2420](https://github.com/mendersoftware/mender-cli/commit/1af242065d7f3dab47469ddaa08e9b20cdaaa8bc))  by @kjaskiewiczz




### Client/useradm


- Useradm client supporting /login
 ([5182401](https://github.com/mendersoftware/mender-cli/commit/5182401d3d8e2feb157fb7fc3f2c93dd6a446c9c)) 
- Introduce verbose logging
 ([30127d4](https://github.com/mendersoftware/mender-cli/commit/30127d4e9ac1cf962dd40bf38b2f5639fdfeafbd)) 




### Cmd


- Error handling on root
 ([688816d](https://github.com/mendersoftware/mender-cli/commit/688816d3bc3d180c862b651ec5a68d05eb97eeec)) 
- Store token in user's homedir
 ([598235a](https://github.com/mendersoftware/mender-cli/commit/598235a89de5efca29f2a5088ee52213ac5ad6a2)) 


  at ~/.mendersoftware/authtoken
- Make 'token' argument global
 ([7361c94](https://github.com/mendersoftware/mender-cli/commit/7361c94d37f13857e9f3756fcf7f242af8f7b873))  by @kjaskiewiczz
- New artifacts command with upload subcommand
 ([9d3e3ee](https://github.com/mendersoftware/mender-cli/commit/9d3e3ee60fd1feb2d9c7b1a223a260b9371c9d09))  by @kjaskiewiczz




### Cmd/login


- Login subcommand
 ([f8d8282](https://github.com/mendersoftware/mender-cli/commit/f8d828232cfa811f34afbb8780e6623ba313b549)) 
- Add logging
 ([950a2c6](https://github.com/mendersoftware/mender-cli/commit/950a2c6ff0a9d05274e162fc824b904e01118966)) 
- Password prompt note
 ([3eb5a84](https://github.com/mendersoftware/mender-cli/commit/3eb5a84abc097197e442d9c70ff95d12866759de)) 
- Tweak permissions for authtoken dir/file
 ([97776d7](https://github.com/mendersoftware/mender-cli/commit/97776d72c17800f8e3ed65d14f0e367e6387b074)) 




### Cmd/root


- Hide arg names under shared consts
 ([f7365fd](https://github.com/mendersoftware/mender-cli/commit/f7365fda7655c8daa664d6350c3e61f2c3cbcbf0)) 
- Add skip-verify flag
 ([442e464](https://github.com/mendersoftware/mender-cli/commit/442e46404e12bd9feb2ae9b7430d82cf78ef80aa)) 


  skip SSL cert verification on demand - needed for everything
  done on the demo setup.
- Add verbose flag and setup logger
 ([3635d98](https://github.com/mendersoftware/mender-cli/commit/3635d98ce5865e19c70e4e57c9cb1d5eef329eec)) 




### Cmd/util


- Default authtoken path is ~/.mender
 ([fbaac4f](https://github.com/mendersoftware/mender-cli/commit/fbaac4f67de8ad81b26f697c4106cf4c34ee45ee)) 




### Log


- Add trivial logging layer
 ([0360400](https://github.com/mendersoftware/mender-cli/commit/0360400eacccca29f210e567a453802f1c04f8f3)) 


  acts as a global logger, prints errors to stderr,
  regular info msgs to stdout, and extra verbose msgs
  if flag is set.




### Vendor


- Howeyc/gopass
 ([c1f477b](https://github.com/mendersoftware/mender-cli/commit/c1f477bf52c4397bf6ec81be60b1be15b9cb57e6)) 
- Progress bar package
 ([73f4bb5](https://github.com/mendersoftware/mender-cli/commit/73f4bb5e6eec15aaccb993b9ed9f8575470061f2))  by @kjaskiewiczz




### Wording


- Replace 'backend' with 'server'
 ([c75a7f8](https://github.com/mendersoftware/mender-cli/commit/c75a7f875a4a55d3560f54c96a70127c38beb18e)) 





## v0.9.0 - 2018-04-05


### Cmd


- Spf13/cobra root command
 ([9c7e056](https://github.com/mendersoftware/mender-cli/commit/9c7e05699a19f9a35ac9d7f5c7efef32560be291)) 


  execute the command in main.
  
  the command is a stripped-down version of the autogenerated code:
  - removed example flags and some comments
  - removed config file/env handling; all params will be passed via
  cli args for now, maybe at some point we'll introduce extra config
  
  added a persitent --server flag (=available to subcommands, all need
  the url).
  
  also added a dummy second command - 'login'. apparently with just the
  'root' command, the nice well-structured app help is never displayed.




### Main


- Dummy binary for travis testing
 ([d0fbecd](https://github.com/mendersoftware/mender-cli/commit/d0fbecd1ac412fcea88541072a8cd2c52a0447ef)) 
- Build test
 ([4bc68aa](https://github.com/mendersoftware/mender-cli/commit/4bc68aae6dea147001ebef9b130a71fc5e372247)) 




### Travis


- Minimal travis file for testing
 ([a8512a7](https://github.com/mendersoftware/mender-cli/commit/a8512a77085bf959154eac14c923a8194776b855)) 
- Basic tests
 ([0b1c84c](https://github.com/mendersoftware/mender-cli/commit/0b1c84c88b6774cad5ea4af97177e4bc25436087)) 


  we don't want to invest in unit tests, but we need at least the license check -
  so perform it automatically as always.
  
  also include gocyclo, go fmt, go vet checks
  
  other stuff will be added incrementally (e.g. with acceptance tests - separate task).
- Setup github releases and s3 releases
 ([58f47bf](https://github.com/mendersoftware/mender-cli/commit/58f47bf1084c30762f4c155c9dd84b888082b948)) 


  general note: the travis 'deploy' step runs only on merges/tags (not PRs).
  
  a new JOB_TYPE=deploy is introduced - it runs via a job matrix and
  emits os/platform-specific binaries (2 total). it's separate from the basic
  tests job, which indeed only runs tests.
  
  a release consists of those 2 binaries (mender-cli.OS.ARCH) (ideas for
  the future: add bash completion files, markdown docs/man pages...).
  
  a github release is pushed strictly when a commit is tagged.
  
  because GH Releases don't support master/latest releases, we also release to s3 (TODO:
  a production, public s3 bucket; currently a dev account is used for the mender-cli-test bucket).
  
  an s3 release is pushed:
  - on every merge (a numbered build + 'latest' copy)
  - on tag (/release/TAG)
- Production s3 bucket
 ([6c8fca6](https://github.com/mendersoftware/mender-cli/commit/6c8fca6fd01a52c9103cb8254dd1819ef0122cfa)) 




### Vendor


- Init
 ([59a855a](https://github.com/mendersoftware/mender-cli/commit/59a855ab31ef82496a1a5d84935ce8de39fe179b)) 


  ...with mendertesting
- Spf13/cobra
 ([2fbe1c0](https://github.com/mendersoftware/mender-cli/commit/2fbe1c0d8099901ed43e9ad0addb08355203d101)) 





---
