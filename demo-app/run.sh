set -ex
export ETCD_ENDPOINTS=http://localhost:2001,http://localhost:2002,http://localhost:2003
go run app.go
