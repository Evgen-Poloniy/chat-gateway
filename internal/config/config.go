package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// Acceptable logger levels.
const (
	TraceLevel = "trace"
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

// Acceptable logger formats.
const (
	TextFormat = "text"
	JsonFormat = "json"
	Stdout     = "stdout"
	Stderr     = "stderr"
)

// Config with tags from cleanenv library.
type ServerConfig struct {
	Host                    string        `env:"API_HOST" env-required:"true"`
	Port                    string        `env:"API_PORT" env-required:"true"`
	MaxHeaderBytes          int           `yaml:"max_header_bytes" env-default:"1048576"`
	ReadTimeout             time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout            time.Duration `yaml:"write_timeout" env-default:"10s"`
	TimeForGracefulShutdown time.Duration `yaml:"time_for_graceful_shutdown" env-default:"10s"`
	ReadHeaderTimeout       time.Duration `yaml:"read_header_timeout" env-default:"2s"`
	IdleTimeout             time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Logger config from config.yaml.
type LoggerConfig struct {
	Level  string   `yaml:"level" env-default:"info" validate:"oneof=trace debug info warn error panic fatal"`
	Format string   `yaml:"format" env-default:"json" validate:"oneof=text json"`
	Output string   `yaml:"output" env-default:"stdout" validate:"oneof=stdout stderr"`
	Files  []string `yaml:"files"`
}

// CORSConfig is the config for CORS policy.
type CORSConfig struct {
	AllowedOrigin    string        `yaml:"allowed_origin"`
	AllowCredentials bool          `yaml:"allow_credentials"`
	AllowedHeaders   []string      `yaml:"allowed_headers"`
	AllowedMethods   []string      `yaml:"allowed_methods"`
	MaxAge           time.Duration `yaml:"max_age" validate:"required,gte=1m,lte=24h"`
}

// Database config from env and config.yaml.
type DatabaseConfig struct {
	Host                string        `env:"DB_HOST" env-required:"true"`
	Port                string        `env:"DB_PORT" env-required:"true"`
	Username            string        `env:"DB_USER" env-required:"true"`
	Password            string        `env:"DB_PASSWORD" env-required:"true"`
	DBName              string        `env:"DB_NAME" env-required:"true"`
	SSLMode             string        `env:"SSL_MODE" env-required:"true" validate:"oneof=disable require"`
	MaxOpenConns        int           `yaml:"max_open_conns"`
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	ConnMaxLifetime     time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleLifetime time.Duration `yaml:"conn_max_idle_lifetime"`
}

// Message broker config from env and config.yaml.
type MessageBrokerConfig struct {
	BootstrapServers           string `env:"KAFKA_BOOTSTRAP_SERVERS" env-required:"true"`
	User                       string `env:"KAFKA_USER"`
	Password                   string `env:"KAFKA_PASSWORD"`
	SASLMechanism              string `env:"KAFKA_SASL_MECHANISM" validate:"oneof=PLAIN SCRAM-SHA-256 SCRAM-SHA-512"`
	SecurityProtocol           string `env:"KAFKA_SECURITY_PROTOCOL" validate:"oneof=PLAINTEXT SASL_PLAINTEXT SASL_SSL SSL"`
	Acks                       string `yaml:"acks"`
	EnableIdempotence          bool   `yaml:"enable_idempotence"`
	Retries                    int    `yaml:"retries"`
	MaxInFlightRequestsPerConn int    `yaml:"max_in_flight_requests_per_connection"`
	LingerMs                   int    `yaml:"linger_ms"`
	BatchNumMessages           int    `yaml:"batch_num_messages"`
	CompressionType            string `yaml:"compression_type" validate:"oneof=none gzip snpappy lz4 zstd"`
	QueueBufferingMaxMessages  int    `yaml:"queue_buffering_max_messages"`
	MessageTimeout             int    `yaml:"message_timeout"`
}

// Dataclass with all configs.
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Logger        LoggerConfig        `yaml:"logger"`
	CORS          CORSConfig          `yaml:"cors"`
	Database      DatabaseConfig      `yaml:"database"`
	MessageBroker MessageBrokerConfig `yaml:"message-broker"`
}

// Load config from config/config.yaml.
func LoadConfig(path string) (Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return Config{}, err
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return Config{}, err
	}

	return config, nil
}
