package shortener

type OriginURL struct {
	URL string `json:"url"`
}

type Shlink struct {
	Result string `json:"result,omitempty"`
}
