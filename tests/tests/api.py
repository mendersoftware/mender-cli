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
import json

import requests

import crypto


def make_api_url(url, path):
    return os.path.join(url, path if not path.startswith("/") else path[1:])


class Deployments:
    def __init__(
        self, token, url="https://mender-api-gateway/api/management/v1/deployments"
    ):
        self.url = url
        self.token = token

    def get_artifacts(self):
        auth = {"Authorization": "Bearer {}".format(self.token)}
        return requests.get(
            make_api_url(self.url, "/artifacts"), verify=False, headers=auth
        )

    def get_releases(self, name):
        auth = {"Authorization": "Bearer {}".format(self.token)}
        return requests.get(
            make_api_url(self.url, "/deployments/releases?name=" + name),
            verify=False,
            headers=auth,
        )


class DevicesDevauth:
    def __init__(self, url="https://mender-api-gateway/api/devices/v1/authentication"):
        self.url = url

    def auth_req(self, id_data, pubkey, privkey):
        payload = {
            "id_data": json.dumps(id_data),
            "pubkey": pubkey,
        }
        signature = crypto.auth_req_sign_rsa(json.dumps(payload), privkey)
        return payload, {"X-MEN-Signature": signature}

    def post_auth_request(self, authset, pubkey, privkey):
        body, sighdr = self.auth_req(authset, pubkey, privkey)
        return requests.post(
            make_api_url(self.url, "/auth_requests"),
            verify=False,
            headers=sighdr,
            json=body,
        )


class ManagementDevauth:
    def __init__(
        self, token, url="https://mender-api-gateway/api/management/v2/devauth"
    ):
        self.url = url
        self.token = token

    def get_devices(self):
        auth = {"Authorization": "Bearer {}".format(self.token)}
        return requests.get(
            make_api_url(self.url, "/devices"), verify=False, headers=auth
        )

    def set_device_auth_status(self, device_id, authset_id, auth_status):
        body = {
            "status": auth_status,
        }
        auth = {"Authorization": "Bearer {}".format(self.token)}
        return requests.put(
            make_api_url(self.url, f"/devices/{device_id}/auth/{authset_id}/status"),
            verify=False,
            headers=auth,
            json=body,
        )
