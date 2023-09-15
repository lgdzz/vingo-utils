package config

type System struct {
	Service struct {
		Name      string `yaml:"name" json:"name"`
		Debug     bool   `yaml:"debug" json:"debug"`
		Port      uint   `yaml:"port" json:"port"`
		Copyright string `yaml:"copyright" json:"copyright"`
	} `yaml:"service" json:"service"`
	Super struct {
		Enable   bool   `yaml:"enable" json:"enable"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"super" json:"super"`
	Auth struct {
		Secret string `yaml:"secret" json:"secret"`
		SSO    bool   `yaml:"sso" json:"sso"`
		Log    bool   `yaml:"log" json:"log"`
		Lock   struct {
			Enable bool   `yaml:"enable" json:"enable"`
			Ticket string `yaml:"ticket" json:"ticket"`
			Bad    uint8  `yaml:"bad" json:"bad"`
			Time   uint   `yaml:"time" json:"time"`
		} `yaml:"lock" json:"lock"`
		Strength int `yaml:"strength" json:"strength"`
	} `yaml:"auth" json:"auth"`
	Right struct {
		Enable bool   `yaml:"enable" json:"enable"`
		Ticket string `yaml:"ticket" json:"ticket"`
	} `yaml:"right" json:"right"`
	Account struct {
		Many bool `yaml:"many" json:"many"`
	} `yaml:"account" json:"account"`
	Cli bool `yaml:"cli" json:"cli"`
}
