---
## 2.0.0 - 2025-11-11


### Bug fixes


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
- Bump golang from 1.24.2 to 1.24.3
 ([21dce15](https://github.com/mendersoftware/mender-cli/commit/21dce159503f75dcd55555d85cbe74e74749e3d8))  by @dependabot[bot]


  Bumps golang from 1.24.2 to 1.24.3.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.24.3
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang from 1.24.3 to 1.24.4
 ([a5db4f2](https://github.com/mendersoftware/mender-cli/commit/a5db4f23a4ef424612901838bdad679cf3326e93))  by @dependabot[bot]


  Bumps golang from 1.24.3 to 1.24.4.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.24.4
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump github.com/go-viper/mapstructure/v2 from 2.2.1 to 2.3.0
 ([5bf7678](https://github.com/mendersoftware/mender-cli/commit/5bf76786bc57815de2945cfa557848bf8f84c653))  by @dependabot[bot]


  Bumps [github.com/go-viper/mapstructure/v2](https://github.com/go-viper/mapstructure) from 2.2.1 to 2.3.0.
  - [Release notes](https://github.com/go-viper/mapstructure/releases)
  - [Changelog](https://github.com/go-viper/mapstructure/blob/main/CHANGELOG.md)
  - [Commits](https://github.com/go-viper/mapstructure/compare/v2.2.1...v2.3.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/go-viper/mapstructure/v2
    dependency-version: 2.3.0
    dependency-type: indirect
  ...
- Bump golang from 1.24.4 to 1.24.5
 ([fb0eafb](https://github.com/mendersoftware/mender-cli/commit/fb0eafbe28da48f237b4bc431a7aa61fb04fbe40))  by @dependabot[bot]


  Bumps golang from 1.24.4 to 1.24.5.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.24.5
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump github.com/go-viper/mapstructure/v2 from 2.3.0 to 2.4.0
 ([ee231ae](https://github.com/mendersoftware/mender-cli/commit/ee231ae8c4489128315029fcc95b7988f4b46d75))  by @dependabot[bot]


  Bumps [github.com/go-viper/mapstructure/v2](https://github.com/go-viper/mapstructure) from 2.3.0 to 2.4.0.
  - [Release notes](https://github.com/go-viper/mapstructure/releases)
  - [Changelog](https://github.com/go-viper/mapstructure/blob/main/CHANGELOG.md)
  - [Commits](https://github.com/go-viper/mapstructure/compare/v2.3.0...v2.4.0)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/go-viper/mapstructure/v2
    dependency-version: 2.4.0
    dependency-type: indirect
  ...
- Bump github.com/ulikunitz/xz from 0.5.12 to 0.5.14
 ([1579e4c](https://github.com/mendersoftware/mender-cli/commit/1579e4c7e03806adc2676ee78ae66d4e24b30e38))  by @dependabot[bot]


  Bumps [github.com/ulikunitz/xz](https://github.com/ulikunitz/xz) from 0.5.12 to 0.5.14.
  - [Commits](https://github.com/ulikunitz/xz/compare/v0.5.12...v0.5.14)
  
  ---
  updated-dependencies:
  - dependency-name: github.com/ulikunitz/xz
    dependency-version: 0.5.14
    dependency-type: indirect
  ...
- Bump golang from 1.24.5 to 1.25.0
 ([8c7bd26](https://github.com/mendersoftware/mender-cli/commit/8c7bd268e2ef62a4b6323c3dd07cd9eddca86175))  by @dependabot[bot]


  Bumps golang from 1.24.5 to 1.25.0.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.25.0
    dependency-type: direct:production
    update-type: version-update:semver-minor
  ...
- Bump golang from 1.25.0 to 1.25.1
 ([34af25b](https://github.com/mendersoftware/mender-cli/commit/34af25be565c22f20e50a72eca7e3ac18126624f))  by @dependabot[bot]


  Bumps golang from 1.25.0 to 1.25.1.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.25.1
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump golang from 1.25.1 to 1.25.3
 ([2e1cfe4](https://github.com/mendersoftware/mender-cli/commit/2e1cfe4ff7da4f4046839a15105743ada29d46a5))  by @dependabot[bot]


  Bumps golang from 1.25.1 to 1.25.3.
  
  ---
  updated-dependencies:
  - dependency-name: golang
    dependency-version: 1.25.3
    dependency-type: direct:production
    update-type: version-update:semver-patch
  ...
- Bump tests/mender_server from `a0df4cb` to `cb41636`
 ([c262bd1](https://github.com/mendersoftware/mender-cli/commit/c262bd17f0a1c3084b9ee94689924b88531ab5aa))  by @dependabot[bot]


  Bumps [tests/mender_server](https://github.com/mendersoftware/mender-server) from `a0df4cb` to `cb41636`.
  - [Release notes](https://github.com/mendersoftware/mender-server/releases)
  - [Commits](https://github.com/mendersoftware/mender-server/compare/a0df4cb00305fbf3ef10cbe5ecb1c45c93b01853...cb41636922a9755426c1d41f22105954f0cd8c1e)
  
  ---
  updated-dependencies:
  - dependency-name: tests/mender_server
    dependency-version: cb41636922a9755426c1d41f22105954f0cd8c1e
    dependency-type: direct:production
  ...
- Bump python from 3.13-slim to 3.14-slim
 ([953826c](https://github.com/mendersoftware/mender-cli/commit/953826c7a014666cc5646c2e4f3f3dc00a2b8ecf))  by @dependabot[bot]


  Bumps python from 3.13-slim to 3.14-slim.
  
  ---
  updated-dependencies:
  - dependency-name: python
    dependency-version: 3.14-slim
    dependency-type: direct:production
  ...




### Refac


- Move acceptance tests to `tests/acceptance`
 ([cc53528](https://github.com/mendersoftware/mender-cli/commit/cc5352895c5e94b6001f0d76fc8e91e642c39205))  by @lluiscampos







## mender-cli 1.12.0

_Released 12.28.2023_

### Statistics

| Developers with the most changesets | |
|---|---|
| Fredrik Flornes Ellertsen | 2 (20.0%) |
| Manuel Zedel | 2 (20.0%) |
| Roberto Giovanardi | 1 (10.0%) |
| Lluis Campos | 1 (10.0%) |
| Daniel Skinstad Drabitzius | 1 (10.0%) |
| Luis Ramirez Vargas | 1 (10.0%) |
| Peter Grzybowski | 1 (10.0%) |
| Steven Leadbeater | 1 (10.0%) |

| Developers with the most changed lines | |
|---|---|
| Steven Leadbeater | 313 (47.5%) |
| Peter Grzybowski | 254 (38.5%) |
| Manuel Zedel | 52 (7.9%) |
| Daniel Skinstad Drabitzius | 16 (2.4%) |
| Fredrik Flornes Ellertsen | 13 (2.0%) |
| Roberto Giovanardi | 4 (0.6%) |
| Lluis Campos | 4 (0.6%) |
| Luis Ramirez Vargas | 3 (0.5%) |

| Developers with the most lines removed | |
|---|---|
| Manuel Zedel | 19 (17.3%) |
| Fredrik Flornes Ellertsen | 6 (5.5%) |

Developers with the most signoffs (total 1)
|---|---|
| Kristian Amlie | 1 (100.0%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 7 (70.0%) |
| fellerts@fastmail.com | 2 (20.0%) |
| stevenleadbeater@live.co.uk | 1 (10.0%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 333 (50.5%) |
| stevenleadbeater@live.co.uk | 313 (47.5%) |
| fellerts@fastmail.com | 13 (2.0%) |

| Employers with the most signoffs (total 1) | |
|---|---|
| Northern.tech | 1 (100.0%) |

| Employers with the most hackers (total 8) | |
|---|---|
| Northern.tech | 6 (75.0%) |
| stevenleadbeater@live.co.uk | 1 (12.5%) |
| fellerts@fastmail.com | 1 (12.5%) |

### Changelogs

#### mender-cli (1.12.0)

New changes in mender-cli since 1.11.1:

##### Bug fixes

* Keep stdout clean
* Use term's built-in IsTerminal

##### Features

* Adding artifact downloads to CLI
* Add meaningful error message for conflict artifact


## mender-cli 1.11.1

_Released 10.18.2023_

### Statistics

| Developers with the most changesets | |
|---|---|
| Peter Grzybowski | 1 (100.0%) |

| Developers with the most changed lines | |
|---|---|
| Peter Grzybowski | 272 (100.0%) |

| Developers with the most signoffs (total 1) | |
|---|---|
| Fabio Tranchitella | 1 (100.0%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 1 (100.0%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 272 (100.0%) |

| Employers with the most signoffs (total 1) | |
|---|---|
| Northern.tech | 1 (100.0%) |

| Employers with the most hackers (total 1) | |
|---|---|
| Northern.tech | 1 (100.0%) |

### Changelogs

#### mender-cli (1.11.1)

New changes in mender-cli since 1.11.0:

##### Features

* post metadata with the direct upload with skip verify
  ([MEN-6696](https://northerntech.atlassian.net/browse/MEN-6696))


## mender-cli 1.11.0

_Released 07.28.2023_

### Statistics

A total of 494 lines added, 1194 removed (delta -700)

| Developers with the most changesets | |
|---|---|
| Lluis Campos | 5 (35.7%) |
| Sven Schermer | 3 (21.4%) |
| Peter Grzybowski | 3 (21.4%) |
| Krzysztof Jaskiewicz | 2 (14.3%) |
| Esteban Agüero Pérez | 1 (7.1%) |

| Developers with the most changed lines | |
|---|---|
| Krzysztof Jaskiewicz | 979 (67.2%) |
| Peter Grzybowski | 396 (27.2%) |
| Lluis Campos | 55 (3.8%) |
| Sven Schermer | 25 (1.7%) |
| Esteban Agüero Pérez | 1 (0.1%) |

| Developers with the most lines removed | |
|---|---|
| Krzysztof Jaskiewicz | 961 (80.5%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 10 (71.4%) |
| Disruptive Technologies | 3 (21.4%) |
| estape11@gmail.com | 1 (7.1%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 1430 (98.2%) |
| Disruptive Technologies | 25 (1.7%) |
| estape11@gmail.com | 1 (0.1%) |

| Employers with the most hackers (total 5) | |
|---|---|
| Northern.tech | 4 (60.0%) |
| Disruptive Technologies | 1 (20.0%) |
| estape11@gmail.com | 1 (20.0%) |

### Changelogs

#### mender-cli (1.11.0)

New changes in mender-cli since 1.10.0:

##### Bug fixes

* Hide errors when help flag is present
  ([MEN-6357](https://northerntech.atlassian.net/browse/MEN-6357))
* Allow to use --token to specify API token
  ([MEN-6357](https://northerntech.atlassian.net/browse/MEN-6357))
* List devices in raw mode not in stdout

##### Features

* direct upload.
  ([MEN-6338](https://northerntech.atlassian.net/browse/MEN-6338))

##### Other

* Update container base image to golang 1.20
* formatting (go fmt).
  ([MEN-6338](https://northerntech.atlassian.net/browse/MEN-6338))


## mender-cli 1.10.0

_Released 02.20.2023_

### Statistics

A total of 11 lines added, 10 removed (delta 1)

| Developers with the most changesets | |
|---|---|
| Lluis Campos | 4 (66.7%) |
| Clément Péron | 1 (16.7%) |
| Alex Miliukov | 1 (16.7%) |

| Developers with the most changed lines | |
|---|---|
| Lluis Campos | 6 (50.0%) |
| Clément Péron | 4 (33.3%) |
| Alex Miliukov | 2 (16.7%) |

| Developers with the most lines removed | |
|---|---|
| Lluis Campos | 1 (10.0%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 5 (83.3%) |
| peron.clem@gmail.com | 1 (16.7%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 8 (66.7%) |
| peron.clem@gmail.com | 4 (33.3%) |

| Employers with the most hackers (total 3) | |
|---|---|
| Northern.tech | 2 (66.7%) |
| peron.clem@gmail.com | 1 (33.3%) |

### Changelogs

#### mender-cli (1.10.0)

New changes in mender-cli since 1.9.0:

##### Bug fixes

* allow to get the server from mender-clirc configuration file


## mender-cli 1.9.0

_Released 09.25.2022_

### Statistics

A total of 171 lines added, 188 removed (delta -17)

| Developers with the most changesets | |
|---|---|
| Fabio Tranchitella | 4 (44.4%) |
| Ole Petter Orhagen | 2 (22.2%) |
| Lluis Campos | 1 (11.1%) |
| Alex Miliukov | 1 (11.1%) |
| Manuel Zedel | 1 (11.1%) |

| Developers with the most changed lines | |
|---|---|
| Fabio Tranchitella | 204 (85.7%) |
| Ole Petter Orhagen | 14 (5.9%) |
| Alex Miliukov | 13 (5.5%) |
| Manuel Zedel | 6 (2.5%) |
| Lluis Campos | 1 (0.4%) |

| Developers with the most lines removed | |
|---|---|
| Fabio Tranchitella | 39 (20.7%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 9 (100.0%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 238 (100.0%) |

| Employers with the most hackers (total 5) | |
|---|---|
| Northern.tech | 5 (100.0%) |

### Changelogs

#### mender-cli (1.9.0)

New changes in mender-cli since 1.8.0:

##### Bug fixes

* stop the port-forward command on errors when reading the websocket
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))
* automatically handle reconnections in port-forward
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))
* use a mutex lock per connection instead of a global one
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))

##### Features

* add option to specify a JWT or personal access token
  ([MEN-5660](https://northerntech.atlassian.net/browse/MEN-5660))


## mender-cli 1.8.2

_Released 03.10.2023_

### Statistics

A total of 17 lines added, 4 removed (delta 13)

| Developers with the most changesets | |
|---|---|
| Ole Petter Orhagen | 1 (50.0%) |
| Clément Péron | 1 (50.0%) |

| Developers with the most changed lines | |
|---|---|
| Ole Petter Orhagen | 13 (76.5%) |
| Clément Péron | 4 (23.5%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 1 (50.0%) |
| peron.clem@gmail.com | 1 (50.0%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 13 (76.5%) |
| peron.clem@gmail.com | 4 (23.5%) |

| Employers with the most hackers (total 2) | |
|---|---|
| Northern.tech | 1 (50.0%) |
| peron.clem@gmail.com | 1 (50.0%) |

### Changelogs

#### mender-cli (1.8.2)

New changes in mender-cli since 1.8.1:

##### Bug fixes

* allow to get the server from mender-clirc configuration file


## mender-cli 1.8.1

_Released 10.19.2022_

### Statistics

A total of 57 lines added, 30 removed (delta 27)

| Developers with the most changesets | |
|---|---|
| Fabio Tranchitella | 4 (100.0%) |

| Developers with the most changed lines | |
|---|---|
| Fabio Tranchitella | 58 (100.0%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 4 (100.0%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 58 (100.0%) |

| Employers with the most hackers (total 1) | |
|---|---|
| Northern.tech | 1 (100.0%) |

### Changelogs

#### mender-cli (1.8.1)

New changes in mender-cli since 1.8.0:

##### Bug fixes

* stop the port-forward command on errors when reading the websocket
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))
* automatically handle reconnections in port-forward
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))
* use a mutex lock per connection instead of a global one
  ([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565))


## mender-cli 1.8.0

_Released 06.14.2022_

### Statistics

A total of 484 lines added, 172 removed (delta 312)

| Developers with the most changesets | |
|---|---|
| Lluis Campos | 4 (36.4%) |
| Fabio Tranchitella | 2 (18.2%) |
| Ruben Schwarz | 2 (18.2%) |
| Mikael Torp-Holte | 1 (9.1%) |
| Maciej Tomczuk | 1 (9.1%) |
| Kristian Amlie | 1 (9.1%) |

| Developers with the most changed lines | |
|---|---|
| Lluis Campos | 199 (40.6%) |
| Fabio Tranchitella | 115 (23.5%) |
| Ruben Schwarz | 114 (23.3%) |
| Maciej Tomczuk | 59 (12.0%) |
| Mikael Torp-Holte | 2 (0.4%) |
| Kristian Amlie | 1 (0.2%) |

| Top changeset contributors by employer | |
|---|---|
| Northern.tech | 9 (81.8%) |
| SOTEC | 2 (18.2%) |

| Top lines changed by employer | |
|---|---|
| Northern.tech | 376 (76.7%) |
| SOTEC | 114 (23.3%) |

| Employers with the most hackers (total 6) | |
|---|---|
| Northern.tech | 5 (83.3%) |
| SOTEC | 1 (16.7%) |

### Changelogs

#### mender-cli (1.8.0)

New changes in mender-cli since 1.7.0:

* New command `mender-cli artifacts delete ID` to delete an
  artifact
* Raw mode for devices list command
* disable `Failed to parse flags: unknown flag` error messages
  ([MEN-5428](https://northerntech.atlassian.net/browse/MEN-5428))
* improve auth and device error messages for the troubleshoot commands
  ([MEN-5428](https://northerntech.atlassian.net/browse/MEN-5428))

##### Dependabot bumps

* Aggregated Dependabot Changelogs:
  * Bumps golang from 1.16.2-alpine3.12 to 1.16.3-alpine3.12.
  * Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.7 to 3.0.8.
      - [Release notes](https://github.com/cheggaaa/pb/releases)
      - [Commits](https://github.com/cheggaaa/pb/compare/v3.0.7...v3.0.8)
  * Bumps golang from 1.16.3-alpine3.12 to 1.16.4-alpine3.12.
  * Bumps golang from 1.16.4-alpine3.12 to 1.16.5-alpine3.12.

      ```
      updated-dependencies:
      - dependency-name: golang
        dependency-type: direct:production
        update-type: version-update:semver-patch
      ```
  * Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.7.1 to 1.8.0.
      - [Release notes](https://github.com/spf13/viper/releases)
      - [Commits](https://github.com/spf13/viper/compare/v1.7.1...v1.8.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/viper
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.8.0 to 1.8.1.
      - [Release notes](https://github.com/spf13/viper/releases)
      - [Commits](https://github.com/spf13/viper/compare/v1.8.0...v1.8.1)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/viper
        dependency-type: direct:production
        update-type: version-update:semver-patch
      ```
  * Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.1.3 to 1.2.1.
      - [Release notes](https://github.com/spf13/cobra/releases)
      - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
      - [Commits](https://github.com/spf13/cobra/compare/v1.1.3...v1.2.1)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/cobra
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/google/uuid](https://github.com/google/uuid) from 1.2.0 to 1.3.0.
      - [Release notes](https://github.com/google/uuid/releases)
      - [Commits](https://github.com/google/uuid/compare/v1.2.0...v1.3.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/google/uuid
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.8.1 to 1.9.0.
      - [Release notes](https://github.com/spf13/viper/releases)
      - [Commits](https://github.com/spf13/viper/compare/v1.8.1...v1.9.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/viper
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.9.0 to 1.10.0.
      - [Release notes](https://github.com/spf13/viper/releases)
      - [Commits](https://github.com/spf13/viper/compare/v1.9.0...v1.10.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/viper
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.2.1 to 1.3.0.
      - [Release notes](https://github.com/spf13/cobra/releases)
      - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
      - [Commits](https://github.com/spf13/cobra/compare/v1.2.1...v1.3.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/cobra
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.10.0 to 1.10.1.
      - [Release notes](https://github.com/spf13/viper/releases)
      - [Commits](https://github.com/spf13/viper/compare/v1.10.0...v1.10.1)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/viper
        dependency-type: direct:production
        update-type: version-update:semver-patch
      ```
  * Bumps [github.com/gorilla/websocket](https://github.com/gorilla/websocket) from 1.4.2 to 1.5.0.
      - [Release notes](https://github.com/gorilla/websocket/releases)
      - [Commits](https://github.com/gorilla/websocket/compare/v1.4.2...v1.5.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/gorilla/websocket
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```
  * Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.3.0 to 1.4.0.
      - [Release notes](https://github.com/spf13/cobra/releases)
      - [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
      - [Commits](https://github.com/spf13/cobra/compare/v1.3.0...v1.4.0)

      ```
      updated-dependencies:
      - dependency-name: github.com/spf13/cobra
        dependency-type: direct:production
        update-type: version-update:semver-minor
      ```


## mender-cli 1.7.0

_Released 04.16.2021_

### Changelogs

#### mender-cli (1.7.0)

New changes in mender-cli since 1.6.0:

* Fix: Respect the --server flag from config everywhere
* `mender-cli --record <my-file> terminal <DEVICE-ID>` records
the terminal session into a local file.
([MEN-4318](https://northerntech.atlassian.net/browse/MEN-4318))
* `mender-cli --playback <my-file> terminal` playbacks the
previously recorded terminal session from a local file.
([MEN-4318](https://northerntech.atlassian.net/browse/MEN-4318))
* New command `mender-cli devices list` to list all devices
from /devauth/devices endpoint. The amount of detail can be controlled
using cli parameter `-d/--detail`, same as for other commands.
* Previously, the --generate-autocomplete call would silently ignore
errors, when the autocomplete directory was not present. This explicitly logs
the errors returned during autocomplete script generation.
* New command port-forward: port-forward TCP and UDP ports from the device
* Add filetransfer upload and download support
([MEN-4323](https://northerntech.atlassian.net/browse/MEN-4323))
* Aggregated Dependabot Changelogs:
* Bumps golang from 1.15.6-alpine3.12 to 1.15.8-alpine3.12.
* Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.1.1 to 1.1.3.
- [Release notes](https://github.com/spf13/cobra/releases)
- [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
- [Commits](https://github.com/spf13/cobra/compare/v1.1.1...v1.1.3)
* Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.5 to 3.0.6.
- [Release notes](https://github.com/cheggaaa/pb/releases)
- [Commits](https://github.com/cheggaaa/pb/compare/v3.0.5...v3.0.6)
* Bumps golang from 1.15.8-alpine3.12 to 1.16.0-alpine3.12.
* Bumps golang from 1.16.0-alpine3.12 to 1.16.2-alpine3.12.
* Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.6 to 3.0.7.
- [Release notes](https://github.com/cheggaaa/pb/releases)
- [Commits](https://github.com/cheggaaa/pb/compare/v3.0.6...v3.0.7)

## mender-cli 1.6.1

_Released 16.04.2021_

### Changelogs

#### mender-cli (1.6.1)

New changes in mender-cli since 1.6.0:

* Fix: Respect the --server flag from config everywhere

## mender-cli 1.6.0

_Released 01.20.2021_

### Changelogs

#### mender-cli (1.6.0)

New changes in mender-cli since 1.5.0:

* Fix login with password, and improve the configuration file handling
* Add 'artifact list' command to list Artifacts on the Mender server
* Prompt for the users username if not provided on the CLI
* Add '--version' option to display the current mender-cli version
* New CLI command "terminal" to access a device's remote terminal
* mender-cli requires now golang 1.15 or newer
* Aggregated Dependabot Changelogs:
* Bumps golang from 1.14-alpine3.12 to 1.15.1-alpine3.12.
* Bump golang from 1.14-alpine3.12 to 1.15.1-alpine3.12
* Bumps golang from 1.15.1-alpine3.12 to 1.15.2-alpine3.12.
* Bump golang from 1.15.1-alpine3.12 to 1.15.2-alpine3.12
* Bumps [github.com/cheggaaa/pb/v3](https://github.com/cheggaaa/pb) from 3.0.4 to 3.0.5.
- [Release notes](https://github.com/cheggaaa/pb/releases)
- [Commits](https://github.com/cheggaaa/pb/compare/v3.0.4...v3.0.5)
* Bump github.com/cheggaaa/pb/v3 from 3.0.4 to 3.0.5
* Bumps golang from 1.15.2-alpine3.12 to 1.15.3-alpine3.12.
* Bumps [github.com/spf13/cobra](https://github.com/spf13/cobra) from 1.0.0 to 1.1.1.
- [Release notes](https://github.com/spf13/cobra/releases)
- [Changelog](https://github.com/spf13/cobra/blob/master/CHANGELOG.md)
- [Commits](https://github.com/spf13/cobra/compare/v1.0.0...v1.1.1)
* Bump golang from 1.15.2-alpine3.12 to 1.15.3-alpine3.12
* Bump github.com/spf13/cobra from 1.0.0 to 1.1.1
* Bumps golang from 1.15.3-alpine3.12 to 1.15.4-alpine3.12.
* Bump golang from 1.15.3-alpine3.12 to 1.15.4-alpine3.12
* Bumps golang from 1.15.4-alpine3.12 to 1.15.5-alpine3.12.
* Bump golang from 1.15.4-alpine3.12 to 1.15.5-alpine3.12
* Bumps golang from 1.15.5-alpine3.12 to 1.15.6-alpine3.12.
* Bump golang from 1.15.5-alpine3.12 to 1.15.6-alpine3.12

## mender-cli 1.5.1

_Released 01.21.2021_

### Changelogs

#### mender-cli (1.5.1)

New changes in mender-cli since 1.5.0:

* Fixed login with password


## mender-cli 1.5.0

_Released 09.11.2020_

### Changelogs

#### mender-cli (1.5.0)

New changes in mender-cli since 1.4.0:

* Add: Make the server flag default to hosted Mender
* Add: Bash auto-completion functionality
* Add: Zsh auto-completion support
* Add: Configuration file functionality
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

## mender-cli 1.4.1

_Released 16.04.2021_

### Changelogs

#### mender-cli (1.4.1)

New changes in mender-cli since 1.4.0:

* Bump golang version to 1.14-alpine3.12


## mender-cli 1.4.0

_Released 07.15.2020_

### Changelogs

#### mender-cli (1.4.0)

New changes in mender-cli since 1.3.0:

* Support for two factor authentication token for login
([MEN-3176](https://northerntech.atlassian.net/browse/MEN-3176))
* Change the name of the two-factor auth option.

## mender-cli 1.3.0

_Released 03.05.2020_

### Changelogs

#### mender-cli (1.3.0)

New changes in mender-cli since 1.2.0:

* Build and publish Mac OS X binary for `mender-cli`

## mender-cli 1.2.0

_Released 09.16.2019_

### Changelogs

#### mender-cli (1.2.0)

New changes in mender-cli since 1.1.0:

* Store login token in XDG Basedir Spec Cache-directory
([MEN-2387](https://northerntech.atlassian.net/browse/MEN-2387))

---
