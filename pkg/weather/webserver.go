package weather

import (
	"net/http"

	"fmt"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func (wa *WebApi) Start() error {
	glog.Infof("Starting webapi, ResourceRoot:%s", wa.ResourceRoot)

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(wa.ResourceRoot))))
	router.HandleFunc("/city/{name}", wa.City).Methods("GET")
	router.HandleFunc("/detail/{.*}", wa.Detail)
	http.Handle("/", router)

	router.Use(wa.AccessLog)
	http.ListenAndServe(fmt.Sprintf("%s", wa.Port), nil)

	return nil
}
