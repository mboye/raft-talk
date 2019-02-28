#!/bin/sh
ETCDCTL_API=3 etcdctl --endpoints "http://node-1:2379,http://node-2:2379,http://node-3:2379" $@
