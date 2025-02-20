package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `env:"ENV" env-required:"true"`
	StoragePath string `env:"STORAGE_PATH" env-required:"true"`
	SecretKey   string `env:"SECRET_KEY" env-required:"true"`
	HTTPServer  `env-prefix:"HTTP_"`
	Services    `env-prefix:"SERVICES_"`
	Redis       `env-prefix:"REDIS_"`
}

type HTTPServer struct {
	Address           string        `env:"ADDRESS"`
	Port              string        `env:"API_PORT"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT"`
}

type Services struct {
	AddressGRPC string `env:"ADDRESS_GRPC_SERVER"`
	PortGRPC    string `env:"PORT_GRPC_SERVER"`
}

type Redis struct {
	Address  string        `env:"HOST"`
	Port     string        `env:"PORT"`
	Password string        `env:"PASSWORD"`
	TTLKeys  time.Duration `env:"TTL_KEYS"`
}

// MustLoad loads the configuration from a file specified via a command-line flag.
// If the configuration file does not exist or an error occurs while reading it, the program terminates with a fatal error.
// Upon successful loading of the configuration, the function returns a pointer to the Config struct.
func MustLoad() *Config {
	var cfg Config
	var filePath string

	flag.StringVar(&filePath, "config", "", "path to config file")
	flag.Parse()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("env file does not exist: %s", filePath)
	}

	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	log.Println("configuration file successfully loaded")
	return &cfg
}
