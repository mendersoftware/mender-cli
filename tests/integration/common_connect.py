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

import redo

import testutils.api.deviceconnect as deviceconnect
from testutils.api.client import ApiClient


def wait_for_connect(env, devid):
    devconn = ApiClient(host=env.server.host, base_url=deviceconnect.URL_MGMT,)

    connected = 0
    for _ in redo.retrier(attempts=12, sleeptime=5):
        res = devconn.call(
            "GET",
            deviceconnect.URL_MGMT_DEVICE,
            headers={"Authorization": f"Bearer {env.server.auth_token}"},
            path_params={"id": devid},
        )
        if not (res.status_code == 200 and res.json()["status"] == "connected"):
            connected = 0
            continue
        connected += 1
        if connected >= 2:
            break
    else:
        assert False, "timed out waiting for /connect"
