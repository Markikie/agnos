package agnos

var Env struct {
	Port     string `env:"PORT" envDefault:"8080"`
	Database struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"agnos"`
		Password string `env:"DB_PASSWORD" envDefault:"password"`
		DBName   string `env:"DB_NAME" envDefault:"agnos"`
	} `envPrefix:"DB_"`
}
