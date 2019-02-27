package indexhandler

import (
	"context"
	"html/template"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

const backgroundColorConfigKey = "app-config/background-color"
const defaultBackgroundColor = "#ffffff"

type indexHandler struct {
	clusterClient   *clientv3.Client
	indexTemplate   *template.Template
	backgroundColor string
	mutex           sync.Mutex
}

// New returns a new request handler for the index page
func New(clusterClient *clientv3.Client) http.Handler {
	handler := &indexHandler{
		clusterClient:   clusterClient,
		backgroundColor: defaultBackgroundColor}

	var err error
	handler.indexTemplate, err = template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		log.WithError(err).Fatal("failed to parse index template")
	}

	handler.loadBackgroundColor()
	go handler.watchBackgroundColor()
	return handler
}

func (h *indexHandler) loadBackgroundColor() {
	backgroundColor := defaultBackgroundColor
	response, err := h.clusterClient.Get(context.Background(), backgroundColorConfigKey)
	if err != nil {
		log.WithError(err).Error("failed to get background color from etcd")
	} else {
		found := false
		for _, kv := range response.Kvs {
			if string(kv.Key) == backgroundColorConfigKey {
				found = true
				backgroundColor = string(kv.Value)
				log.WithField("background_color", backgroundColor).Info("loaded background color from etcd")
				break
			}
		}

		if !found {
			log.Error("background color configuration not found in etcd")
		}
	}

	h.mutex.Lock()
	h.backgroundColor = backgroundColor
	h.mutex.Unlock()
	log.WithField("background_color", backgroundColor).Info("current background color")
}

func (h *indexHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.mutex.Lock()
	h.indexTemplate.ExecuteTemplate(response, "index.html", h.backgroundColor)
	log.WithFields(log.Fields{
		"path":       request.URL.Path,
		"method":     request.Method,
		"user_agent": request.Header.Get("user-agent")}).Info("handled request")
	h.mutex.Unlock()
}

func (h *indexHandler) watchBackgroundColor() {
	log.Info("watching for background color changes in etcd")
	changes := h.clusterClient.Watch(context.Background(), backgroundColorConfigKey, clientv3.WithPrevKV())
	for change := range changes {
		for _, event := range change.Events {
			if string(event.Kv.Key) == backgroundColorConfigKey {
				oldValue := ""
				if event.PrevKv != nil {
					oldValue = string(event.PrevKv.Value)
				}

				newValue := string(event.Kv.Value)
				log.WithFields(log.Fields{
					"old_value": oldValue,
					"new_value": newValue}).Info("background color changed")

				h.mutex.Lock()
				h.backgroundColor = newValue
				h.mutex.Unlock()
				break
			}
		}
	}
}
