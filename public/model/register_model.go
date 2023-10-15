package public_model

type RegisterModel struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
