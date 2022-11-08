package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"remi/internal/entities"
	"remi/pkg/golibs/idutil"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock := NewMock()
	repo := UserRepository{DB: db}

	now := time.Now()
	u := &entities.User{
		ID:        idutil.NewID(),
		Name:      "name",
		Username:  "username",
		Password:  "password",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	testCases := []TestCase{
		{
			name:        "happy case",
			req:         u,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(id,username,password,name,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)")).
					WithArgs(u.ID, u.Username, u.Password, u.Name, u.CreatedAt, u.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:        "exec error",
			req:         u,
			expectedErr: fmt.Errorf("r.DB.ExecContext: %w", sql.ErrNoRows),
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(id,username,password,name,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)")).
					WithArgs(u.ID, u.Username, u.Password, u.Name, u.CreatedAt, u.UpdatedAt).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "no row affected",
			req:         u,
			expectedErr: fmt.Errorf("can't insert user"),
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(id,username,password,name,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)")).
					WithArgs(u.ID, u.Username, u.Password, u.Name, u.CreatedAt, u.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		err := repo.Create(ctx, testCase.req.(*entities.User))
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
		}
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock := NewMock()
	repo := UserRepository{DB: db}

	arg := "name"
	testCases := []TestCase{
		{
			name:        "happy case",
			req:         arg,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE username = $1")).
					WithArgs(arg).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "created_at", "updated_at"}).AddRow(idutil.NewID(), "username", "password", "name", time.Now(), time.Now()))
			},
		},
		{
			name:        "exec error",
			req:         arg,
			expectedErr: sql.ErrNoRows,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE username = $1")).
					WithArgs(arg).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		movie, err := repo.FindByUsername(ctx, testCase.req.(string))
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
			assert.NotNil(t, movie)
		}
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := NewMock()
	repo := UserRepository{DB: db}

	arg := "id"
	testCases := []TestCase{
		{
			name:        "happy case",
			req:         arg,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE id = $1")).
					WithArgs(arg).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "created_at", "updated_at"}).AddRow(idutil.NewID(), "username", "password", "name", time.Now(), time.Now()))
			},
		},
		{
			name:        "exec error",
			req:         arg,
			expectedErr: sql.ErrNoRows,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE id = $1")).
					WithArgs(arg).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		movie, err := repo.FindByID(ctx, testCase.req.(string))
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
			assert.NotNil(t, movie)
		}
	}
}

func TestUserRepository_List(t *testing.T) {
	db, mock := NewMock()
	repo := UserRepository{DB: db}

	args := &ListUsersArgs{
		IDs: []string{"id-1", "id-2"},
	}

	testCases := []TestCase{
		{
			name:        "happy case",
			req:         args,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE id = ANY($1::_TEXT)")).
					WithArgs(pq.StringArray(args.IDs)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "created_at", "updated_at"}).AddRow(idutil.NewID(), "username", "password", "name", time.Now(), time.Now()))
			},
		},
		{
			name:        "exec error",
			req:         args,
			expectedErr: fmt.Errorf("r.QueryContext: %w", sql.ErrNoRows),
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,username,password,name,created_at,updated_at FROM users WHERE id = ANY($1::_TEXT)")).
					WithArgs(pq.StringArray(args.IDs)).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		movie, err := repo.List(ctx, testCase.req.(*ListUsersArgs))
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
			assert.NotNil(t, movie)
		}
	}
}
