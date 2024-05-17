package dto

import "time"

type UserITRegisterControllerResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type UserITLoginControllerResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ListUserItemDto struct {
	ID        string    `json:"userId"`
	Name      string    `json:"name"`
	NIP       string    `json:"nip"`
	CreatedAt time.Time `json:"createdAt"`
}
