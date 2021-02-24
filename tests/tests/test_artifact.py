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
import artifact
from common import single_user, expect_output, DEFAULT_TOKEN_PATH
import api
import docker
import s3


@pytest.yield_fixture(scope="function")
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


@pytest.yield_fixture(scope="function")
def valid_artifact():
    path = '/tests/foo-artifact'
    artifact.create_artifact_file(path)
    yield path
    os.remove(path)

@pytest.yield_fixture(scope="function")
def clean_deployments_db():
    yield
    r = docker.exec('mender-mongo', \
                    docker.BASE_COMPOSE_FILES, \
                    'mongo', 'deployment_service', '--eval', 'db.dropDatabase()')
    assert r.returncode == 0, r.stderr

@pytest.yield_fixture(scope="function")
def clean_mender_storage():
    yield
    s3.cleanup_mender_storage()

class TestArtifactUpload:
    @pytest.mark.usefixtures('clean_deployments_db', 'clean_mender_storage')
    def test_ok(self, logged_in_single_user, valid_artifact):
        c = cli.Cli()
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
        c = cli.Cli()
        r = c.run('--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  'artifacts', 'upload', \
                  '--description', 'foo',
                  valid_artifact)

        assert r.returncode!=0
        expect_output(r.stderr, 'FAILURE', 'Please Login first')
