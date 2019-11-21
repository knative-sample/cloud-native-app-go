package weather

import (
	"net/http"

	"fmt"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/tracing"
)

func (wa *WebApi) Start() error {
	glog.Infof("Starting webapi, ResourceRoot:%s", wa.ResourceRoot)
	wa.tracer = tracing.GetTracer(wa.ServiceName, wa.InstanceIp, wa.ZipKinEndpoint)

	router := mux.NewRouter()
	router.PathPrefix("/index").Handler(http.StripPrefix("/index", http.FileServer(http.Dir(wa.ResourceRoot))))
	router.HandleFunc("/api/city/list", wa.CityList).Methods("GET")
	router.HandleFunc("/api/city/detail/{name}/{date}", wa.Detail)
	http.Handle("/", router)
	router.Use(wa.AccessLog)
	router.Use(wa.TraceLog)
	http.ListenAndServe(fmt.Sprintf(":%s", wa.Port), nil)

	return nil
}
