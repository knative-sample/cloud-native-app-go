package app

import (
	"strings"

	"fmt"
	"os"

	"net"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-app-go/weather/cmd/weather/app/options"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/version"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/weather"
	"github.com/spf13/cobra"
)

// start edas api
func NewCommandStartServer(stopCh <-chan struct{}) *cobra.Command {
	ops := &options.Options{}
	mainCmd := &cobra.Command{
		Short: "AppOS",
		Long:  "Application Operating System",
		RunE: func(c *cobra.Command, args []string) error {
			glog.V(2).Infof("NewCommandStartServer main:%s", strings.Join(args, " "))
			run(stopCh, ops)
			return nil
		},
	}

	ops.SetOps(mainCmd)
	return mainCmd
}

func run(stopCh <-chan struct{}, ops *options.Options) {
	vs := version.Version().Info("Application Operating System")
	if ops.Version {
		fmt.Println(vs)
		os.Exit(0)
	}

	cityHost, cityPort, err := net.SplitHostPort(ops.CityService)
	if err != nil {
		glog.Fatalf("parse CityService:%s error:%s", ops.CityService)
	}
	detailHost, detailPort, err := net.SplitHostPort(ops.DetailService)
	if err != nil {
		glog.Fatalf("parse DetailService:%s error:%s", ops.DetailService)
	}

	wa := weather.WebApi{
		Port:         ops.Port,
		ResourceRoot: ops.ResourceRoot,
		CityService: &weather.Service{
			Host: cityHost,
			Port: cityPort,
		},
		DetailService: &weather.Service{
			Host: detailHost,
			Port: detailPort,
		},
	}
	go func() {
		if err := wa.Start(); err != nil {
			glog.Fatalf("start Webserver error:%s", err.Error())
		}
	}()
	<-stopCh
}
