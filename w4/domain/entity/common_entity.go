package entity

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

func (st SortType) String() string {
	return string(st)
}

func (st SortType) Valid() bool {
	return st == SortTypeAsc || st == SortTypeDesc
}

type Location struct {
	Lat float64 `json:"lat" validate:"required,latitude"`
	Lon float64 `json:"long" validate:"required,longitude"`
}
