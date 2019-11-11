package detail

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/db"
)

func NewCity(tsc *db.TableStoreConfig) *DetailManager {
	return &DetailManager{TableStoreConfig: tsc}
}

func (cm *DetailManager) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/area/{adcode}/{date}", cm.getAreaDetail).Methods("GET")
	http.Handle("/", router)

	router.Use(cm.AccessLog)
	http.ListenAndServe(fmt.Sprintf(":%s", cm.Port), nil)

	return nil
}
