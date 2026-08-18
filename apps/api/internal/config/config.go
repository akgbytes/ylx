package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Email    EmailConfig
	Log      LogConfig
}

func Load() (*Config, error) {
	server, err := loadServerConfig()
	if err != nil {
		return nil, err
	}

	database, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	redis, err := loadRedisConfig()
	if err != nil {
		return nil, err
	}

	auth, err := loadAuthConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server:   server,
		Database: database,
		Redis:    redis,
		Auth:     auth,
		Email:    loadEmailConfig(),
		Log:      loadLogConfig(),
	}

	cfg.normalize()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) normalize() {
	c.Server.normalize()
	c.Database.normalize()
	c.Redis.normalize()
	c.Auth.normalize()
	c.Email.normalize()
	c.Log.normalize()
}

func (c *Config) validate() error {
	if err := c.Server.validate(); err != nil {
		return err
	}

	if err := c.Database.validate(); err != nil {
		return err
	}

	if err := c.Redis.validate(); err != nil {
		return err
	}

	if err := c.Auth.validate(); err != nil {
		return err
	}

	if err := c.Email.validate(); err != nil {
		return err
	}

	if err := c.Log.validate(); err != nil {
		return err
	}

	return nil
}
