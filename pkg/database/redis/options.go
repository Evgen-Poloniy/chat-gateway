package redis

import "time"

// Options holds configuration parameters for the Redis connection.
type Options struct {
	addr         string
	password     string
	db           int
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Option defines a function type used to configure the Options struct.
type Option func(*Options)

// WithAddr configures the Redis network address combining host and port.
func WithAddr(host, port string) Option {
	return func(o *Options) {
		o.addr = host + ":" + port
	}
}

// WithPassword sets the password used for Redis authentication.
func WithPassword(password string) Option {
	return func(o *Options) {
		o.password = password
	}
}

// WithDB specifies the target Redis database index.
func WithDB(db int) Option {
	return func(o *Options) {
		o.db = db
	}
}

// WithDialTimeout sets the maximum amount of time to establish a connection.
func WithDialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = d
	}
}

// WithReadTimeout sets the timeout duration for socket reads.
func WithReadTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = d
	}
}

// WithWriteTimeout sets the timeout duration for socket writes.
func WithWriteTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = d
	}
}
