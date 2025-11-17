package domain

import "time"

type OAuthProvider string

func (o OAuthProvider) String() string {
	return string(o)
}

const OAuthProviderGithub = "github"

type User struct {
	ID           string
	Email        string
	Avatar       string
	PasswordHash string
	Nickname     string
	GithubID     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  time.Time
}

type JwtPayload struct {
	UserID string `json:"user_id"`
}

type User2Token struct {
	AccessToken  string
	RefreshToken string
}

type OAuthUserInfo struct {
	Provider string
	ID       string
	Login    string
	Nickname string
	Email    string
}
