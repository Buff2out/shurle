package server

type Settings struct {
	Socket      string `env:"SERVER_ADDRESS,required"`
	Prefix      string `env:"BASE_URL,required"`
	ShURLsJSON  string `env:"FILE_STORAGE_PATH,required"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}
