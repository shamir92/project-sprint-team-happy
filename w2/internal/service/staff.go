package service

type Staff struct {
	ID          string
	PhoneNumber string
	Name        string
	Password    string
}

type CreateStafIn struct {
	PhoneNumber string
	Name        string
	Password    string
}

func (s *Service) CreateStaff(in CreateStafIn) (Staff, error) {
	insertUserQuery := `
		INSERT INTO users(phone_number, name, password) 
		VALUES($1, $2, $3, $4) 
		RETURNING id`

	staff := Staff{
		PhoneNumber: in.Password,
		Name:        in.Name,
	}

	err := s.db.QueryRow(insertUserQuery, in.Name, in.PhoneNumber, in.Password).Scan(&staff.ID)

	if err != nil {
		return Staff{}, err
	}

	return staff, nil
}
