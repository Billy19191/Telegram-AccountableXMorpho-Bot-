package config

// type Config struct {
// 	redisHost   string
// 	redisPort   int
// 	redisClient *redis.Client
// }

// func LoadConfig() (*Config, error) {
// 	redisHost := os.Getenv("REDIS_HOST")
// 	redisPort := os.Getenv("REDIS_PORT")

// 	if redisHost == "" {
// 		redisHost = "localhost"
// 	}

// 	if redisPort == "" {
// 		redisPort = "6379"
// 	}

// 	portNum, err := strconv.Atoi(redisPort)
// 	if err != nil {
// 		return nil, fmt.Errorf("REDIS_PORT must be a number: %w", err)
// 	}

// 	c := &Config{
// 		redisHost: redisHost,
// 		redisPort: portNum,
// 	}

// 	client, err := c.initRedisCache()
// 	if err != nil {
// 		return nil, err
// 	}
// 	c.redisClient = client

// 	return c, nil
// }

// func (c *Config) initRedisCache() (*redis.Client, error) {

// 	addr := fmt.Sprintf("%s:%d", c.redisHost, c.redisPort)
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     addr,
// 		Password: "",
// 		DB:       0,
// 	})

// 	ctx := context.Background()

// 	_, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
// 	}

// 	log.Println("✅ Redis initialized!")

// 	return client, nil
// }
