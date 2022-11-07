package up

import (
	"strings"
	"time"

	"remi/pkg/xerror"
)

type CreateMovieRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

func (r *CreateMovieRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "name can't be null")
	}
	if strings.TrimSpace(r.Link) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "link can't be null")
	}
	if strings.TrimSpace(r.Description) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "description can't be null")
	}

	return nil
}

type CreateMovieResponse struct {
	ID string `json:"id"`
}

type GetMovieByUserRequest struct {
	ID string `json:"id"`
}

func (r *GetMovieByUserRequest) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "id can't be null")
	}

	return nil
}

type GetMovieByUserResponse struct {
	Movie
}

type ListMoviesByUserRequest struct {
	Offset *int `json:"offset"`
	Limit  *int `json:"limit"`
}

func (r *ListMoviesByUserRequest) Validate() error {
	return nil
}

type ListMoviesByUserResponse struct {
	Movies       []*Movie      `json:"movies"`
	OffsetPaging *OffsetPaging `json:"paging"`
}

type ListMoviesRequest struct {
	Offset *int `json:"offset"`
	Limit  *int `json:"limit"`
}

func (r *ListMoviesRequest) Validate() error {
	return nil
}

type ListMoviesResponse struct {
	Movies       []*Movie      `json:"movies"`
	OffsetPaging *OffsetPaging `json:"paging"`
}

type Movie struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Thumbnail   string    `json:"thumbnail"`
	SharedBy    string    `json:"shared_by"`
	SharedAt    time.Time `json:"shared_at"`
}
