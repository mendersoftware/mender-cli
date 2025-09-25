# Copyright 2025 Northern.tech AS
#
#    All Rights Reserved


from os import path
import sys

import os
import xdist
import logging
import subprocess
import shutil
import pytest

from filelock import FileLock

sys.path += [
    path.join(path.dirname(__file__), "mender_server/backend/tests/integration/")
]

from server import Server
from setup import (
    Env,
    gateway_up,
    clients_up,
    wait_for_devices,
    project_name,
    project_name_client,
)

import helpers

THIS_DIR = os.path.dirname(os.path.abspath(__file__))

logging.getLogger("requests").setLevel(logging.CRITICAL)
logging.getLogger("urllib3").setLevel(logging.CRITICAL)
logging.getLogger("redo").setLevel(logging.INFO)
logger = logging.getLogger()

collect_ignore = ["mender_server"]


def pytest_addoption(parser):
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )


@pytest.fixture(scope="session", autouse=True)
def check_variables():
    if "MENDER_GATEWAY_IMAGE" not in os.environ:
        pytest.fail(f"MENDER_GATEWAY_IMAGE is not set")


def _extract_fs_from_image(request, client_compose_file, filename):
    if os.path.exists(os.path.join(THIS_DIR, filename)):
        return filename

    d = os.path.join(THIS_DIR, "output")

    def cleanup():
        shutil.rmtree(d, ignore_errors=True)

    request.addfinalizer(cleanup)

    image_type = "client"
    if "gateway" in os.path.basename(client_compose_file):
        image_type = "gateway"

    with helpers.docker_lock:
        image = (
            subprocess.check_output(
                [
                    "docker",
                    "compose",
                    "-f",
                    client_compose_file,
                    "config",
                    "--images",
                    image_type,
                ],
                env=os.environ,
            )
            .decode("UTF-8")
            .strip()
            .split("\n")[0]
        )

        subprocess.check_call(["mkdir", "-p", "output/"])
        subprocess.check_call(
            [
                "docker",
                "run",
                "--rm",
                "--entrypoint",
                "/extract_fs",
                "--volume",
                d + ":/output",
                image,
            ]
        )

    shutil.move(
        os.path.join(d, "core-image-full-cmdline-qemux86-64.ext4"),
        os.path.join(THIS_DIR, filename),
    )

    return filename


def _image(request, compose_file, filename):
    return _extract_fs_from_image(
        request, os.path.join(THIS_DIR, compose_file), filename
    )


# https://pytest-xdist.readthedocs.io/en/latest/how-to.html#making-session-scoped-fixtures-execute-only-once
def run_if_non_existent(request, tmp_path_factory, compose_file, filename):
    """
    Makes session scoped fixtures only run once
    to avoid conflict while running multiple workers with xdist
    """
    root_tmp_dir = tmp_path_factory.getbasetemp().parent
    file_lock = root_tmp_dir / filename

    with FileLock(f"{file_lock}.lock"):
        if file_lock.is_file() and os.path.exists(os.path.join(THIS_DIR, filename)):
            return filename
    return _image(request, compose_file, filename)


@pytest.fixture(scope="session")
def client_device_image(request, tmp_path_factory, worker_id):
    compose_file = "docker-compose.client.yml"
    filename = f"client-device-image-full-cmdline-qemux86-64.ext4"

    if worker_id == "master":
        return _image(request, compose_file, filename)

    return run_if_non_existent(request, tmp_path_factory, compose_file, filename)


@pytest.fixture(scope="session")
def gateway_device_image(request, tmp_path_factory, worker_id):
    compose_file = "docker-compose.gateway.yml"
    filename = f"gateway-device-image-full-cmdline-qemux86-64.ext4"

    if worker_id == "master":
        return _image(request, compose_file, filename)

    return run_if_non_existent(request, tmp_path_factory, compose_file, filename)


@pytest.fixture(scope="session")
def broken_update_image():
    image = "broken_update.ext4"
    subprocess.check_call(["dd", "if=/dev/urandom", f"of={image}", "bs=10M", "count=5"])
    return image


def stop_open_source_server():
    helpers.docker_compose_stop(
        project_name=project_name, files=["mender_server/docker-compose.yml"],
    )


def start_open_source_server():
    helpers.docker_compose_start(
        project_name=project_name, files=["mender_server/docker-compose.yml"],
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
        files=["docker-compose.client.yml", "docker-compose.gateway.yml",],
    )


@pytest.fixture(scope="function")
def standard_setup_one_client_bootstrapped_with_gateway():
    env = Env()
    env.server = Server()

    env.gateway = gateway_up()
    env.devices = clients_up(1)
    env.device = env.devices[0]

    wait_for_devices(env)

    env.server.accept_devices(env.devices + [env.gateway])

    return env


@pytest.fixture(scope="function")
def standard_setup_two_clients_bootstrapped_with_gateway():
    env = Env()
    env.server = Server()

    env.gateway = gateway_up()
    env.devices = clients_up(2)
    env.device = env.devices[0]

    wait_for_devices(env)

    env.server.accept_devices(env.devices + [env.gateway])

    return env
