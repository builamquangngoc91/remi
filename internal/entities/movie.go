package entities

import "time"

// Movie reflects movies data from DB
type Movie struct {
	ID          string
	Name        string
	Description string
	Link        string
	Thumbnail   string
	SharedBy    string
	SharedAt    *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type Movies []*Movie

func (e *Movie) FieldMap() (fields []string, values []interface{}) {
	return []string{
			"id",
			"name",
			"description",
			"link",
			"thumbnail",
			"shared_by",
			"shared_at",
			"created_at",
			"updated_at",
			"deleted_at",
		}, []interface{}{
			&e.ID,
			&e.Name,
			&e.Description,
			&e.Link,
			&e.Thumbnail,
			&e.SharedBy,
			&e.SharedAt,
			&e.CreatedAt,
			&e.UpdatedAt,
			&e.DeletedAt,
		}
}

func (e *Movie) TableName() string {
	return "movies"
}
