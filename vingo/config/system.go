package config

type System struct {
	Service SystemService `yaml:"service" json:"service"`
	Super   SystemSuper   `yaml:"super" json:"super"`
	Auth    SystemAuth    `yaml:"auth" json:"auth"`
	Right   SystemRight   `yaml:"right" json:"right"`
	Account SystemAccount `yaml:"account" json:"account"`
	Cli     bool          `yaml:"cli" json:"cli"`
}

type SystemService struct {
	Name      string `yaml:"name" json:"name"`
	Debug     bool   `yaml:"debug" json:"debug"`
	Port      uint   `yaml:"port" json:"port"`
	Copyright string `yaml:"copyright" json:"copyright"`
}

type SystemSuper struct {
	Enable   bool   `yaml:"enable" json:"enable"`
	Password string `yaml:"password" json:"password"`
}

type SystemAuth struct {
	Secret   string         `yaml:"secret" json:"secret"`
	SSO      bool           `yaml:"sso" json:"sso"`
	Log      bool           `yaml:"log" json:"log"`
	Lock     SystemAuthLock `yaml:"lock" json:"lock"`
	Strength int            `yaml:"strength" json:"strength"`
}

type SystemAuthLock struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Ticket string `yaml:"ticket" json:"ticket"`
	Bad    uint8  `yaml:"bad" json:"bad"`
	Time   uint   `yaml:"time" json:"time"`
}

type SystemRight struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Ticket string `yaml:"ticket" json:"ticket"`
}

type SystemAccount struct {
	Many bool `yaml:"many" json:"many"`
}
