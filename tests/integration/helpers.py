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

import filelock
import logging

import logging
import subprocess

import hashlib


docker_lock = filelock.FileLock(".docker_lock")


def md5sum(fname):
    hash_md5 = hashlib.md5()
    with open(fname, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_md5.update(chunk)
    return hash_md5.hexdigest()


def get_mac_address(device):
    result = device.run(
        "/usr/share/mender/identity/mender-device-identity", hide=True, warn_only=True
    ).strip()
    return result.split("=")[-1]


def get_device_id(device, server):
    mac_address = get_mac_address(device)
    device = next(
        device
        for device in server.get_accepted_devices()
        if device["identity_data"]["mac"] == mac_address
    )
    return device["id"]


def docker_compose_start(project_name, files):
    with docker_lock:
        cmd = ["docker", "compose", "--project-name", project_name]
        for file in files:
            cmd.extend(["--file", file])
        cmd += ["up", "--detach"]
        subprocess.check_call(cmd)


def docker_compose_stop(project_name, files):
    with docker_lock:
        cmd = ["docker", "compose", "--project-name", project_name]
        for file in files:
            cmd.extend(["--file", file])
        cmd += ["down", "--volumes", "--remove-orphans"]
        subprocess.check_call(cmd)
