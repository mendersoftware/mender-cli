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


@pytest.fixture(scope="function")
def cleanup_token(request):
    yield
    os.remove(request.param)


class TestLogin:
    @pytest.mark.parametrize("cleanup_token", [DEFAULT_TOKEN_PATH], indirect=True)
    def test_ok(self, single_user, cleanup_token):
        c = cli.MenderCliCoverage()

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

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        expect_output(r.stderr, "login successful")

    @pytest.mark.parametrize("cleanup_token", ["/tests/authtoken"], indirect=True)
    def test_ok_custom_path(self, single_user, cleanup_token):
        c = cli.MenderCliCoverage()

        custom_path = "/tests/authtoken"

        r = c.run(
            "login",
            "--server",
            "https://mender-api-gateway",
            "--skip-verify",
            "--token",
            "/tests/authtoken",
            "--username",
            "user@tenant.com",
            "--password",
            "youcantguess",
        )

        assert r.returncode == 0, r.stderr

        self.__check_token_at(custom_path)
        expect_output(r.stderr, "login successful")

    @pytest.mark.parametrize("cleanup_token", [DEFAULT_TOKEN_PATH], indirect=True)
    def test_ok_verbose(self, single_user, cleanup_token):
        c = cli.MenderCliCoverage()

        r = c.run(
            "login",
            "--server",
            "https://mender-api-gateway",
            "--skip-verify",
            "--verbose",
            "--username",
            "user@tenant.com",
            "--password",
            "youcantguess",
        )

        assert r.returncode == 0, r.stderr

        self.__check_token_at(DEFAULT_TOKEN_PATH)
        expect_output(
            r.stderr, "creating directory", "saved token to", "login successful"
        )

    def test_error_wrong_creds(self, single_user):
        c = cli.MenderCliCoverage()

        r = c.run(
            "login",
            "--server",
            "https://mender-api-gateway",
            "--skip-verify",
            "--username",
            "notfound@tenant.com",
            "--password",
            "youcantguess",
        )

        assert r.returncode != 0

        expect_output(r.stderr, "FAILURE: login failed with status 401")

    def _write_mender_cli_conf(self, conf):
        """Simple utility function to write the configuration file"""
        with open(os.getenv("HOME") + "/.mender-clirc", "w") as f:
            try:
                f.write(conf)
            except Exception as e:
                pytest.fail("Failed to create configuration file: {}".format(e))

    def test_login_from_configuration_file(self, single_user):
        """test that the configuration file parameters are respected"""
        # Wrong username and password
        conf = """
        {
            "username": "foobar",
            "password": "barbaz",
            "server": "https://mender-api-gateway"
        }
        """

        self._write_mender_cli_conf(conf)

        c = cli.MenderCliCoverage()
        r = c.run("login", "--skip-verify")

        assert r.returncode != 0, r.stderr

        expect_output(r.stderr, "FAILURE: login failed with status 401")

        # correct username and password, wrong server
        conf = """
        {
            "username": "user@tenant.com",
            "password": "youcantguess",
            "server": "https://wrong.server.com"
        }
        """

        self._write_mender_cli_conf(conf)

        r = c.run("login", "--skip-verify")

        assert r.returncode != 0, r.stderr

        expect_output(r.stderr, "FAILURE:", "request failed")

        # correct username, password and server
        conf = """
        {
            "username": "user@tenant.com",
            "password": "youcantguess",
            "server": "https://mender-api-gateway"
        }
        """

        self._write_mender_cli_conf(conf)

        r = c.run("login", "--skip-verify")

        assert r.returncode == 0, r.stderr

    def test_login_from_stdin(self, single_user):
        """test login reading password from stdin"""
        conf = """
        {
            "username": "user@tenant.com",
            "server": "https://mender-api-gateway"
        }
        """

        self._write_mender_cli_conf(conf)

        c = cli.MenderCliCoverage()
        r = c.run_and_enter_password("login", "--skip-verify", password="youcantguess")

        assert r.returncode == 0, r.stderr

    def test_configuration_parameter_override(self, single_user):
        """test that parameters listed in the configuration file can be overridden on the CLI

        All parameters in the configuration are wrong, and should therefore be
        overridden by the CLI parameters

        """

        conf = """
        {
            "username": "foobar",
            "password": "barbaz",
            "server": "https://wrong.server.com"
        }
        """

        self._write_mender_cli_conf(conf)

        c = cli.MenderCliCoverage()
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

    def __check_token_at(self, path):
        assert os.path.isfile(path)
