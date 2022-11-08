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
	"github.com/stretchr/testify/assert"
)

func TestMovieRepository_Create(t *testing.T) {
	db, mock := NewMock()
	repo := MovieRepository{DB: db}

	now := time.Now()
	m := &entities.Movie{
		ID:          idutil.NewID(),
		Name:        "movie-1",
		Description: "description of movie-1",
		Link:        "link of movie-1",
		Thumbnail:   "thumbnail of movie-1",
		SharedBy:    "1",
		SharedAt:    &now,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		DeletedAt:   nil,
	}

	testCases := []TestCase{
		{
			name:        "happy case",
			req:         m,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO movies(id,name,description,link,thumbnail,shared_by,shared_at,created_at,updated_at,deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")).
					WithArgs(m.ID, m.Name, m.Description, m.Link, m.Thumbnail, m.SharedBy, m.SharedAt, m.CreatedAt, m.UpdatedAt, m.DeletedAt).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:        "exec error",
			req:         m,
			expectedErr: fmt.Errorf("r.DB.ExecContext: %w", sql.ErrNoRows),
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO movies(id,name,description,link,thumbnail,shared_by,shared_at,created_at,updated_at,deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")).
					WithArgs(m.ID, m.Name, m.Description, m.Link, m.Thumbnail, m.SharedBy, m.SharedAt, m.CreatedAt, m.UpdatedAt, m.DeletedAt).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "no row affected",
			req:         m,
			expectedErr: fmt.Errorf("can't insert movie"),
			setup: func(ctx context.Context) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO movies(id,name,description,link,thumbnail,shared_by,shared_at,created_at,updated_at,deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")).
					WithArgs(m.ID, m.Name, m.Description, m.Link, m.Thumbnail, m.SharedBy, m.SharedAt, m.CreatedAt, m.UpdatedAt, m.DeletedAt).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		err := repo.Create(ctx, testCase.req.(*entities.Movie))
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
		}
	}
}

func TestMovieRepository_FindByIDAndUserID(t *testing.T) {
	db, mock := NewMock()
	repo := MovieRepository{DB: db}

	type Args struct {
		ID     string
		UserID string
	}
	args := &Args{
		ID:     "id",
		UserID: "user-id",
	}

	testCases := []TestCase{
		{
			name:        "happy case",
			req:         args,
			expectedErr: nil,
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,description,link,thumbnail,shared_by,shared_at,created_at,updated_at,deleted_at FROM movies WHERE id = $1 AND shared_by = $2")).
					WithArgs(args.ID, args.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "link", "thumbnail", "shared_by", "shared_at", "created_at", "updated_at", "deleted_at"}).AddRow(idutil.NewID(), "name", "description", "link", "thumbnail", "1", time.Now(), time.Now(), time.Now(), nil))
			},
		},
		{
			name:        "exec error",
			req:         args,
			expectedErr: fmt.Errorf("row.Scan: %w", sql.ErrNoRows),
			setup: func(ctx context.Context) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,description,link,thumbnail,shared_by,shared_at,created_at,updated_at,deleted_at FROM movies WHERE id = $1 AND shared_by = $2")).
					WithArgs(args.ID, args.UserID).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		testCase.setup(ctx)
		args := testCase.req.(*Args)
		movie, err := repo.FindByIDAndUserID(ctx, args.ID, args.UserID)
		if testCase.expectedErr != nil {
			assert.Equal(t, testCase.expectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, testCase.expectedErr, err)
			assert.NotNil(t, movie)
		}
	}
}
