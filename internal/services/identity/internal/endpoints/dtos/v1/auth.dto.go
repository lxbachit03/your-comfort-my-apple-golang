package v1dto

type LoginAccountRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterAccountRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
}

type LoginAccountResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
	ExpiresIn    int    `json:"expires_in"`
}

type RegisterAccountResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
