package dto

type CreateUserNurseDtoResponse struct {
	UserID string `json:"userId"`
	NIP    int    `json:"nip"`
	Name   string `json:"name"`
}

type NurseLoginDtoResponse struct {
	UserID      string `json:"userId"`
	NIP         int    `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
