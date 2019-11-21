package weather

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/city"
	"github.com/openzipkin/zipkin-go"

	//zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"encoding/json"

	"strings"

	"github.com/knative-sample/cloud-native-app-go/weather/pkg/detail"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
)

func (wa *WebApi) CityList(w http.ResponseWriter, r *http.Request) {
	currentSpan := wa.NewSpan("GetCityList", r.Context())
	defer currentSpan.Finish()

	addr := fmt.Sprintf("%s:%s", wa.CityService.Host, wa.CityService.Port)
	if wa.tracer == nil {
		glog.Fatalf("wa tracer is nill")
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(wa.tracer)))
	if err != nil {
		glog.Fatalf("grpc.Dial error:%s", err.Error())
	}
	defer conn.Close()

	client := city.NewCityManagerClient(conn)

	cts, err := client.CityList(zipkin.NewContext(context.Background(), currentSpan), &city.CityQuery{Citycode: "010"})
	if err != nil {
		glog.Fatalf("could not greet: %v", err)
	}
	cityInfos, _ := json.Marshal(cts)
	fmt.Fprintf(w, string(cityInfos))
}

func (wa *WebApi) Detail(w http.ResponseWriter, r *http.Request) {
	currentSpan := wa.NewSpan("GetAreaList", r.Context())
	defer currentSpan.Finish()

	childSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	defer childSpan.Finish()

	// 1. get city areas 2. foreach area get weather info
	params := strings.TrimPrefix(r.URL.Path[1:], "api/city/detail/")
	vars := strings.Split(params, "/")
	citycode := vars[0]
	date := vars[1]
	glog.Infof("citycode: %s", citycode)
	glog.Infof("date: %s", date)
	areaChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	areas, err := wa.getAreas(citycode, areaChildSpan)
	if err != nil {
		glog.Errorf("getAreas error:%s", err.Error())
	}
	defer areaChildSpan.Finish()

	detailResult := []*detail.DetailInfo{}
	for _, a := range areas {
		detailChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
		d, err := wa.getDetail(a.Citycode, date, detailChildSpan)
		if err != nil {
			glog.Errorf("getAreas error:%s", err.Error())
		}
		if d.Name == "" {
			continue
		}
		detailChildSpan.Finish()

		detailResult = append(detailResult, d)
	}

	dbts, _ := json.Marshal(detailResult)
	fmt.Fprintf(w, string(dbts))
}

func (wa *WebApi) getAreas(cityCode string, currentSpan zipkin.Span) ([]*city.Area, error) {
	//addr := "127.0.0.1:9090"
	addr := fmt.Sprintf("%s:%s", wa.CityService.Host, wa.CityService.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(wa.tracer)))
	if err != nil {
		glog.Fatalf("grpc.Dial error:%s", err.Error())
	}
	defer conn.Close()

	client := city.NewCityManagerClient(conn)

	cts, err := client.AreaList(zipkin.NewContext(context.Background(), currentSpan), &city.AreaQuery{Citycode: cityCode})
	if err != nil {
		glog.Fatalf("AreaList error: %v", err)
	}
	glog.Infof("getAreas areas: %s", cts)
	return cts.Areas, nil
}

func (wa *WebApi) getDetail(cityCode, date string, currentSpan zipkin.Span) (*detail.DetailInfo, error) {
	//addr := "127.0.0.1:9091"
	addr := fmt.Sprintf("%s:%s", wa.DetailService.Host, wa.DetailService.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(wa.tracer)))
	if err != nil {
		glog.Fatalf("grpc.Dial error:%s", err.Error())
	}
	defer conn.Close()

	client := detail.NewDetailClient(conn)

	cts, err := client.GetDetail(zipkin.NewContext(context.Background(), currentSpan), &detail.DetailQuery{Citycode: cityCode, Date: date})
	if err != nil {
		glog.Fatalf("GetDetail error: %v", err)
	}
	glog.Infof("getDetail areas: %s", cts)
	return cts, nil
}
