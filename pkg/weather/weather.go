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

	"github.com/knative-sample/cloud-native-app-go/weather/pkg/detail"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
)

func (wa *WebApi) CityList(w http.ResponseWriter, r *http.Request) {
	currentSpan := wa.NewSpan("GetCityList", r.Context())
	defer currentSpan.Finish()

	// TODO 换成具体的服务地址
	addr := "127.0.0.1:9090"
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
	glog.Infof("Greeting cityes: %s", cts)

	fmt.Fprintf(w, "City Hello, %s!", r.URL.Path[1:])
}

func (wa *WebApi) Detail(w http.ResponseWriter, r *http.Request) {
	currentSpan := wa.NewSpan("GetAreaList", r.Context())
	defer currentSpan.Finish()

	childSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	defer childSpan.Finish()

	// TODO 1. get city areas 2. foreach area get weather info
	areaChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
	areas, err := wa.getAreas("0571", areaChildSpan)
	if err != nil {
		glog.Errorf("getAreas error:%s", err.Error())
	}
	defer areaChildSpan.Finish()

	detailResult := []*detail.DetailInfo{}
	for _, a := range areas {
		detailChildSpan := wa.tracer.StartSpan("GetDetail", zipkin.Parent(currentSpan.Context()))
		d, err := wa.getDetail(a.Citycode, "2019-11-12", detailChildSpan)
		if err != nil {
			glog.Errorf("getAreas error:%s", err.Error())
		}
		detailChildSpan.Finish()

		detailResult = append(detailResult, d)
	}

	dbts, _ := json.Marshal(detailResult)
	fmt.Fprintf(w, "Detail Info, %s!", dbts)
}

func (wa *WebApi) getAreas(cityCode string, currentSpan zipkin.Span) ([]*city.Area, error) {
	// TODO 换成具体的服务地址
	addr := "127.0.0.1:9090"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(wa.tracer)))
	if err != nil {
		glog.Fatalf("grpc.Dial error:%s", err.Error())
	}
	defer conn.Close()

	client := city.NewCityManagerClient(conn)

	cts, err := client.AreaList(zipkin.NewContext(context.Background(), currentSpan), &city.AreaQuery{Citycode: "0571"})
	if err != nil {
		glog.Fatalf("AreaList error: %v", err)
	}
	glog.Infof("getAreas areas: %s", cts)
	return cts.Areas, nil
}

func (wa *WebApi) getDetail(cityCode, date string, currentSpan zipkin.Span) (*detail.DetailInfo, error) {
	// TODO 换成具体的服务地址
	addr := "127.0.0.1:9091"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(wa.tracer)))
	if err != nil {
		glog.Fatalf("grpc.Dial error:%s", err.Error())
	}
	defer conn.Close()

	client := detail.NewDetailClient(conn)

	cts, err := client.GetDetail(zipkin.NewContext(context.Background(), currentSpan), &detail.DetailQuery{Citycode: "0571", Date: date})
	if err != nil {
		glog.Fatalf("GetDetail error: %v", err)
	}
	glog.Infof("getDetail areas: %s", cts)
	return cts, nil
}
