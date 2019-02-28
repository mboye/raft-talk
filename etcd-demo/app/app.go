package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/mboye/raft-talk/etcd-demo/app/indexhandler"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

func main() {

	endpointsArg, found := os.LookupEnv("ETCD_ENDPOINTS")
	if !found {
		log.Fatal("environment variable missing: ETCD_ENDPOINTS")
	}
	endpoints := strings.Split(endpointsArg, ",")

	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	clusterClient, err := clientv3.New(cfg)
	if err != nil {
		log.WithField("err", err).Fatal("failed to create etcd client")
	}

	http.Handle("/", indexhandler.New(clusterClient))
	log.Info("starting server on port 8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	log.WithError(err).Error("server stopped")
}
