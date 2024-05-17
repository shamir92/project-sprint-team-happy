package dto

type IdentityDetailDto struct {
	IdentityNumber      int64  `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type CreatedByDto struct {
	Nip    int64  `json:"nip"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type MedicalRecordItemDto struct {
	IdentityDetail IdentityDetailDto `json:"identityDetail"`
	Symptoms       string            `json:"symptoms"`
	Medications    string            `json:"medications"`
	CreatedAt      string            `json:"createdAt"`
	CreatedBy      CreatedByDto      `json:"createdBy"`
}
