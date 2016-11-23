package defaults

import "time"

const (
	// RetryDelay defines the interval between retry attempts
	RetryDelay = 5 * time.Second
	// RetryAttempts defines the maximum number of retry attempts
	RetryAttempts = 100

	// SSHConnectTimeout defines the timeout for establishing an SSH connection
	SSHConnectTimeout = 20 * time.Second
)
