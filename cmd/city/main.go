package main

import (
	"flag"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-app-go/weather/cmd/city/app"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/utils/logs"
	"github.com/knative-sample/cloud-native-app-go/weather/pkg/utils/signals"
)

//go:generate protoc -I ../../pkg/city --go_out=plugins=grpc:../../pkg/city ../../pkg/city/city.proto
func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	stopCh := signals.SetupSignalHandler()

	// Start runner
	cmd := app.NewCommandStartServer(stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})

	if err := cmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
