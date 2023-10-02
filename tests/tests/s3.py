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
from minio import Minio


def cleanup_mender_storage(bucket="mender-artifact-storage"):
    client = Minio(
        "minio.s3.docker.mender.io:9000",
        access_key="minio",
        secret_key="minio123",
        region="us-east-1",
        secure=False,
    )

    objs = client.list_objects(bucket)
    client.remove_objects(bucket, objs)
