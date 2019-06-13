package handler

type githubUser struct {
	Bio       string
	Name      string
	Email     string
	Login     string
	Website   string `json:"blog"`
	Location  string
	AvatarURL string `json:"avatar_url"`
}
