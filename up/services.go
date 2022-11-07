package up

import (
	"context"
)

type UserService interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

type MovieService interface {
	Create(context.Context, *CreateMovieRequest) (*CreateMovieResponse, error)
	GetMovieByUser(context.Context, *GetMovieByUserRequest) (*GetMovieByUserResponse, error)
	ListMoviesByUser(context.Context, *ListMoviesByUserRequest) (*ListMoviesByUserResponse, error)
	ListMovies(context.Context, *ListMoviesRequest) (*ListMoviesResponse, error)
}
