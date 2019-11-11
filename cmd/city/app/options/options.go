package options

import (
	"github.com/spf13/cobra"
)

type Options struct {
	ResourceRoot  string
	CityService   string
	DetailService string
	Port          string
	Version       bool
}

func (s *Options) SetOps(ac *cobra.Command) {
	ac.Flags().StringVar(&s.ResourceRoot, "resource", "/var/html", "html resource root")
	ac.Flags().StringVar(&s.CityService, "city-service", "city:8080", "application Name")
	ac.Flags().StringVar(&s.DetailService, "detail-service", "detail:8080", "Application so File Path")
	ac.Flags().StringVar(&s.Port, "port", "8080", "http listen port")
	ac.Flags().BoolVar(&s.Version, "version", false, "Print version information")
}
