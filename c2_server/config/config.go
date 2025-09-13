package config

// ServerConfig holds the configuration for the C2 server
type ServerConfig struct {
	Host           string
	Port           string
	TLSEnabled     bool
	CertFile       string
	KeyFile        string
	AllowedOrigins []string
}

// DefaultConfig returns the default server configuration
func DefaultConfig() *ServerConfig {
	return &ServerConfig{
		Host:           "0.0.0.0",
		Port:           "8080",
		TLSEnabled:     false,
		CertFile:       "",
		KeyFile:        "",
		AllowedOrigins: []string{"*"}, // In production, this should be restricted
	}
}

// GetAddress returns the full address string (host:port)
func (c *ServerConfig) GetAddress() string {
	return c.Host + ":" + c.Port
}