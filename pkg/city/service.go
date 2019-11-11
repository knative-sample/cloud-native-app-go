package city

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/db"
)

func NewCity(tsc *db.TableStoreConfig) *CityManager {
	return &CityManager{TableStoreConfig: tsc}
}

func (cm *CityManager) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/list", cm.cityList).Methods("GET")
	http.Handle("/", router)

	router.Use(cm.AccessLog)
	http.ListenAndServe(fmt.Sprintf(":%s", cm.Port), nil)

	return nil
}
