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
import redo
import subprocess

from contextlib import contextmanager
from multiprocessing import Process
from tempfile import NamedTemporaryFile

import filelock

from DNS import DnsRequest, SocketError

from common_connect import wait_for_connect
from helpers import md5sum, get_device_id

# Function must be defined outside of class so it can be pickled
def port_forward(auth_token, server_url, dev_id, port_mapping, *port_mappings):
    return subprocess.run(
        [
            "mender-cli",
            "--skip-verify",
            "--server",
            server_url,
            "--token-value",
            auth_token,
            "port-forward",
            dev_id,
            port_mapping,
        ]
        + list(port_mappings),
        check=True,
        capture_output=True,
    )


@contextmanager
def port_forward_process(auth_token, server_url, dev_id, port_mapping, *port_mappings):
    pfw = Process(
        target=port_forward,
        args=(auth_token, server_url, dev_id, port_mapping) + port_mappings,
    )
    try:
        pfw.start()
        yield pfw
    finally:
        if pfw.is_alive():
            pfw.terminate()
            pfw.join(timeout=5)
            if pfw.is_alive():
                pfw.kill()


@redo.retriable(sleeptime=10, attempts=6)
def dns_request(name, qtype, server, port):
    req = DnsRequest(name=name, qtype=qtype, server=server, port=port)
    try:
        response = req.req()
        assert len(response.answers) >= 1, response.show()
    except SocketError as err:
        raise err


@redo.retriable(sleeptime=10, attempts=6)
def run_scp(command):
    proc = subprocess.run(command, check=True, capture_output=True)
    assert proc.returncode == 0, (proc.stdout, proc.stderr)


class BaseTestPortForward:
    """Tests the port forward functionality"""

    def do_test_portforward(self, env, devid):
        # wait for the device to connect via websocket
        wait_for_connect(env, devid)

        server_url = "https://" + env.server.host
        auth_token = env.server.auth_token

        # Acquire lock to avoid enterprise and open-source to bind the same port
        with filelock.FileLock(".test_portforward_lock"):
            tcp_port = 9922
            udp_port = 9953
            with port_forward_process(
                auth_token,
                server_url,
                devid,
                f"{tcp_port}:22",
                f"udp/{udp_port}:8.8.8.8:53",
            ):

                # verify the UDP port-forward querying the Google's DNS server
                dns_request(
                    name="mender.io", qtype="A", server="localhost", port=udp_port
                )
                dns_request(
                    name="mender.io", qtype="MX", server="localhost", port=udp_port
                )

                # verify the TCP port-forward using scp to upload and download files
                try:
                    # create a 1KB random file
                    f = NamedTemporaryFile(delete=False)
                    f.write(os.urandom(1024))
                    f.close()

                    # upload the file using scp
                    run_scp(
                        [
                            "scp",
                            "-O",
                            "-o",
                            "StrictHostKeyChecking=no",
                            "-o",
                            "UserKnownHostsFile=/dev/null",
                            "-P",
                            str(tcp_port),
                            f.name,
                            f"root@localhost:{f.name}",
                        ]
                    )

                    # download the file using scp
                    run_scp(
                        [
                            "scp",
                            "-O",
                            "-o",
                            "StrictHostKeyChecking=no",
                            "-o",
                            "UserKnownHostsFile=/dev/null",
                            "-P",
                            str(tcp_port),
                            f"root@localhost:{f.name}",
                            f.name + ".download",
                        ]
                    )
                    # assert the files are not corrupted
                    assert md5sum(f.name) == md5sum(f.name + ".download")
                finally:
                    os.unlink(f.name)
                    if os.path.isfile(f.name + ".download"):
                        os.unlink(f.name + ".download")


class TestPortForwardOpenSource(BaseTestPortForward):
    def test_portforward(self, standard_setup_one_client_bootstrapped):
        devid = get_device_id(
            standard_setup_one_client_bootstrapped.device,
            standard_setup_one_client_bootstrapped.server,
        )

        self.do_test_portforward(standard_setup_one_client_bootstrapped, devid)
