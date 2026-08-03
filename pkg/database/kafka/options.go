package kafka

type Config struct {
	BootstrapServers []string
	Acks             string
	Retries          int
	LingerMs         int
	BatchNumMessages int
	CompressionType  string
	MessageTimeoutMs int

	SecurityProtocol string
	SASLMechanism    string
	SASLUsername     string
	SASLPassword     string
}

type Option func(cfg *Config) error

// WithSecurityProtocol sets the communication protocol (e.g., SASL_SSL, SASL_PLAINTEXT).
// It ignores empty strings or standard unencrypted PLAINTEXT.
func WithSecurityProtocol(protocol string) Option {
	return func(cfg *Config) error {
		if protocol == "" || protocol == "PLAINTEXT" {
			return nil
		}
		cfg.SecurityProtocol = protocol
		return nil
	}
}

// WithSASLMechanism sets the SASL authentication mechanism (e.g., PLAIN, SCRAM-SHA-512).
func WithSASLMechanism(mechanism string) Option {
	return func(cfg *Config) error {
		if mechanism == "" {
			return nil
		}
		cfg.SASLMechanism = mechanism
		return nil
	}
}

// WithSASLUsername sets the username for SASL authentication.
func WithSASLUsername(username string) Option {
	return func(cfg *Config) error {
		if username == "" {
			return nil
		}
		cfg.SASLUsername = username
		return nil
	}
}

// WithSASLPassword sets the password for SASL authentication.
func WithSASLPassword(password string) Option {
	return func(cfg *Config) error {
		if password == "" {
			return nil
		}
		cfg.SASLPassword = password
		return nil
	}
}
