package server

type ServiceInfo struct {
	Name string
	ServerIp string
	Ports []Port
}

type Port struct {
	Protocol string
	Name string
	Url string
	Target string
}
