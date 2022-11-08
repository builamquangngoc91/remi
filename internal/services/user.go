package services

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"remi/internal/entities"
	"remi/internal/repositories"
	"remi/pkg/crypto"
	"remi/pkg/golibs/idutil"
	"remi/pkg/xerror"
	"remi/up"

	"github.com/dgrijalva/jwt-go"
)

var _ up.UserService = &UserService{}

type UserService struct {
	userRepo *repositories.UserRepository
	jwtKey   string
	url      string
}

func NewUserService(db *sql.DB, jwtKey, url string) *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(db),
		jwtKey:   jwtKey,
		url:      url,
	}
}

func (s *UserService) Register(ctx context.Context, req *up.RegisterRequest) (*up.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, xerror.Error(xerror.InvalidArgument, err)
	}

	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, xerror.Error(xerror.Internal, err)
	}

	if user != nil {
		return nil, xerror.Error(xerror.InvalidArgument, fmt.Errorf("user exists with the given username"))
	}

	password, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, fmt.Errorf("crypto.HashPassword: %w", err))
	}

	now := time.Now()
	err = s.userRepo.Create(ctx, &entities.User{
		ID:        idutil.NewID(),
		Username:  req.Username,
		Password:  password,
		Name:      req.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
	})
	if err != nil {
		return nil, xerror.Error(xerror.Internal, err)
	}

	return &up.RegisterResponse{}, nil
}

func (s *UserService) Login(ctx context.Context, req *up.LoginRequest) (*up.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, xerror.Error(xerror.Internal, fmt.Errorf("s.userRepo.FindByUsername: %w", err))
		}
		return nil, xerror.Error(xerror.UnAuthorized, fmt.Errorf("incorrect username/pwd: %w", err))
	}

	if user == nil || !crypto.CheckPasswordHash(req.Password, user.Password) {
		return nil, xerror.Error(xerror.UnAuthorized, fmt.Errorf("incorrect username/pwd"))
	}

	token, err := s.createToken(user.ID, user.Username)
	if err != nil {
		return nil, xerror.Error(xerror.Internal, err)
	}

	return &up.LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Token:    token,
	}, nil
}

func (s *UserService) createToken(id, username string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = id
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(s.jwtKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

type Data struct {
	URL string
}

func (s *UserService) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/sign_in.html"))

	data := Data{
		URL: s.url,
	}
	tmpl.Execute(w, data)
}

func (s *UserService) GetRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/sign_up.html"))

	data := Data{
		URL: s.url,
	}
	tmpl.Execute(w, data)
}

func (s *UserService) GetHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))

	data := Data{
		URL: s.url,
	}
	tmpl.Execute(w, data)
}
