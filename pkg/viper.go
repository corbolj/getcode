package getcode

type ViperConfig struct {
	DefaultUser      string              `yaml:"default_user"`
	GitCredentalPath string              `yaml:"git_credental_path"`
	OrgUsers         map[string][]string `yaml:"org_users"`
	ProjectPath      string              `yaml:"project_path"`
}
