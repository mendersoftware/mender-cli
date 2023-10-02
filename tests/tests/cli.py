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
import subprocess
import tempfile

MENDER_CLI = "tests/mender-cli"
MENDER_CLI_COVERAGE = "tests/mender-cli-test"
COVERAGE_DIR = "/tests/cover"


class Cli:
    """Simple wrapper for subprocess"""

    def __init__(self, path=MENDER_CLI):
        self.path = path

    def run(self, *argv):
        """Returns a CompletedProcess wrapped in CliResult"""
        args = [self.path] + list(argv)
        completed = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return CliResult(completed)

    def run_and_enter_password(self, *argv, password=""):
        """Returns a CompletedProcess wrapped in CliResult"""
        args = [self.path] + list(argv)
        p = subprocess.Popen(
            args, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE
        )
        stdout, stderr = p.communicate(input=password.encode() + b"\n")
        return CliResult(p, stdout=stdout, stderr=stderr)


class CliResult:
    """Wraps CompletedProcess, decodes output to strings"""

    def __init__(self, completed_process, stdout=None, stderr=None):
        self.completed_process = completed_process

        self.returncode = completed_process.returncode
        self.stdout = (
            self.completed_process.stdout.decode("utf-8") if stdout is None else stdout
        )
        self.stderr = (
            self.completed_process.stderr.decode("utf-8") if stderr is None else stderr
        )


class MenderCliCoverage(Cli):
    def __init__(self, path=MENDER_CLI_COVERAGE):
        self.path = path

    def _next_coverage_file(self):
        coverfile = tempfile.NamedTemporaryFile(
            delete=False, dir=COVERAGE_DIR, prefix="coverage-acceptance-", suffix=".txt"
        )
        coverfile.close()
        return coverfile.name

    def run(self, *argv):
        args = [
            self.path,
            "-test.coverprofile=" + self._next_coverage_file(),
            "-acceptance-tests",
            "-test.run=TestMain",
            "-cli-args=" + " ".join(list(argv)),
        ]
        completed = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return CliResult(completed)
