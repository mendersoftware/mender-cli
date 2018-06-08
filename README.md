[![Build Status](https://travis-ci.org/mendersoftware/mender-cli.svg?branch=master)](https://travis-ci.org/mendersoftware/mender-cli)
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

* Join our [Google
  group](https://groups.google.com/a/lists.mender.io/forum/#!forum/mender)
* Follow us on [Twitter](https://twitter.com/mender_io?target=_blank). Please
  feel free to tweet us questions.
* Fork us on [Github](https://github.com/mendersoftware)
* Email us at [contact@mender.io](mailto:contact@mender.io)
