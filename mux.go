package smux

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Config is used to tune the Smux session
type Config struct {
	// KeepAliveInterval is how often to send a NOP command to the remote
	KeepAliveInterval time.Duration

	// KeepAliveTimeout is how long the session
	// will be closed if no data has arrived
	KeepAliveTimeout time.Duration

	// MaxFrameSize is used to control the maximum
	// frame size to sent to the remote
	MaxFrameSize int

	// MaxFrameTokens is used to control the maximum
	// number of frame in the buffer pool
	MaxFrameTokens int
}

// DefaultConfig is used to return a default configuration
func DefaultConfig() *Config {
	return &Config{
		KeepAliveInterval: 10 * time.Second,
		KeepAliveTimeout:  30 * time.Second,
		MaxFrameSize:      2048,
		MaxFrameTokens:    2048,
	}
}

// VerifyConfig is used to verify the sanity of configuration
func VerifyConfig(config *Config) error {
	if config.KeepAliveInterval == 0 {
		return errors.New("keep-alive interval must be positive")
	}
	if config.KeepAliveTimeout < config.KeepAliveInterval {
		return fmt.Errorf("keep-alive timeout must be larger than keep-alive interval")
	}
	if config.MaxFrameSize <= 0 {
		return errors.New("max frame size must be positive")
	}
	if config.MaxFrameSize > 65535 {
		return errors.New("max frame size must not be larger than 65535")
	}
	if config.MaxFrameTokens <= 0 {
		return errors.New("max frame tokens must be positive")
	}
	return nil
}

// Server is used to initialize a new server-side connection.
func Server(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := VerifyConfig(config); err != nil {
		return nil, err
	}
	return newSession(config, conn, false), nil
}

// Client is used to initialize a new client-side connection.
func Client(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := VerifyConfig(config); err != nil {
		return nil, err
	}
	return newSession(config, conn, true), nil
}
