package model

type Url struct {
	AuthUrl  string
	TokenUrl string
}
type Config struct {
	RedirectUrl  string
	ClientID     string
	ClientSecret string
	Endpoint     Url
	Scopes       []string
}
