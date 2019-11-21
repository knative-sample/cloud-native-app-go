package city

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
	UnimplementedCityManagerServer
	Port             string
	ZipKinEndpoint   string
	ServiceName      string
	InstanceIp       string
	TableStoreConfig *db.TableStoreConfig
	tracer           *zipkin.Tracer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) CityList(ctx context.Context, in *CityQuery) (*Citys, error) {
	glog.Infof("CityList Received query: %s", in.GetCitycode())

	if parent := zipkin.SpanFromContext(ctx); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("city_list", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}
	cities, err := s.TableStoreConfig.QueryCities()
	if err != nil {
		log.Printf("QueryCities error %s", err.Error())
		return nil, err
	}
	cs := make([]*City, 0)
	for _, city := range cities {
		c := &City{
			Name:     city.Name,
			Citycode: city.Citycode,
		}
		cs = append(cs, c)
	}
	return &Citys{Citys: cs}, nil
}

func (s *Server) AreaList(ctx context.Context, in *AreaQuery) (*Areas, error) {
	glog.Infof("AreaList Received query: %s", in.GetCitycode())

	if parent := zipkin.SpanFromContext(ctx); parent != nil {
		//tracer := tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
		subSpan := s.tracer.StartSpan("city_area_list", zipkin.Parent(parent.Context()))
		defer subSpan.Finish()
		//do some operations
		time.Sleep(time.Millisecond * 10)
	}

	areas, err := s.TableStoreConfig.QueryAreaByCitycode(in.Citycode)
	if err != nil {
		log.Printf("QueryAreaByCitycode error %s", err.Error())
		return nil, err
	}
	as := make([]*Area, 0)
	for _, city := range areas {
		area := &Area{
			Name:     city.Name,
			Citycode: city.Adcode,
		}
		as = append(as, area)
	}
	return &Areas{Areas: as}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.tracer = tracing.GetTracer(s.ServiceName, s.InstanceIp, s.ZipKinEndpoint)
	gs := grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(s.tracer)))
	RegisterCityManagerServer(gs, s)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
