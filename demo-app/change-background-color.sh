#!/bin/bash
set -ex
etcdctl --endpoints="http://localhost:2001,http://localhost:2002,http://localhost:2003" put  "app-config/background-color" "$1"
