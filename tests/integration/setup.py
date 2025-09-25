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

import time
import random
import logging
import subprocess

from helpers import docker_lock

from testutils.infra.device import MenderDevice

logger = logging.getLogger(__name__)


project_name = f"mender"
project_name_client = f"virtual_device{random.randint(1000, 9999)}"


class Env:
    def __init__(self):
        self.server = None

        self.devices = []
        self.gateway = None

        self.auth = None

    def get_virtual_network_host_ip(self):
        container = f"{project_name_client}-client-1"
        cmd = [
            "docker",
            "inspect",
            "-f",
            "{{range .NetworkSettings.Networks}}{{.Gateway}}{{end}}",
            container,
        ]
        with docker_lock:
            output = subprocess.check_output(cmd)
        return output.decode().strip()


def clients_up(number_of_clients):
    with docker_lock:
        subprocess.check_call(
            [
                "docker",
                "compose",
                "-p",
                project_name_client,
                "-f",
                "docker-compose.client.yml",
                "up",
                "-d",
                "--scale",
                f"client={number_of_clients}",
            ]
        )
        clients = []
        for device_number in range(1, number_of_clients + 1):
            device = f"{project_name_client}-client-{device_number}"
            client_ip = subprocess.check_output(
                [
                    "docker",
                    "inspect",
                    "-f",
                    "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
                    device,
                ],
                text=True,
            ).strip()
            clients.append(MenderDevice(f"{client_ip}:8822"))
    return clients


def wait_for_devices(env):
    # Give the device some time to start so we don't get
    # stuck in a 60 second ssh exception right away
    time.sleep(15)
    for device in env.devices:
        logger.info(f"Waiting for '{device.host}'")
        device.ssh_is_opened()
