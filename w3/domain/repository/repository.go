package repository

func whereOrAnd(paramsNumber int) string {
	if paramsNumber == 1 {
		return " WHERE "
	}

	return " AND "
}
