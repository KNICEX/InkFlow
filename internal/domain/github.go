package domain

type GithubInfo struct {
	Id           int64
	Username     string
	AccessToken  string
	RefreshToken string
	AvatarUrl    string
}
