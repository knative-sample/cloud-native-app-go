package weather

type WebApi struct {
	Port          string
	CityService   *Service
	DetailService *Service
	ResourceRoot  string
}

type Service struct {
	Host string
	Port string
}
