package service

type Staff struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"-"`
}

type CreateStaffRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

func (s *Service) StaffCreate(in CreateStaffRequest) (Staff, error) {
	insertUserQuery := `
		INSERT INTO users(phone_number, name, password) 
		VALUES($1, $2, $3) 
		RETURNING user_id`

	staff := Staff{
		PhoneNumber: in.PhoneNumber,
		Name:        in.Name,
	}

	err := s.db.QueryRow(insertUserQuery, in.Name, in.PhoneNumber, in.Password).Scan(&staff.UserID)

	if err != nil {
		return Staff{}, err
	}

	return staff, nil
}
