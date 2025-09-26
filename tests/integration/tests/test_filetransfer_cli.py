# Copyright 2025 Northern.tech AS
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
#

import os
import subprocess

from tempfile import NamedTemporaryFile

from common_connect import wait_for_connect
from helpers import md5sum, get_device_id


class BaseTestFileTransferCLI:
    """Tests the file transfer functionality"""

    def do_test_filetransfer_cli(self, env, devid):
        # wait for the device to connect via websocket
        wait_for_connect(env, devid)

        # authenticate with mender-cli
        server_url = "https://" + env.server.host
        username = env.server.username
        password = env.server.password
        p = subprocess.Popen(
            [
                "mender-cli",
                "--skip-verify",
                "--server",
                server_url,
                "login",
                "--username",
                username,
                "--password",
                password,
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        stdout, stderr = p.communicate()
        exit_code = p.wait()
        assert exit_code == 0, (stdout, stderr)

        # upload and download files using mender-cli
        try:
            # create a 40MB random file
            f = NamedTemporaryFile(delete=False)
            for _ in range(40 * 1024):
                f.write(os.urandom(1024))
            f.close()

            # upload the file using mender-cli
            p = subprocess.Popen(
                [
                    "mender-cli",
                    "--skip-verify",
                    "--server",
                    server_url,
                    "cp",
                    f.name,
                    devid + ":/tmp/random.bin",
                ],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            stdout, stderr = p.communicate()
            exit_code = p.wait()
            assert exit_code == 0, (stdout, stderr)

            # download the file using mender-cli
            p = subprocess.Popen(
                [
                    "mender-cli",
                    "--skip-verify",
                    "--server",
                    server_url,
                    "cp",
                    devid + ":/tmp/random.bin",
                    f.name + ".download",
                ],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            stdout, stderr = p.communicate()
            exit_code = p.wait()
            assert exit_code == 0, (stdout, stderr)

            # assert the files are not corrupted
            assert md5sum(f.name) == md5sum(f.name + ".download")

            # upload the file to a directory (fail)
            p = subprocess.Popen(
                [
                    "mender-cli",
                    "--skip-verify",
                    "--server",
                    server_url,
                    "cp",
                    f.name,
                    devid + ":/tmp/path-does-not-exist/random.bin",
                ],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            stdout, stderr = p.communicate()
            exit_code = p.wait()
            assert exit_code == 1, (stdout, stderr)
            assert b"failed to create target file" in stderr, (stdout, stderr)

        finally:
            os.unlink(f.name)
            if os.path.isfile(f.name + ".download"):
                os.unlink(f.name + ".download")

        # download a file which doesn't exist (fail)
        p = subprocess.Popen(
            [
                "mender-cli",
                "--skip-verify",
                "--server",
                server_url,
                "cp",
                "/this/file/does/not/exist",
                devid + ":/tmp/path-does-not-exist/random.bin",
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        stdout, stderr = p.communicate()
        exit_code = p.wait()
        assert exit_code == 1, (stdout, stderr)
        assert b"no such file or directory" in stderr, (stdout, stderr)

        # upload a file which doesn't exist (fail)
        p = subprocess.Popen(
            [
                "mender-cli",
                "--skip-verify",
                "--server",
                server_url,
                "cp",
                "/this/file/does/not/exist",
                devid + ":/tmp/test.bin",
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        stdout, stderr = p.communicate()
        exit_code = p.wait()
        assert exit_code == 1, (stdout, stderr)
        assert b"no such file or directory" in stderr, (stdout, stderr)


class TestFileTransferCLI(BaseTestFileTransferCLI):
    """Tests the file transfer functionality"""

    def test_filetransfer_cli(self, standard_setup_one_client_bootstrapped):
        devid = get_device_id(
            standard_setup_one_client_bootstrapped.device,
            standard_setup_one_client_bootstrapped.server,
        )

        self.do_test_filetransfer_cli(standard_setup_one_client_bootstrapped, devid)
