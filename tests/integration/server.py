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
import uuid
import subprocess
import redo
import logging

logger = logging.getLogger(__name__)

from testutils.common import create_user
from testutils.api.client import ApiClient
from testutils.api import (
    deployments,
    useradm,
    tenantadm,
    deviceauth,
)

from setup import project_name
from helpers import get_mac_address, docker_lock


def get_server_host(container=f"{project_name}-traefik-1"):
    with docker_lock:
        return (
            subprocess.check_output(
                "docker inspect "
                + container
                + "  --format='{{range .NetworkSettings.Networks}}{{.Gateway}}{{end}}'",
                shell=True,
            )
            .decode()
            .strip()
        )


class Server:
    def __init__(self, plan="os"):
        self.containers_namespace = "mender"
        self.host = get_server_host()
        self.deployment_id = ""

        self.api_dev_deploy = ApiClient(deployments.URL_MGMT, self.host)
        self.devauthm = ApiClient(deviceauth.URL_MGMT, host=self.host)
        self.tenantadm = ApiClient(tenantadm.URL_MGMT, host=self.host)

        self.get_tenant_username_and_password(plan=plan)
        self.auth_token = self.get_auth_token()

        self.num_accepted_devices = 0

    def get_auth_token(self):
        r = ApiClient(useradm.URL_MGMT, self.host).call(
            "POST", useradm.URL_LOGIN, auth=(self.username, self.password)
        )
        assert r.status_code == 200
        assert r.text is not None
        assert r.text != ""

        return r.text

    def get_pending_devices(self):
        qs_params = {"page": 1, "per_page": 64, "status": "pending"}
        result = self.devauthm.with_auth(self.auth_token).call(
            "GET", deviceauth.URL_MGMT_DEVICES, qs_params=qs_params
        )
        assert result.status_code == 200
        return result.json()

    def get_accepted_devices(self):
        qs_params = {"page": 1, "per_page": 64, "status": "accepted"}
        r = self.devauthm.with_auth(self.auth_token).call(
            "GET", deviceauth.URL_MGMT_DEVICES, qs_params=qs_params
        )
        assert r.status_code == 200
        return r.json()

    def get_tenant_token(self):
        r = self.tenantadm.with_auth(self.auth_token).call(
            "GET", tenantadm.URL_MGMT_THIS_TENANT
        )
        assert r.status_code == 200
        return r.json()["tenant_token"]

    def accept_devices(self, devices, max_sleeptime=60, max_attempts=10):
        def _accept_devices():
            device_ids = []
            mac_addresses = [get_mac_address(device) for device in devices]
            for device in self.get_pending_devices():
                if device["identity_data"]["mac"] in mac_addresses:
                    device_id = device["id"]
                    auth_set_id = device["auth_sets"][0]["id"]
                    r = self.devauthm.with_auth(self.auth_token).call(
                        "PUT",
                        deviceauth.URL_AUTHSET_STATUS,
                        deviceauth.req_status("accepted"),
                        path_params={"did": device_id, "aid": auth_set_id},
                    )
                    assert r.status_code == 204
                    self.num_accepted_devices += 1
                    device_ids.append(device_id)

            if self.num_accepted_devices != len(devices):
                logging.info(
                    f"Accepted devices: {self.num_accepted_devices}, expected: {len(devices)}"
                )
                raise ValueError

            return device_ids

        return redo.retry(
            _accept_devices,
            max_sleeptime=max_sleeptime,
            sleeptime=5,
            sleepscale=2,
            attempts=max_attempts,
        )

    def create_deployment(self, artifact_name, device_ids):
        logger.info("Creating deployment")
        response = self.api_dev_deploy.with_auth(self.auth_token).call(
            "POST",
            deployments.URL_DEPLOYMENTS,
            body={
                "name": artifact_name,
                "artifact_name": artifact_name,
                "devices": device_ids,
            },
        )
        assert response.status_code == 201, f"{response.text} {response.status_code}"
        self.deployment_id = os.path.basename(response.headers["Location"])
        return self.deployment_id

    def get_tenant_username_and_password(self, plan):
        _ = plan
        uuidv4 = str(uuid.uuid4())
        self.username, self.password = (
            "some.user+" + uuidv4 + "@example.com",
            "secretsecret",
        )
        create_user(
            self.username,
            self.password,
            containers_namespace=self.containers_namespace,
        )
        return None, self.username, self.password

    def upload_image(self, filename):
        self.api_dev_deploy.headers = {}
        response = self.api_dev_deploy.with_auth(self.auth_token).call(
            "POST",
            deployments.URL_DEPLOYMENTS_ARTIFACTS,
            files=(
                ("description", (None)),
                ("size", (None, str(os.path.getsize(filename)))),
                (
                    "artifact",
                    (filename, open(filename, "rb"), "application/octet-stream"),
                ),
            ),
        )
        assert response.status_code == 201, f"{response.text} {response.status_code}"

    def check_expected_statistics(
        self,
        deployment_id,
        expected_status,
        expected_mender_clients,
        max_sleeptime=60,
        max_attempts=10,
    ):
        def _check_expected_statistics():
            response = self.api_dev_deploy.with_auth(self.auth_token).call(
                "GET", deployments.URL_DEPLOYMENTS_ID.format(id=deployment_id)
            )
            assert expected_mender_clients == int(
                response.json()["statistics"]["status"][expected_status]
            )

        return redo.retry(
            _check_expected_statistics,
            max_sleeptime=max_sleeptime,
            sleeptime=5,
            sleepscale=2,
            attempts=max_attempts,
        )

    def check_expected_status(
        self, expected_status, deployment_id, max_sleeptime=60, max_attempts=10
    ):
        def _check_expected_status():
            response = self.api_dev_deploy.with_auth(self.auth_token).call(
                "GET", deployments.URL_DEPLOYMENTS_ID.format(id=deployment_id)
            )
            assert (
                response.status_code == 200
            ), f"{response.text} {response.status_code}"
            assert response.json()["status"] == expected_status

        return redo.retry(
            _check_expected_status,
            max_sleeptime=max_sleeptime,
            sleeptime=5,
            sleepscale=2,
            attempts=max_attempts,
        )

    def get_deployment_logs(self, device_id, deployment_id):
        deployment_logs = f"/deployments/{deployment_id}/devices/{device_id}/log"
        response = self.api_dev_deploy.with_auth(self.auth_token).call(
            "GET", deployment_logs
        )
        assert response.status_code == 200, f"{response.text} {response.status_code}"
        return response.text

    def abort_deployment(self, deployment_id):
        set_status = f"/deployments/{deployment_id}/status"
        response = self.api_dev_deploy.with_auth(self.auth_token).call(
            "PUT", set_status, body={"status": "aborted"},
        )
        assert response.status_code == 204, f"{r.text} {r.status_code}"
