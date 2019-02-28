package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type leaderElectorResponse struct {
	Name string `json:"name"`
}

func main() {
	log.Info("leader election app starting")
	indexTemplate, err := template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		log.WithError(err).Fatal("failed to parse index template")
	}

	podName, found := os.LookupEnv("HOSTNAME")
	if !found {
		log.Fatal("environment variable missing: HOSTNAME")
	}

	handler := func(response http.ResponseWriter, request *http.Request) {
		indexTemplate.ExecuteTemplate(response, "index.html", isLeader(podName))
	}

	http.HandleFunc("/", handler)

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	log.WithError(err).Fatal("server stopped")
}

func isLeader(podName string) bool {
	response, err := http.Get("http://localhost:4040/")
	if err != nil {
		log.WithError(err).Error("failed to check leadership status")
		return false
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Error("failed to parse leader elector response")
		return false
	}

	electorResponse := leaderElectorResponse{}
	if err := json.Unmarshal(body, &electorResponse); err != nil {
		log.WithError(err).Error("failed to unmarshal leader elector response")
		return false
	}

	result := electorResponse.Name == podName
	log.WithField("is_leader", result).Info("pod leadership status")
	return result
}
