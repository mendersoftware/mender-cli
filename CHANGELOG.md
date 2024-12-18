---
## 1.12.0-build4.1 - 2024-12-18


### Bug Fixes


- *(tests)* Use mongosh instead of mongo command
 ([03ade99](https://github.com/mendersoftware/mender-cli/commit/03ade99fa8c16d7e1e35a219d6f18c3921aef652))  by @kjaskiewiczz


  Legacy "mongo" shell has been removed in MongoDB 6.0






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





## 1.10.0-build7 - 2022-12-19


### Bug Fixes


- *(portforward)* Allow to get the server from mender-clirc configuration file
 ([a0d6c05](https://github.com/mendersoftware/mender-cli/commit/a0d6c05c7cbfa765e8cb74c6d74ab9522d951db3))  by @clementperon


  At the moment portforward doesn't take into account the mender clirc
  configuration file.
  
  Duplicate what's is done in mender terminal to get the proper configuration.






## 1.9.0-build8 - 2022-09-13


### Bug Fixes


- Stop the port-forward command on errors when reading the websocket
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([6231746](https://github.com/mendersoftware/mender-cli/commit/6231746cc6a211002dbb0c089bc053718c5c230b))  by @tranchitella


  There is no way to recover from failures when reading the websocket, as
  gorilla/websocket will panic with `repeated read on failed websocket
  connection`. Exit the command instead.
- Automatically handle reconnections in port-forward
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([77072de](https://github.com/mendersoftware/mender-cli/commit/77072de43c2898160127d77df8e99e8e7a481ac4))  by @tranchitella


  In case of errors, instead of stopping the process and exiting,
  automatically reconnect to the device and continue port-forwarding the
  connections. This improves the user experience in case of temporary
  errors when port-forwarding an HTTP server, where each request is a new
  connection. Additionally, ignore errors on connection close.
- Use a mutex lock per connection instead of a global one
([MEN-5565](https://northerntech.atlassian.net/browse/MEN-5565)) ([33ed7ae](https://github.com/mendersoftware/mender-cli/commit/33ed7aebeec428cecdd4291810126301ff437f20))  by @tranchitella




### Features


- Add option to specify a JWT or personal access token
([MEN-5660](https://northerntech.atlassian.net/browse/MEN-5660)) ([fa0134b](https://github.com/mendersoftware/mender-cli/commit/fa0134b4a56e0c46e9e7b6313de43713de8c1cf7))  by @tranchitella





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
([MEN-4305](https://northerntech.atlassian.net/browse/MEN-4305)) ([662cd74](https://github.com/mendersoftware/mender-cli/commit/662cd748dccb3c2983dbce248e9e16374a2b3727))  by @tranchitella




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
