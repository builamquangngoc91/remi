package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"remi/internal/entities"
	"remi/pkg/golibs/database"
)

type MovieRepository struct {
	*sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db,
	}
}

func (r *MovieRepository) Create(ctx context.Context, u *entities.Movie) error {
	fields, values := u.FieldMap()
	placeHolders := database.GeneratePlaceholders(len(fields))

	stmt := fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s)`, u.TableName(), strings.Join(fields, ","), placeHolders)
	result, err := r.DB.ExecContext(ctx, stmt, values...)
	if err != nil {
		return fmt.Errorf("r.DB.ExecContext: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rowAffected != 1 {
		return fmt.Errorf("can't insert movie")
	}

	return err
}

// Get find movie by id
func (r *MovieRepository) FindByIDAndUserID(ctx context.Context, id, userID string) (*entities.Movie, error) {
	movie := &entities.Movie{}
	fields, values := movie.FieldMap()

	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1 AND shared_by = $2 AND deleted_at IS NULL`, strings.Join(fields, ","), movie.TableName())
	row := r.QueryRowContext(ctx, stmt, id, userID)

	if err := row.Scan(values...); err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return movie, nil
}

type ListMoviesArgs struct {
	UserID *string
	Offset *int
	Limit  *int
}

// List find movies
func (r *MovieRepository) List(ctx context.Context, args *ListMoviesArgs) (ms entities.Movies, _ error) {
	movie := &entities.Movie{}
	fields, _ := movie.FieldMap()

	limit := 10
	if args.Limit != nil {
		limit = *args.Limit
	}

	offset := 0
	if args.Offset != nil {
		offset = *args.Offset
	}

	stmt := fmt.Sprintf(`SELECT %s FROM %s 
	WHERE ($1::TEXT IS NULL OR shared_by = $1::TEXT) AND
	deleted_at IS NULL
	ORDER BY created_at DESC
	LIMIT %d
	OFFSET %d`, strings.Join(fields, ","), movie.TableName(), limit, offset)
	rows, err := r.QueryContext(ctx, stmt, args.UserID)
	if err != nil {
		return nil, fmt.Errorf("r.QueryContext: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		m := &entities.Movie{}
		_, values := m.FieldMap()
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		ms = append(ms, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return ms, nil
}
