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

import requests


class Deployments:
    def __init__(self, token, url='https://mender-api-gateway/api/management/v1/deployments'):
        self.url=url
        self.token=token

    def get_artifacts(self):
        auth = {'Authorization': 'Bearer {}'.format(self.token)}
        return requests.get(self.make_api_url('/artifacts'), verify=False, headers=auth)

    def make_api_url(self, path):
        return os.path.join(self.url,
                            path if not path.startswith("/") else path[1:])
