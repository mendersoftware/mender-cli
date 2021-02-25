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
import os
from pathlib import Path

import pytest

import cli
from common import single_user, expect_output, DEFAULT_TOKEN_PATH
import api
import docker
import crypto


@pytest.fixture(scope="function")
def logged_in_single_user(single_user):
    c = cli.Cli()
    r = c.run(
        "login",
        "--server",
        "https://mender-api-gateway",
        "--skip-verify",
        "--username",
        "user@tenant.com",
        "--password",
        "youcantguess",
    )

    assert r.returncode == 0, r.stderr
    yield
    os.remove(DEFAULT_TOKEN_PATH)


@pytest.fixture(scope="function")
def clean_devauth_db():
    yield
    r = docker.exec(
        "mender-mongo",
        docker.BASE_COMPOSE_FILES,
        "mongo",
        "deviceauth",
        "--eval",
        "db.dropDatabase()",
    )
    assert r.returncode == 0, r.stderr


@pytest.mark.usefixtures("clean_devauth_db")
class TestDevicesList:
    def test_ok(self, logged_in_single_user):
        # Create a device
        authset = crypto.rand_id_data()
        keypair = crypto.get_keypair_rsa()
        dapi = api.DevicesDevauth()
        r = dapi.post_auth_request(authset, keypair[1], keypair[0])
        assert r.status_code == 401

        # Device should be listed as pending
        c = cli.Cli()
        r = c.run(
            "--server", "https://mender-api-gateway", "--skip-verify", "devices", "list"
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stdout, "Status: pending")

        # Get and accept the device
        token = Path(DEFAULT_TOKEN_PATH).read_text()
        mapi = api.ManagementDevauth(token)
        r = mapi.get_devices()
        assert r.status_code == 200
        devices = r.json()
        device_id = devices[0]["id"]
        authset_id = devices[0]["auth_sets"][0]["id"]
        r = mapi.set_device_auth_status(device_id, authset_id, "accepted")
        assert r.status_code == 204

        # Device should now be listed as accepted
        c = cli.Cli()
        r = c.run(
            "--server", "https://mender-api-gateway", "--skip-verify", "devices", "list"
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stdout, "Status: accepted")
