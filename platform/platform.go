package platform

type DeploymentPlatform interface {
	CreateApp() (string, error)
	CreateAddOns(serviceImages []string) error
	Deploy() error
}
