# Copyright 2025 Northern.tech AS
#
#    All Rights Reserved
#

import pytest
import filelock
import tempfile
import random
import logging

import logging
import tempfile
import subprocess


artifact_lock = filelock.FileLock(".artifact_modification_lock")
docker_lock = filelock.FileLock(".docker_lock")
logger = logging.getLogger(__name__)


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


def install_community_update_module(device, module):
    url = f"https://raw.githubusercontent.com/mendersoftware/mender-update-modules/master/{module}/module/{module}"
    device.run("mkdir -p /usr/share/mender/modules/v3")
    device.run(f"wget -P /usr/share/mender/modules/v3 {url}")
    device.run(f"chmod +x /usr/share/mender/modules/v3/{module}")


def get_script_artifact(script, artifact_name, output_path, device_type="qemux86-64"):
    with tempfile.NamedTemporaryFile(suffix="testdeployment") as tf:
        tf.write(script)
        tf.seek(0)
        out = tf.read()
        logger.info(f"Script: {out}")
        script_path = tf.name

        cmd = [
            "mender-artifact",
            "write",
            "module-image",
            "-T",
            "script",
            "-n",
            artifact_name,
            "-t",
            device_type,
            "-o",
            output_path,
            "-f",
            script_path,
        ]

        logger.info(f"Executing command: {cmd}")
        subprocess.check_call(cmd)

        return output_path


def make_rootfs_artifact(
    image, device_type, artifact_name, artifact_filename,
):

    cmd = [
        "mender-artifact",
        "write",
        "rootfs-image",
        "-f",
        image,
        "-t",
        device_type,
        "-n",
        artifact_name,
        "-o",
        artifact_filename,
    ]

    logger.info(f"Running: {cmd}")
    subprocess.check_call(cmd)

    return artifact_filename


def update_procedure(
    install_image=None,
    devices=None,
    device_type="qemux86-64",
    deployment_triggered_callback=lambda: None,
    env=None,
):

    with artifact_lock:
        artifact_name = f"mender-{random.randint(0, 99999999)}"

        logger.debug("randomized image id: " + artifact_name)

        with tempfile.NamedTemporaryFile() as artifact_file:
            created_artifact = make_rootfs_artifact(
                install_image, device_type, artifact_name, artifact_file.name,
            )

            devices_ids = [get_device_id(device, env.server) for device in devices]
            if created_artifact:
                env.server.upload_image(created_artifact)
                if devices is None:
                    devices = list(set([device_id for device_id in devices_ids]))
                deployment_id = env.server.create_deployment(
                    artifact_name=artifact_name, device_ids=devices_ids
                )
            else:
                logger.warn("failed to create artifact")
                pytest.fail("error creating artifact")

    deployment_triggered_callback()
    return deployment_id, artifact_name


def update_image(
    device,
    host_ip,
    expected_mender_clients=1,
    install_image=None,
    devices=None,
    deployment_triggered_callback=lambda: None,
    env=None,
):
    """
        Perform a successful upgrade, and assert that deployment status/logs are correct.

        A reboot is performed, and running partitions have been swapped.
        Deployment status will be set as successful for device.
        Logs will not be retrieved, and result in 404.
    """

    previous_inactive_part = device.get_passive_partition()
    with device.get_reboot_detector(host_ip) as reboot:
        deployment_id, _ = update_procedure(
            install_image,
            devices=devices,
            deployment_triggered_callback=deployment_triggered_callback,
            env=env,
        )
        reboot.verify_reboot_performed()

        assert (
            device.get_active_partition()[-1] == previous_inactive_part[-1]
        ), "device did not flip partitions during update"

    env.server.check_expected_statistics(
        deployment_id, "success", expected_mender_clients
    )
    return deployment_id


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
