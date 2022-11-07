package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"remi/internal/entities"
	"remi/pkg/golibs/database"

	"github.com/lib/pq"
)

type UserRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *entities.User) error {
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
		return fmt.Errorf("can't insert user")
	}

	return err
}

// FindByUsername find user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	user := &entities.User{}
	fields, values := user.FieldMap()

	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE username = $1`, strings.Join(fields, ","), user.TableName())
	row := r.QueryRowContext(ctx, stmt, username)

	if err := row.Scan(values...); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID find user by id
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	user := &entities.User{}
	fields, values := user.FieldMap()

	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, strings.Join(fields, ","), user.TableName())
	row := r.QueryRowContext(ctx, stmt, id)

	if err := row.Scan(values...); err != nil {
		return nil, err
	}

	return user, nil
}

type ListUsersArgs struct {
	IDs []string
}

// List find movies
func (r *UserRepository) List(ctx context.Context, args *ListUsersArgs) (us entities.Users, _ error) {
	user := &entities.User{}
	fields, _ := user.FieldMap()

	stmt := fmt.Sprintf(`SELECT %s FROM %s 
	WHERE id = ANY($1::_TEXT)`, strings.Join(fields, ","), user.TableName())
	rows, err := r.QueryContext(ctx, stmt, pq.StringArray(args.IDs))
	if err != nil {
		return nil, fmt.Errorf("r.QueryContext: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		u := &entities.User{}
		_, values := u.FieldMap()
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		us = append(us, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return us, nil
}
