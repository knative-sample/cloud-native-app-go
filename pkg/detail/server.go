package detail

import (
	"context"
	"log"
	"net"

	"time"

	"fmt"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/db"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/tracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedDetailServer
	Port             string
	ZipKinEndpoint   string
	ServiceName      string
	InstanceIp       string
	TableStoreConfig *db.TableStoreConfig
	tracer           *zipkin.Tracer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) GetDetail(ctx context.Context, in *DetailQuery) (*DetailInfo, error) {
	glog.Infof("Received query: %s", in.GetCitycode())

	if parent := zipkin.SpanFromContext(ctx); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("detail_sub_span", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}
	w, err := s.TableStoreConfig.QueryWeather(in.GetCitycode(), in.Date)
	if err != nil {
		log.Printf("QueryWeather error %s", err.Error())
		return nil, err
	}
	return &DetailInfo{
		Adcode:       in.GetCitycode(),
		Name:         w.City,
		Date:         w.Date,
		Daypower:     w.Daypower,
		Daytemp:      w.Daytemp,
		Dayweather:   w.Dayweather,
		Daywind:      w.Daywind,
		Nightpower:   w.Nightpower,
		Nighttemp:    w.Nighttemp,
		Nightweather: w.Nightweather,
		Nightwind:    w.Nightwind,
		Province:     w.Province,
		Reporttime:   w.Reporttime,
		Week:         w.Week,
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.tracer = tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
	gs := grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(s.tracer)))
	RegisterDetailServer(gs, s)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
