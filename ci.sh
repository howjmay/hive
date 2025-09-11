#!/bin/bash

git pull --rebase origin main
mkdir -p workspace
./hive --sim ethereum/eest --client wasp-client -docker.pull -sim.parallelism 10  > workspace/test.log 2>&1