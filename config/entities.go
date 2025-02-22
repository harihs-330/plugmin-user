package config

// EnvConfig represents the configuration structure for the application.
type EnvConfig struct {
	Debug             bool     `default:"true" split_words:"true"`
	Port              int      `default:"8080" split_words:"true"`
	DB                Database `split_words:"true"`
	AcceptedVersions  []string `required:"true" split_words:"true"`
	Gmail             Gmail    `split_words:"true"`
	SignupURL         string   `required:"true" split_words:"true"`
	UserURL           string   `required:"true" split_words:"true"`
	SecretKey         string   `required:"true" split_words:"true"`
	AddPermissionsURL string   `required:"true" split_words:"true"`
	ResetPasswordURL  string   `required:"true" split_words:"true"`
	APIBuilderURL     string   `required:"true" split_words:"true"`
}

// Database represents the configuration for the database connection.
type Database struct {
	Driver    string
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}

type Gmail struct {
	Name     string
	Address  string
	Password string
}
