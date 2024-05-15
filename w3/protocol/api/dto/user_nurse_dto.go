package dto

type CreateUserNurseDtoResponse struct {
	UserID string `json:"userId"`
	NIP    string `json:"nip"`
	Name   string `json:"name"`
}
