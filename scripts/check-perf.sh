#!/bin/sh
/scripts/etcdctl.sh del --prefix /etcdctl-check-perf/
/scripts/etcdctl.sh check perf
