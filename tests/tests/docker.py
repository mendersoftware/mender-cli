#!/usr/bin/python
# Copyright 2018 Mender Software AS
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        https://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
import cli

BASE_COMPOSE_FILES = [
    "/tests/integration/docker-compose.yml",
    "/tests/integration/docker-compose.storage.minio.yml",
    "/tests/integration/docker-compose.testing.yml",
]

def exec(service, files, *argv):
    c = cli.Cli('/usr/bin/docker-compose')

    # set by run-test-environment
    args = ['-p', 'acceptance-tests']

    for f in files:
        args += ['-f', f]

    args += ['exec', '-T', service] + list(argv)

    return c.run(*args)
