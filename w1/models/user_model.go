package models

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetAllUsers() []User {
	// Normally, you would query the database here
	return []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
}
