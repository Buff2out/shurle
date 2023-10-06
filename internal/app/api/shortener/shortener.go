package shortener

type OriginURL struct {
	URL string `json:"url"`
}

type Shlink struct {
	Result string `json:"result,omitempty"`
}

type ShURLFile struct {
	UId         string `json:"uid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
