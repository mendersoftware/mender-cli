#!/usr/bin/python
# Copyright 2021 Northern.tech AS
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
from pathlib import Path
import os
import socket

import pytest

import docker

USER_HOME = str(Path.home())
DEFAULT_TOKEN_PATH = os.path.join(USER_HOME, ".cache", "mender", "authtoken")


def expect_output(stream, *expected):
    if isinstance(stream, list):
        stream = "\n".join(stream)
    for e in expected:
        assert e in stream, f'expected string "{e}" not found in stream: {stream}'
