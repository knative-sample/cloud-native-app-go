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

	return &DetailInfo{
		Adcode:       in.GetCitycode(),
		Name:         "西湖区",
		Date:         "2019-09-26",
		Daypower:     "5",
		Daytemp:      "30",
		Dayweather:   "晴",
		Daywind:      "东",
		Nightpower:   "5",
		Nighttemp:    "19",
		Nightweather: "晴",
		Nightwind:    "东",
		Province:     "浙江",
		Reporttime:   "2019-09-26 22:49:20",
		Week:         "4",
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
