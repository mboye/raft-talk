#!/bin/sh
docker run --rm -it --network raft-talk_default --volume "$(PWD)/scripts:/scripts" quay.io/coreos/etcd:v3.3.12 sh
