package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/alexellis/faas/gateway/metrics"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

type Camera struct {
	Name              string `json:"Name"`
	Status            bool   `json:"status"`
	Screenshot        string `json:"screenshot"`
	UsesMobileNetwork bool   `json:"usesMobileNetwork"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	cameraRouter := mux.NewRouter().StrictSlash(true)
	//cameraRouter.HandleFunc("/", indexPage)
	cameraRouter.HandleFunc("/", returnAllCameras)
	cameraRouter.HandleFunc("/camera/{name}", returnCameraStatus)
	cameraRouter.HandleFunc("/screenshot/{name}", returnCameraScreenshot)
	metricsHandler := metrics.PrometheusHandler()
	cameraRouter.Handle("/metrics", metricsHandler)

	host, exists := os.LookupEnv("HOST")
	if !exists {
		host = "0.0.0.0"
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	log.WithFields(logrus.Fields{
		"host": host,
		"port": port,
	}).Info("Potty server is listening...")
	log.Fatal(http.ListenAndServe(host+":"+port, cameraRouter))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func returnAllCameras(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"cameraCount": len(Cameras),
	}).Info("Returning camera list.")
	json.NewEncoder(w).Encode(Cameras)
}

func returnCameraStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnCameraStatus")
	vars := mux.Vars(r)
	key := vars["name"]

	for _, camera := range Cameras {
		if camera.Name == key {
			if camera.Status {
				if camera.UsesMobileNetwork {
					time.Sleep(10 * time.Second)
				}
				json.NewEncoder(w).Encode("Running...")
			} else {
				json.NewEncoder(w).Encode("Failed to reach camera...")
			}
		}
	}
}

func returnCameraScreenshot(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["name"]

	log.WithFields(logrus.Fields{
		"name": key,
	}).Info("Getting screenshot...")

	for _, camera := range Cameras {
		if camera.Name == key {
			if camera.Status {
				serviceName := "returnCameraScreenshot"
				labels := prometheus.Labels{"http_code": "200", "function_name": serviceName}
				StreamingCount.With(labels).Add(1)
				if camera.UsesMobileNetwork {
					sleepTime := rand.Intn(10)
					log.WithFields(logrus.Fields{
						"name":           key,
						"connectionTime": sleepTime,
					}).Warn("Camera uses mobile network.")
					time.Sleep(time.Duration(sleepTime) * time.Second)
				}
				log.WithFields(logrus.Fields{
					"name": key,
				}).Info("Snapshot taken.")
				json.NewEncoder(w).Encode(camera.Screenshot)
			} else {
				log.WithFields(logrus.Fields{
					"name": key,
				}).Error("Camera status is OFFLINE")
				json.NewEncoder(w).Encode("Failed to reach camera...")
			}
		}
	}
}

var Cameras []Camera
var StreamingCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "number_of_streamings",
		Help: "Number of stramings",
	},
	[]string{"function_name", "http_code"},
)
var log = logrus.New()

func main() {

	Cameras = []Camera{
		Camera{Name: "Bagoly", Status: true, Screenshot: "bagoly.txt", UsesMobileNetwork: true},
		Camera{Name: "Cinke", Status: false, Screenshot: "cinke.txt", UsesMobileNetwork: false},
		Camera{Name: "Rozsdafarku", Status: true, Screenshot: "rozsdafarku", UsesMobileNetwork: false},
	}

	prometheus.Register(StreamingCount)

	client, err := elastic.NewClient(elastic.SetURL("http://okd-5mthh-worker-tb667.apps.okd.codespring.ro:30029"), elastic.SetSniff(false))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(client, "http://okd-5mthh-worker-tb667.apps.okd.codespring.ro:30029", logrus.DebugLevel, "gyeka-000001")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	handleRequests()
}
