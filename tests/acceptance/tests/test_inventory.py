#!/usr/bin/python
# Copyright 2026 Northern.tech AS
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
from common import expect_output, DEFAULT_TOKEN_PATH
import api
import crypto
import json


@pytest.fixture(scope="function")
def logged_in_single_user(single_user):
    c = cli.Cli()
    r = c.run(
        "login",
        "--server",
        "https://docker.mender.io",
        "--skip-verify",
        "--username",
        "user@tenant.com",
        "--password",
        "youcantguess",
    )

    assert r.returncode == 0, r.stderr
    yield
    os.remove(DEFAULT_TOKEN_PATH)


def create_device(attributes):
    authset = crypto.rand_id_data()
    keypair = crypto.get_keypair_rsa()

    token = Path(DEFAULT_TOKEN_PATH).read_text()
    mapi = api.ManagementDevauth(token)
    r = mapi.preauthorize(authset, keypair[1])
    assert r.status_code == 201

    dapi = api.DevicesDevauth()
    r = dapi.post_auth_request(authset, keypair[1], keypair[0])
    assert r.status_code == 200
    device_token = r.text

    r = mapi.get_devices()
    assert r.status_code == 200
    device = next(d for d in r.json() if d["identity_data"] == authset)

    r = api.DevicesInventory().patch_attributes(device_token, attributes)
    assert r.status_code == 200, r.text
    return device["id"], authset


@pytest.fixture(scope="function")
def device_with_inventory(logged_in_single_user):
    yield create_device(
        [
            {"name": "device_type", "value": "test-device"},
            {"name": "hostname", "value": "pos-1"},
            {"name": "ipv4_eth0", "value": ["10.0.0.1/24", "10.0.0.2/24"]},
        ]
    )


class TestDevicesInventoryList:
    def test_ok(self, device_with_inventory):
        device_id, authset = device_with_inventory

        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--filter",
            "device_type=test-device",
            "--filter",
            "identity/status=accepted",
            "--columns",
            "hostname,identity/mac,ipv4_eth0",
            "--output",
            "csv",
        )
        assert r.returncode == 0, r.stderr
        expect_output(
            r.stdout,
            "id,hostname,identity/mac,ipv4_eth0\n",
            f"{device_id},pos-1,{authset['mac']},10.0.0.1/24;10.0.0.2/24\n",
        )

        # Device should not match a filter on other device types
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--filter",
            "device_type=in:other,another",
            "--columns",
            "hostname",
            "--output",
            "csv",
        )
        assert r.returncode == 0, r.stderr
        assert r.stdout == "id,hostname\n"

        # Device should be listed when fetching all pages
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--all",
            "--filter",
            "hostname=~pos-",
            "--columns",
            "hostname",
            "--output",
            "csv",
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stdout, f"{device_id},pos-1\n")

        # Device should not be listed by a filter on a missing attribute
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--filter",
            "!hostname?",
            "--output",
            "json",
        )
        assert r.returncode == 0, r.stderr
        assert json.loads(r.stdout) == []

        # Only the selected attributes should be returned in raw mode
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--columns",
            "hostname",
            "--raw",
        )
        assert r.returncode == 0, r.stderr
        raw = json.loads(r.stdout, strict=False)
        assert raw[0]["id"] == device_id
        names = [a["name"] for a in raw[0]["attributes"]]
        assert "hostname" in names
        assert "device_type" not in names

    def test_ok_all_sorted_across_pages(self, logged_in_single_user):
        ids = {}
        for hostname in ["pos-b", "pos-c", "pos-a"]:
            device_id, _ = create_device(
                [
                    {"name": "device_type", "value": "sort-test"},
                    {"name": "hostname", "value": hostname},
                ]
            )
            ids[hostname] = device_id

        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--all",
            "--per-page",
            "1",
            "--filter",
            "device_type=sort-test",
            "--sort",
            "hostname:desc",
            "--columns",
            "hostname",
            "--output",
            "csv",
        )
        assert r.returncode == 0, r.stderr
        assert r.stdout == (
            "id,hostname\n"
            f"{ids['pos-c']},pos-c\n"
            f"{ids['pos-b']},pos-b\n"
            f"{ids['pos-a']},pos-a\n"
        )

    def test_ok_output_file(self, device_with_inventory, tmp_path):
        device_id, _ = device_with_inventory
        out_file = tmp_path / "devices.csv"

        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--all",
            "--progress",
            "--filter",
            "device_type=test-device",
            "--columns",
            "hostname",
            "--output",
            "csv",
            "--output-file",
            str(out_file),
        )
        assert r.returncode == 0, r.stderr
        assert r.stdout == ""
        expect_output(r.stderr, "Fetched 1/1 devices (100%)")
        assert out_file.read_text() == f"id,hostname\n{device_id},pos-1\n"
        assert [p.name for p in tmp_path.iterdir()] == ["devices.csv"]

    def test_ok_count(self, device_with_inventory):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--count",
            "--filter",
            "device_type=test-device",
        )
        assert r.returncode == 0, r.stderr
        assert r.stdout == "1\n"

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--count",
            "--filter",
            "device_type=other",
        )
        assert r.returncode == 0, r.stderr
        assert r.stdout == "0\n"

    def test_error_invalid_filter(self, logged_in_single_user):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--filter",
            "hostname",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "invalid filter")

    def test_error_progress_without_all(self, logged_in_single_user):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
            "--progress",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "--progress requires --all")

    def test_error_no_login(self):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "list",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")


class TestDevicesInventoryShow:
    def test_ok(self, device_with_inventory):
        device_id, _ = device_with_inventory

        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "show",
            device_id,
        )
        assert r.returncode == 0, r.stderr
        expect_output(
            r.stdout, f"ID: {device_id}", "inventory/device_type: test-device"
        )

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "show",
            device_id,
            "--raw",
        )
        assert r.returncode == 0, r.stderr
        raw = json.loads(r.stdout, strict=False)
        assert raw["id"] == device_id

    def test_error_no_login(self):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "show",
            "some_device_id",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")


class TestDevicesInventoryAttributes:
    def test_ok(self, device_with_inventory):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "attributes",
            "--output",
            "csv",
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stdout, "inventory,device_type,1\n")

    def test_error_no_login(self):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "devices",
            "inventory",
            "attributes",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")
