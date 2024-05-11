package entity

type SortType string

const (
	DESC SortType = "desc"
	ASC  SortType = "asc"
)

func (s SortType) String() string {
	return string(s)
}

type FindProductOptionBuilder func(*FindProductOption)

type FindProductOption struct {
	Limit         int
	Offset        int
	IsAvailable   *bool
	InStock       *bool
	ID            string
	Name          string
	Category      string
	SKU           string
	SortPrice     SortType
	SortCreatedAt SortType
}

func WithOffsetAndLimit(offset, limit int) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.Limit = limit
		fpo.Offset = offset
	}
}

func WithIsAvailable(isAvailable *bool) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.IsAvailable = isAvailable
	}
}

func WithInStock(inStock *bool) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.InStock = inStock
	}
}

func WithProductID(id string) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.ID = id
	}
}

func WithProductName(name string) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.Name = name
	}
}

func WithProductCategory(category string) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.Category = category
	}
}

func WithProductSKU(sku string) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.SKU = sku
	}
}

func WithSortPrice(sort SortType) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.SortPrice = sort
	}
}

func WithSortCreatedAt(sort SortType) FindProductOptionBuilder {
	return func(fpo *FindProductOption) {
		fpo.SortCreatedAt = sort
	}
}

type FindCheckoutHistoryBuilder func(*FindCheckoutHistory)

type FindCheckoutHistory struct {
	CustomerID    string
	Limit         int
	Offset        int
	SortCreatedAt SortType
}

func FindCheckoutWithLimitAndOffset(offset, limit int) FindCheckoutHistoryBuilder {
	return func(fch *FindCheckoutHistory) {
		fch.Offset = offset
		fch.Limit = limit
	}
}

func FindCheckoutWithCustomerId(customerId string) FindCheckoutHistoryBuilder {
	return func(fch *FindCheckoutHistory) {
		fch.CustomerID = customerId
	}
}

func FindCheckoutWithSortByCreatedAt(sort SortType) FindCheckoutHistoryBuilder {
	return func(fch *FindCheckoutHistory) {
		fch.SortCreatedAt = sort
	}
}
