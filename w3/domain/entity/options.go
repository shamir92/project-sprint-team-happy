package entity

type SortType string

const (
	ASC  SortType = "asc"
	DESC SortType = "desc"
)

func (s SortType) String() string {
	return string(s)
}

type BrowseMedicalRecordPatientOptionBuilder func(*BrowseMedicalRecordPatientOption)
type BrowseMedicalRecordPatientOption struct {
	Limit          int
	Offset         int
	IdentityNumber *int
	Name           string
	PhoneNumber    *int
	SortCreatedAt  SortType
}

func WithOffsetAndLimit(offset, limit int) BrowseMedicalRecordPatientOptionBuilder {
	return func(b *BrowseMedicalRecordPatientOption) {
		b.Limit = limit
		b.Offset = offset
	}
}

func WithIdentityNumber(identityNumber int) BrowseMedicalRecordPatientOptionBuilder {
	return func(b *BrowseMedicalRecordPatientOption) {
		b.IdentityNumber = &identityNumber
	}
}

func WithName(name string) BrowseMedicalRecordPatientOptionBuilder {
	return func(b *BrowseMedicalRecordPatientOption) {
		b.Name = name
	}
}

func WithPhoneNumber(phoneNumber int) BrowseMedicalRecordPatientOptionBuilder {
	return func(b *BrowseMedicalRecordPatientOption) {
		b.PhoneNumber = &phoneNumber
	}
}

func WithSortCreatedAt(sortCreatedAt SortType) BrowseMedicalRecordPatientOptionBuilder {
	return func(b *BrowseMedicalRecordPatientOption) {
		b.SortCreatedAt = sortCreatedAt
	}
}
