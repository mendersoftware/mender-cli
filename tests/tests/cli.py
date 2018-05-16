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
import subprocess

MENDER_CLI='tests/mender-cli'

class Cli:
    """Simple wrapper for subprocess"""
    def __init__(self, path=MENDER_CLI):
        self.path = path

    def run(self, *argv):
        """Returns a CompletedProcess wrapped in CliResult"""
        args = [self.path] + list(argv)
        completed = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return CliResult(completed)

class CliResult:
    """Wraps CompletedProcess, decodes output to strings"""
    def __init__(self, completed_process):
        self.completed_process = completed_process

        self.returncode = completed_process.returncode
        self.stdout = self.completed_process.stdout.decode('utf-8')
        self.stderr = self.completed_process.stderr.decode('utf-8')
