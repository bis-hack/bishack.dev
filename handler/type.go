package handler

type githubUser struct {
	Email     string
	Name      string
	Login     string
	AvatarURL string `json:"avatar_url"`
	Location  string
	Website   string `json:"blog"`
}
