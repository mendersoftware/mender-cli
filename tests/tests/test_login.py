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

import pytest

from common import single_user, expect_output, DEFAULT_TOKEN_PATH
import cli

@pytest.yield_fixture(scope="function")
def cleanup_token(request):
    yield
    os.remove(request.param)

class TestLogin:
    @pytest.mark.parametrize('cleanup_token', [DEFAULT_TOKEN_PATH], indirect=True)
    def test_ok(self, single_user, cleanup_token):
        c = cli.Cli()

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode == 0, r.stderr

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        expect_output(r.stdout, \
                             'login successful')

    @pytest.mark.parametrize('cleanup_token', ['/tests/authtoken'], indirect=True)
    def test_ok_custom_path(self, single_user, cleanup_token):
        c = cli.Cli()

        custom_path = '/tests/authtoken'

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--token', '/tests/authtoken', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode == 0, r.stderr

        self.__check_token_at(custom_path)
        expect_output(r.stdout, \
                             'login successful')

    @pytest.mark.parametrize('cleanup_token', [DEFAULT_TOKEN_PATH], indirect=True)
    def test_ok_verbose(self, single_user, cleanup_token):
        c = cli.Cli()

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--verbose', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode == 0, r.stderr

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        expect_output(r.stdout, \
                             'creating directory',
                             'saved token to',
                             'login successful')

    def test_error_wrong_creds(self, single_user):
        c = cli.Cli()

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--username', 'notfound@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode != 0

        expect_output(r.stderr, 'FAILURE: login failed with status 401')

    def test_error_no_server(self, single_user):
        c = cli.Cli()

        r = c.run('login', \
                  '--skip-verify', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode != 0

        expect_output(r.stderr, '"server" not set')

    def __check_token_at(self, path):
        assert os.path.isfile(path)
