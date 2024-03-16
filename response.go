package main

type UserResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

type LoginResponse struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	Id           int    `json:"id"`
}

type ChirpResponse struct {
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
	Id       int    `json:"id"`
}
