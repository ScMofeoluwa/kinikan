package platform

type DeploymentPlatform interface {
	ValidateCLI() error
	CreateApp() error
	CreateAddOns(serviceImages []string) error
	DeployApp(serviceImages []string) error
}
