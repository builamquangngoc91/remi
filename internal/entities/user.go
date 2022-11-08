package entities

import "time"

// User reflects users data from DB
type User struct {
	ID        string
	Username  string
	Password  string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Users []*User

func (e *User) FieldMap() (fields []string, values []interface{}) {
	return []string{
			"id",
			"username",
			"password",
			"name",
			"created_at",
			"updated_at",
		}, []interface{}{
			&e.ID,
			&e.Username,
			&e.Password,
			&e.Name,
			&e.CreatedAt,
			&e.UpdatedAt,
		}
}

func (e *User) TableName() string {
	return "users"
}
