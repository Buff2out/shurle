package shortener

type OriginURL struct {
	URL string `json:"url"`
}

type Shlink struct {
	Result string `json:"result,omitempty"`
}

type ShURLFile struct {
	UID         string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Settings struct {
	Socket     string `env:"SERVER_ADDRESS,required"`
	Prefix     string `env:"BASE_URL,required"`
	ShURLsJSON string `env:"FILE_STORAGE_PATH,required"`
}
