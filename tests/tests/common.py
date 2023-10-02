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

import pytest

import docker

USER_HOME = str(Path.home())
DEFAULT_TOKEN_PATH = os.path.join(USER_HOME, ".cache", "mender", "authtoken")


@pytest.fixture(scope="class")
def single_user():
    r = docker.exec(
        "mender-useradm",
        docker.BASE_COMPOSE_FILES,
        "/usr/bin/useradm",
        "create-user",
        "--username",
        "user@tenant.com",
        "--password",
        "youcantguess",
    )

    assert r.returncode == 0, r.stderr
    yield
    clean_useradm_db()


def clean_useradm_db():
    r = docker.exec(
        "mender-mongo",
        docker.BASE_COMPOSE_FILES,
        "mongo",
        "useradm",
        "--eval",
        "db.dropDatabase()",
    )

    assert r.returncode == 0, r.stderr


def expect_output(stream, *expected):
    for e in expected:
        assert e in stream, f'expected string "{e}" not found in stream: {stream}'
