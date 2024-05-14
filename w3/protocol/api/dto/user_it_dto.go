package dto

type UserITRegisterRequest struct {
	NIP      string `json:"nip" validate:"required, numeric, min=6150000000000, max=6159999999999"`
	Name     string `json:"name" validate:"required, min=5, max=50"`
	Password string `json:"password" validate:"required, min=5, max=33"`
}
