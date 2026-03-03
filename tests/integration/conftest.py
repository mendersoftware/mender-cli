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

from os import path
import sys

import xdist
import logging
import pytest

sys.path += [path.join(path.dirname(__file__), "../mender_server/backend/tests/")]

from server import Server
from setup import (
    Env,
    clients_up,
    wait_for_devices,
    project_name,
    project_name_client,
)

import helpers

logging.getLogger("requests").setLevel(logging.CRITICAL)
logging.getLogger("urllib3").setLevel(logging.CRITICAL)
logging.getLogger("redo").setLevel(logging.INFO)


def pytest_addoption(parser):
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )


def stop_open_source_server():
    helpers.docker_compose_stop(
        project_name=project_name,
        files=["../mender_server/docker-compose.yml"],
    )


def start_open_source_server():
    helpers.docker_compose_start(
        project_name=project_name,
        files=["../mender_server/docker-compose.yml"],
    )


# The following functions make sure the server
# is started _once_ per test run, regardless of
# multiple workers with xdist
def pytest_sessionstart(session):
    # If the worker id is master, then we're not running xdist
    master = "master" == xdist.get_xdist_worker_id(session)
    if master or xdist.is_xdist_controller(session):
        start_open_source_server()


def pytest_sessionfinish(session):
    # If the worker id is master, then we're not running xdist
    master = "master" == xdist.get_xdist_worker_id(session)
    if master or xdist.is_xdist_controller(session):
        stop_open_source_server()


@pytest.fixture(scope="function", autouse=True)
def devices_down():
    yield
    helpers.docker_compose_stop(
        project_name=project_name_client,
        files=["docker-compose.client.yml"],
    )


@pytest.fixture(scope="function")
def standard_setup_one_client_bootstrapped():
    env = Env()
    env.server = Server()

    env.devices = clients_up(1)
    env.device = env.devices[0]

    wait_for_devices(env)

    env.server.accept_devices(env.devices)

    return env
