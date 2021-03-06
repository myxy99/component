package xemail

type options struct {
	Host     string `json:"host,omitempty" mapStructure:"host"`
	Port     int    `json:"port" mapStructure:"port"`
	Username string `json:"username,omitempty" mapStructure:"username"`
	Password string `json:"-" mapStructure:"password"`
}

func newEmailOptions() *options {
	return &options{
		Host:     "smtp.yeah.net",
		Port:     465,
		Username: "yangzzzzzz@yeah.net",
		Password: "your-password",
	}
}
