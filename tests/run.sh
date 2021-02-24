#!/bin/bash

DIR=$(readlink -f $(dirname $0))

[ $$ -eq 1 ] && sleep 10

mkdir -p /tests/cover

py.test -s --tb=short \
          --verbose --junitxml=$DIR/results.xml \
          $DIR/tests/test_*.py "$@"
