package main

type UserResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

type ChirpResponse struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}
