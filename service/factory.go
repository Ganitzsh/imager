package service

var tokenStoreFactory = map[StoreType]func(c *Config) (TokenStore, error){
	StoreTypeRedis: func(c *Config) (TokenStore, error) {
		client, err := connectToRedis(c)
		if err != nil {
			return nil, err
		}
		return NewTokenStoreRedis(client)
	},
}

func GetTokenStore(c *Config) (TokenStore, error) {
	if c == nil || c.Store == nil {
		return nil, ErrInvalidConfig
	}
	f := tokenStoreFactory[c.Store.Type]
	if f == nil {
		return nil, ErrInvalidConfig
	}
	return f(c)
}
