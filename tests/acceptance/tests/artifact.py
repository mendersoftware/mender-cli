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
import tempfile
import logging

import cli

MENDER_ARTIFACT_TOOL='tests/mender-artifact'

def create_artifact_file(outpath):
    with tempfile.NamedTemporaryFile(prefix='menderin') as infile:
        logging.info('writing mender artifact to %s', outpath)

        infile.write(b'bogus test data')
        infile.flush()

        c = cli.Cli(MENDER_ARTIFACT_TOOL)
        r = c.run('write', 'rootfs-image',
                  '--device-type', 'device-foo',\
                  '--file', infile.name, \
                  '--artifact-name', 'artifact-foo',\
                  '--output-path', outpath)

        if r.returncode != 0:
            msg = 'mender-artifact failed with code {}, \nstdout: \n{}stderr: {}\n'.format(r.returncode, r.stdout, r.stderr)
            logging.error(msg)
            raise RuntimeError(msg)

        logging.info('mender artifact written to %s', outpath)
