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

import socket

import pytest
import boto3
import docker


@pytest.fixture(scope="session")
def s3_client():
    return boto3.resource("s3")


@pytest.fixture(scope="function")
def clean_s3(s3_client, bucket="mender"):
    bucket = s3_client.Bucket(bucket)
    for item in bucket.objects.all():
        item.delete()


@pytest.fixture(scope="session")
def docker_client():
    return docker.from_env()


@pytest.fixture(scope="session")
def docker_compose_project(docker_client: docker.DockerClient):
    return docker_client.containers.list(filters={"id": socket.gethostname()})[
        0
    ].labels.get("com.docker.compose.project")


@pytest.fixture(scope="session")
def useradm_container(docker_client: docker.DockerClient, docker_compose_project: str):
    return docker_client.containers.list(
        limit=1,
        all=True,
        filters={
            "label": [
                f"com.docker.compose.project={docker_compose_project}",
                "com.docker.compose.service=useradm",
            ]
        },
    )[0]


@pytest.fixture(scope="session")
def mongo_container(docker_client: docker.DockerClient, docker_compose_project: str):
    return docker_client.containers.list(
        limit=1,
        all=True,
        filters={
            "label": [
                f"com.docker.compose.project={docker_compose_project}",
                "com.docker.compose.service=mongo",
            ]
        },
    )[0]


@pytest.fixture(scope="function")
def single_user(useradm_container):
    exit_code, output = useradm_container.exec_run(
        cmd=[
            "/usr/bin/useradm",
            "create-user",
            "--username",
            "user@tenant.com",
            "--password",
            "youcantguess",
        ],
    )

    assert exit_code == 0, f"output: {output}"


_SNIPPET_CLEAN_DB = """
db.adminCommand({
    listDatabases: 1,
    nameOnly: true,
    filter: {
        name: {
            $nin: ["admin", "local", "config"]
        }
    },
}).databases.forEach(cmdRes => {
    db.getSiblingDB(cmdRes.name).dropDatabase();
});
"""


@pytest.fixture(scope="function", autouse=True)
def clean_db(mongo_container):
    exit_code, output = mongo_container.exec_run(
        cmd=[
            "mongosh",
            "--eval",
            _SNIPPET_CLEAN_DB,
        ]
    )

    assert exit_code == 0, f"output: {output}"
