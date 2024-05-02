package models

type CatRace int

type CatSex string

const (
	Persian CatRace = iota
	MaineCoon
	Siamese
	Ragdoll
	Bengal
	Sphynx
	BritishShorthair
	Abyssinian
	ScottishFold
	Birman
)

const (
	CatMale   CatSex = "male"
	CatFemale CatSex = "female"
)

func IsValidCatRace(race string) bool {
	var raceNames = map[string]CatRace{
		"Persian":           Persian,
		"Maine Coon":        MaineCoon,
		"Siamese":           Siamese,
		"Ragdoll":           Ragdoll,
		"Bengal":            Bengal,
		"Sphynx":            Sphynx,
		"British Shorthair": BritishShorthair,
		"Abyssinian":        Abyssinian,
		"Scottish Fold":     ScottishFold,
		"Birman":            Birman,
	}
	_, ok := raceNames[race]

	return ok
}

func IsValidCatSex(sex string) bool {
	switch sex {
	case string(CatMale), string(CatFemale):
		return true
	default:
		return false
	}
}
