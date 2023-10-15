package public_model

type TokenModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshModel struct {
	Token string `json:"token"`
}
