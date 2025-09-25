#!/usr/bin/python
# Copyright 2023 Northern.tech AS
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
import re
from pathlib import Path

import pytest

import cli
import artifact
from common import expect_output, DEFAULT_TOKEN_PATH
import api
import docker


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


@pytest.fixture(scope="function")
def valid_artifact():
    yield "data/test.mender"


class TestArtifactUpload:
    @pytest.mark.usefixtures("clean_s3")
    def test_ok_direct(self, logged_in_single_user, valid_artifact):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "upload",
            "--direct",
            valid_artifact,
        )

        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "upload successful")

    @pytest.mark.usefixtures("clean_s3")
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "upload",
            "--description",
            "foo",
            valid_artifact,
        )

        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "upload successful")

        token = Path(DEFAULT_TOKEN_PATH).read_text()

        dapi = api.Deployments(token)
        r = dapi.get_artifacts()

        assert r.status_code == 200

        artifacts = r.json()
        assert len(artifacts) == 1

        artifact = artifacts[0]

        assert artifact["name"] == "test"
        assert artifact["description"] == "foo"
        assert artifact["device_types_compatible"] == ["test"]

    def test_error_no_login(self, valid_artifact):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "upload",
            "--description",
            "foo",
            valid_artifact,
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")


class TestArtifactDelete:
    @pytest.mark.usefixtures("clean_s3")
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.Cli()

        # upload artifact first
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "upload",
            "--description",
            "foo",
            valid_artifact,
        )

        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "upload successful")

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "list",
        )

        regex = re.compile(
            r"[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}\Z",
            re.I,
        )
        artifact_id = None
        for line in r.stdout.split("\n"):
            match = regex.search(line)
            if match:
                artifact_id = match.group()
                break

        assert r.returncode == 0, r.stderr
        assert artifact_id is not None

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "delete",
            artifact_id,
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "delete successful")

        token = Path(DEFAULT_TOKEN_PATH).read_text()

        dapi = api.Deployments(token)
        r = dapi.get_artifacts()

        assert r.status_code == 200
        artifacts = r.json()
        assert len(artifacts) == 0

    def test_error_no_login(self):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "delete",
            "some_artifact_id",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")


class TestArtifactDownload:
    @pytest.mark.usefixtures("clean_s3")
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.Cli()

        # upload artifact first
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "upload",
            "--description",
            "foo",
            valid_artifact,
        )

        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "upload successful")

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "list",
        )

        regex = re.compile(
            r"[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}\Z",
            re.I,
        )
        artifact_id = None
        for line in r.stdout.split("\n"):
            match = regex.search(line)
            if match:
                artifact_id = match.group()
                break

        assert r.returncode == 0, r.stderr
        assert artifact_id is not None

        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "download",
            artifact_id,
        )
        assert r.returncode == 0, r.stderr
        expect_output(r.stderr, "download successful")

    def test_error_no_login(self):
        c = cli.Cli()
        r = c.run(
            "--server",
            "https://docker.mender.io",
            "--skip-verify",
            "artifacts",
            "download",
            "some_artifact_id",
        )

        assert r.returncode != 0
        expect_output(r.stderr, "FAILURE", "Please Login first")
