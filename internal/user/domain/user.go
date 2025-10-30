package domain

import "time"

type User struct {
	ID           string
	Email        string
	AvatarURL    string
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
	Provider  string
	ID        string
	Login     string
	Nickname  string
	Email     string
	AvatarURL string
}
