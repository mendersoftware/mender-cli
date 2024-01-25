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
from common import single_user, expect_output, DEFAULT_TOKEN_PATH
import api
import docker
import s3


@pytest.fixture(scope="function")
def logged_in_single_user(single_user):
    c = cli.Cli()
    r = c.run('login', \
              '--server', 'https://mender-api-gateway', \
              '--skip-verify', \
              '--username', 'user@tenant.com', \
              '--password', 'youcantguess')

    assert r.returncode == 0, r.stderr
    yield
    os.remove(DEFAULT_TOKEN_PATH)


@pytest.fixture(scope="function")
def valid_artifact():
    path = '/tests/foo-artifact'
    artifact.create_artifact_file(path)
    yield path
    os.remove(path)

@pytest.fixture(scope="function")
def valid_artifact_new():
    path = '/tests/foo-artifact-new'
    artifact.create_artifact_file(path)
    yield path
    os.remove(path)

@pytest.fixture(scope="function")
def clean_deployments_db():
    yield
    r = docker.exec('mender-mongo', \
                    docker.BASE_COMPOSE_FILES, \
                    'mongosh', 'deployment_service', '--eval', 'db.dropDatabase()')
    assert r.returncode == 0, r.stderr

@pytest.fixture(scope="function")
def clean_mender_storage():
    yield
    s3.cleanup_mender_storage()

class TestArtifactUpload:
    @pytest.mark.usefixtures('clean_deployments_db', 'clean_mender_storage')
    def test_ok_direct(self, logged_in_single_user, valid_artifact_new):
        c = cli.MenderCliCoverage()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'upload', \
                  '--direct',
                  valid_artifact_new)

        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'upload successful')

    @pytest.mark.usefixtures('clean_deployments_db', 'clean_mender_storage')
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.MenderCliCoverage()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'upload', \
                  '--description', 'foo',
                  valid_artifact)

        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'upload successful')

        token = Path(DEFAULT_TOKEN_PATH).read_text()

        dapi = api.Deployments(token)
        r = dapi.get_artifacts()

        assert r.status_code == 200

        artifacts = r.json()
        assert len(artifacts) == 1

        artifact = artifacts[0]

        assert artifact['name'] == 'artifact-foo'
        assert artifact['description'] == 'foo'
        assert artifact['device_types_compatible'] == ['device-foo']

    def test_error_no_login(self, valid_artifact):
        c = cli.MenderCliCoverage()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'upload', \
                  '--description', 'foo',
                  valid_artifact)

        assert r.returncode!=0
        expect_output(r.stderr, 'FAILURE', 'Please Login first')


class TestArtifactDelete:
    @pytest.mark.usefixtures('clean_deployments_db', 'clean_mender_storage')
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.MenderCliCoverage()

        # upload artifact first
        r = c.run('--server', 'https://mender-api-gateway', \
            '--skip-verify', \
            'artifacts', 'upload', \
            '--description', 'foo',
            valid_artifact)

        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'upload successful')
        
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'list')

        regex = re.compile(r"[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}\Z", re.I)
        artifact_id = None
        for line in r.stdout.split("\n"):
            match = regex.search(line)
            if match:
                artifact_id = match.group()
                break

        assert r.returncode==0, r.stderr
        assert artifact_id is not None
        
        r = c.run('--server', 'https://mender-api-gateway', \
            '--skip-verify', \
            'artifacts', 'delete',
            artifact_id)
        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'delete successful')

        token = Path(DEFAULT_TOKEN_PATH).read_text()

        dapi = api.Deployments(token)
        r = dapi.get_artifacts()

        assert r.status_code == 200
        artifacts = r.json()
        assert len(artifacts) == 0

    def test_error_no_login(self):
        c = cli.MenderCliCoverage()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'delete', \
                  "some_artifact_id")

        assert r.returncode!=0
        expect_output(r.stderr, 'FAILURE', 'Please Login first')

class TestArtifactDownload:
    @pytest.mark.usefixtures('clean_deployments_db', 'clean_mender_storage')
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.MenderCliCoverage()

        # upload artifact first
        r = c.run('--server', 'https://mender-api-gateway', \
            '--skip-verify', \
            'artifacts', 'upload', \
            '--description', 'foo',
            valid_artifact)

        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'upload successful')

        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'list')

        regex = re.compile(r"[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}\Z", re.I)
        artifact_id = None
        for line in r.stdout.split("\n"):
            match = regex.search(line)
            if match:
                artifact_id = match.group()
                break

        assert r.returncode==0, r.stderr
        assert artifact_id is not None

        r = c.run('--server', 'https://mender-api-gateway', \
            '--skip-verify', \
            'artifacts', 'download',
            artifact_id)
        assert r.returncode==0, r.stderr
        expect_output(r.stdout, 'download successful')

    def test_error_no_login(self):
        c = cli.MenderCliCoverage()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'download', \
                  "some_artifact_id")

        assert r.returncode!=0
        expect_output(r.stderr, 'FAILURE', 'Please Login first')
