package dto

type PatientRegisterControllerResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PatientBrowseControllerResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
