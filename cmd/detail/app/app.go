package app

import (
	"strings"

	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-app-go/weather/cmd/detail/app/options"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/db"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/detail"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/version"
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

	dm := &detail.DetailManager{
		Port:             ops.Port,
		TableStoreConfig: &db.TableStoreConfig{},
	}
	go func() {
		if err := dm.Start(); err != nil {
			glog.Fatalf("start Webserver error:%s", err.Error())
		}
	}()

	<-stopCh
}
