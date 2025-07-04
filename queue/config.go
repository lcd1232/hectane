package queue

import "time"

// See https://github.com/Freeaqingme/dkim
type DKIMConfig struct {
	PrivateKey       string `json:"private-key"`
	Selector         string `json:"selector"`
	Canonicalization string `json:"canonicalization"`
}

// Application configuration.
type Config struct {
	// TODO: support json marshal for duration
	Hostname               string        `json:"hostname"`
	Directory              string        `json:"directory"`
	DisableSSLVerification bool          `json:"disable-ssl-verification"`
	InitialInterval        time.Duration `json:"-"`
	RandomizationFactor    float64       `json:"randomization_factor"`
	Multiplier             float64       `json:"multiplier"`
	MaxInterval            time.Duration `json:"-"`
	MaxAttempts            int           `json:"max_attempts"`
	Jitter                 bool          `json:"jitter"`
	// MessageTTL is the time to live for messages when we're trying to load them.
	MessageTTL time.Duration `json:"-"`

	// Map domain names to DKIM config for that domain
	DKIMConfigs map[string]DKIMConfig `json:"dkim-configs"`
	// ProcessFunc allow you to define custom process function for message.
	ProcessFunc ProcessFunc `json:"-"`
	// ProcessFailFunc allow you to define processor if message failed with 5xx status.
	ProcessFailFunc ProcessFailFunc `json:"-"`
}
