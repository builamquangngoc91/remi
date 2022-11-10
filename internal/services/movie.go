package services

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"remi/internal/entities"
	"remi/internal/repositories"
	"remi/pkg/golibs/idutil"
	"remi/pkg/xerror"
	"remi/up"
)

var _ up.MovieService = &MovieService{}

type MovieService struct {
	movieRepo *repositories.MovieRepository
	userRepo  *repositories.UserRepository
	url       string
}

func NewMovieService(db *sql.DB, url string) *MovieService {
	return &MovieService{
		userRepo:  repositories.NewUserRepository(db),
		movieRepo: repositories.NewMovieRepository(db),
		url:       url,
	}
}

func (s *MovieService) Create(ctx context.Context, req *up.CreateMovieRequest) (*up.CreateMovieResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, xerror.Error(xerror.InvalidArgument, err)
	}

	if !strings.Contains(req.Link, "youtube.com") {
		return nil, xerror.Error(xerror.InvalidArgument, fmt.Errorf("we only support link youtube"))
	}

	_url, err := url.Parse(req.Link)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("url.Parse: %w", err))
	}
	params, err := url.ParseQuery(_url.RawQuery)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("url.ParseQuery: %w", err))
	}

	youtubeVideoID := params["v"][0]

	userID, _ := userIDFromCtx(ctx)

	now := time.Now()
	movieEnt := &entities.Movie{
		ID:          idutil.NewID(),
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Thumbnail:   fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", youtubeVideoID),
		SharedBy:    userID,
		SharedAt:    &now,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if err := s.movieRepo.Create(ctx, movieEnt); err != nil {
		return nil, xerror.Error(xerror.Internal, err)
	}

	return &up.CreateMovieResponse{
		ID: movieEnt.ID,
	}, nil
}

func (s *MovieService) GetMovieByUser(ctx context.Context, req *up.GetMovieByUserRequest) (*up.GetMovieByUserResponse, error) {
	userID, _ := userIDFromCtx(ctx)
	movie, err := s.movieRepo.FindByIDAndUserID(ctx, req.ID, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.movieRepo.Get: %w", err))
		}
		return nil, xerror.Error(xerror.InvalidArgument, fmt.Errorf("movie (%s) not found", req.ID))
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.userRepo.FindByID: %w", err))
	}

	return &up.GetMovieByUserResponse{
		Movie: up.Movie{
			ID:          movie.ID,
			Name:        movie.Name,
			Link:        movie.Link,
			Thumbnail:   movie.Thumbnail,
			Description: movie.Description,
			SharedBy:    user.Name,
			SharedAt:    *movie.SharedAt,
		},
	}, nil
}

func (s *MovieService) ListMoviesByUser(ctx context.Context, req *up.ListMoviesByUserRequest) (resp *up.ListMoviesByUserResponse, _ error) {
	userID, _ := userIDFromCtx(ctx)
	movies, err := s.movieRepo.List(
		ctx,
		&repositories.ListMoviesArgs{
			UserID: &userID,
		},
	)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.movieRepo.List: %w", err))
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.userRepo.FindByID: %w", err))
	}

	resp = &up.ListMoviesByUserResponse{
		OffsetPaging: &up.OffsetPaging{
			Limit:  *req.Limit,
			Offset: *req.Offset,
		},
	}
	for _, movie := range movies {
		resp.Movies = append(resp.Movies, &up.Movie{
			ID:          movie.ID,
			Name:        movie.Name,
			Description: movie.Description,
			Link:        movie.Link,
			Thumbnail:   movie.Thumbnail,
			SharedBy:    user.Name,
			SharedAt:    *movie.SharedAt,
		})
	}

	return resp, nil
}

func (s *MovieService) ListMovies(ctx context.Context, req *up.ListMoviesRequest) (resp *up.ListMoviesResponse, _ error) {
	movies, err := s.movieRepo.List(
		ctx,
		&repositories.ListMoviesArgs{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.movieRepo.List: %w", err))
	}

	userIDs := make([]string, 0, len(movies))
	for _, movie := range movies {
		userIDs = append(userIDs, movie.SharedBy)
	}

	users, err := s.userRepo.List(ctx, &repositories.ListUsersArgs{IDs: userIDs})
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.userRepo.List: %w", err))
	}
	userMap := make(map[string]*entities.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	resp = &up.ListMoviesResponse{
		OffsetPaging: &up.OffsetPaging{
			Offset: *req.Offset,
			Limit:  *req.Limit,
		},
	}
	for _, movie := range movies {
		user := userMap[movie.SharedBy]
		resp.Movies = append(resp.Movies, &up.Movie{
			ID:          movie.ID,
			Name:        movie.Name,
			Description: movie.Description,
			Link:        movie.Link,
			Thumbnail:   movie.Thumbnail,
			SharedBy:    user.Name,
			SharedAt:    *movie.SharedAt,
		})
	}

	return resp, nil
}

func (s *MovieService) LikeMovie(ctx context.Context, req *up.LikeMovieRequest) (*up.LikeMovieResponse, error) {
	
	return nil, nil
}

func (s *MovieService) DislikeMovie(ctx context.Context, req *up.DislikeMovieRequest) (*up.DislikeMovieResponse, error) {
	return nil, nil
}

func (s *MovieService) GetCreateMoviePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/movie_create.html"))

	data := Data{
		URL: s.url,
	}
	tmpl.Execute(w, data)
}

type ViewMovieData struct {
	Name        string
	Link        string
	SharedBy    string
	Description string
}

func (s *MovieService) GetViewMoviePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/movie.html"))

	params := r.URL.Query()
	ids, ok := params["id"]
	if !ok || len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	ctx := context.Background()

	movie, err := s.movieRepo.FindByID(ctx, ids[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	user, err := s.userRepo.FindByID(ctx, movie.SharedBy)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	_url, err := url.Parse(movie.Link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	params, err = url.ParseQuery(_url.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	youtubeVideoID := params["v"][0]

	viewMovieData := ViewMovieData{
		Link:        fmt.Sprintf("https://www.youtube.com/embed/%s", youtubeVideoID),
		Name:        movie.Name,
		Description: movie.Description,
		SharedBy:    user.Name,
	}

	tmpl.Execute(w, viewMovieData)
}
