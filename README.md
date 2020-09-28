[![Build Status](https://gitlab.com/Northern.tech/Mender/mender-cli/badges/master/pipeline.svg)](https://gitlab.com/Northern.tech/Mender/mender-cli/pipelines)
[![Coverage Status](https://coveralls.io/repos/github/mendersoftware/mender-cli/badge.svg?branch=master)](https://coveralls.io/github/mendersoftware/mender-cli?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mendersoftware/mender-cli)](https://goreportcard.com/report/github.com/mendersoftware/mender-cli)


Mender CLI
========================

Mender is an open source over-the-air (OTA) software updater for embedded Linux
devices. Mender comprises a client running at the embedded device, as well as
a server that manages deployments across many devices.

This repository contains a standalone tool that makes it much easier to work
with the [Mender server management APIs](https://docs.mender.io/apis/management-apis).

The goal with `mender-cli` is to simplify integration between the Mender server
and cloud services like continuous integration (CI)/build automation.

Over time `mender-cli` will be extended to simplify the most common use cases
for integrating the Mender server into other backend and cloud systems.


## Getting started

To start using `mender-cli`, we recommend that you begin with the
[documentation section to set up mender-cli](https://docs.mender.io/server-integration/using-the-apis#set-up-mender-cli).


## Downloading the binaries

You can find the latest `mender-cli` binaries in the [Downloads page on Mender
Docs](https://docs.mender.io/downloads).

## Configuration file

The `mender-cli` tool supports having a custom configuration setup. For now it
supports the `username`, `password`, and `server` configuration parameters. The
file must be in the JSON format, and can be located in one of the following
directories:

* `/etc/mender-cli/.mender-clirc`
* `$HOME/.mender-clirc`
* `.` (The directory in which the binary is run from)

Example configuration file:

```json
{
    "username": "foo@bar.com",
    "password": "baz",
    "server"  : "bar.com"
}
```

!!! Note: It is possible to override all configuration file parameters on the command line.

## Autocompletion

Autocompletion can be enabled for the `mender-cli` tool through one of two ways.

1. `sudo make install-autocomplete-scripts`

This is the simplest option, and will automatically generate and install the
auto-completion scripts for both `Bash` and `Zsh`, and install them into the
appropriate directories. That is `Bash` autocompletion goes into the
`/etc/bash_completion.d/` directory, and `Zsh` scripts (if `Zsh` is installed),
gets installed into the `/usr/local/share/zsh/site-functions/` directory.

2. Manually

### Enabling Bash auto-complete manually

The `mender-cli` tool can be enabled to support shell autocompletion, like you
are used to for your regular tools, like `git`, `cd`, etc. In order to enable
this functionality run the `mender-cli` tool with the `--generate` flag. Example:

```console
mender-cli --generate
```

This will output the file `autocomplete.sh` in the `autocomplete` directory in
the directory the binary is run from, so it is recommended to run this from the
`mender-cli` source directory, where this directory already exists.

In order to enable the functionality the current `Bash` shell has to pick up the
completions (i.e., source it). This can be done in one of two ways:

1. Copy the `./autocomplete/autocomplete.sh` file to `/etc/bash_completion.d/`,
   where a new `Bash` shell will automatically source it on invocation.

2. Keep the file where it is, and have each new login shell source the file as a
   part of `.bashrc`. This means that `echo "source
   /path/to/mender-cli/autocomplete/autocomplete.sh" >> ~/.bashrc` needs to be
   present in your `Bash` config.

### Enabling Zsh auto-complete manually

The `mender-cli` tool supports enabling `Zsh` support, just like [Bash
auto-completion](#enabling-bash-auto-complete-manually).

In order to enable `Zsh` auto-completion do:

```console
mender-cli --generate
```

to generate the auto-completion script.

```console
echo $FPATH
```

Choose one of the directories in the `$FPATH`, and then copy the
`./autocomplete/autocomplete.zsh` script into this directory, and rename it to
`_mender-cli`, and restart your shell. Now typing:

```console
mender-cli [Tab]
```

should result in:

```console
$ mender-cli [Tab]
artifacts  -- Operations on mender artifacts.
login      -- Log in to the Mender server (required before other operation
```

## Contributing

We welcome and ask for your contribution. If you would like to contribute to
Mender, please read our guide on how to best get started [contributing code or
documentation](https://github.com/mendersoftware/mender/blob/master/CONTRIBUTING.md).


## License

Mender is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/mendersoftware/mender-cli/blob/master/LICENSE) for
the full license text.


## Security disclosure

We take security very seriously. If you come across any issue regarding
security, please disclose the information by sending an email to
[security@mender.io](security@mender.io). Please do not create a new public
issue. We thank you in advance for your cooperation.


## Connect with us

* Join the [Mender Hub discussion forum](https://hub.mender.io)
* Follow us on [Twitter](https://twitter.com/mender_io). Please
  feel free to tweet us questions.
* Fork us on [Github](https://github.com/mendersoftware)
* Create an issue in the [bugtracker](https://tracker.mender.io/projects/MEN)
* Email us at [contact@mender.io](mailto:contact@mender.io)
* Connect to the [#mender IRC channel on Freenode](http://webchat.freenode.net/?channels=mender)
