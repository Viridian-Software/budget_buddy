package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) AddUserHandler(w http.ResponseWriter, r *http.Request) {

}
