#!/usr/bin/python
# Copyright 2018 Mender Software AS
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        https://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
from pathlib import Path
import os

import pytest

from common import single_user
import cli

USER_HOME = str(Path.home())
DEFAULT_TOKEN_PATH = os.path.join(USER_HOME,'.mender', 'authtoken')


class TestLogin:
    def test_ok(self, single_user):
        c = cli.Cli()

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode == 0, r.stderr

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        self.__expect_output(r.stdout, \
                             'login successful')

    def test_ok_custom_path(self, single_user):
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
        self.__expect_output(r.stdout, \
                             'login successful')

    def test_ok_verbose(self, single_user):
        c = cli.Cli()

        r = c.run('login', \
                  '--server', 'https://mender-api-gateway', \
                  '--skip-verify', \
                  '--verbose', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode == 0, r.stderr

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        self.__expect_output(r.stdout, \
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

        self.__expect_output(r.stderr, 'FAILURE: login failed with status 401')

    def test_error_no_server(self, single_user):
        c = cli.Cli()

        r = c.run('login', \
                  '--skip-verify', \
                  '--username', 'user@tenant.com', \
                  '--password', 'youcantguess')

        assert r.returncode != 0

        self.__expect_output(r.stderr, '"server" not set')

    def __check_token_at(self, path):
        assert os.path.isfile(path)

    def __expect_output(self, stream, *expected):
        for e in expected:
            assert e in stream, 'expected string {} not found in stream'.format(e)
