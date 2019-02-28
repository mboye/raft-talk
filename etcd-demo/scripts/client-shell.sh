#!/bin/sh
docker run --rm -it --network etcd-demo_default --volume "$(PWD)/scripts:/scripts" quay.io/coreos/etcd:v3.3.12 sh
