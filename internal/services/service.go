package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"remi/pkg/xerror"

	"github.com/dgrijalva/jwt-go"
)

type (
	AuthType     int
	ResponseType int
)

const (
	None = AuthType(0)
	User = AuthType(1)

	JSON = ResponseType(0)
	HTML = ResponseType(1)
	File = ResponseType(2)
)

type Decl struct {
	Auth         AuthType
	HandlerFunc  interface{}
	ResponseType ResponseType
}

type RemiService struct {
	jwtKey       string
	userService  *UserService
	movieService *MovieService
	acl          map[string]map[string]Decl
}

func NewRemiService(db *sql.DB, JWTKey string) *RemiService {
	userService := NewUserService(db, JWTKey)
	movieService := NewMovieService(db)

	return &RemiService{
		jwtKey:       JWTKey,
		userService:  userService,
		movieService: movieService,
		acl: map[string]map[string]Decl{
			"/api/v1/register": {
				http.MethodPost: Decl{
					HandlerFunc:  userService.Register,
					Auth:         None,
					ResponseType: JSON,
				},
			},
			"/api/v1/login": {
				http.MethodPost: Decl{
					HandlerFunc:  userService.Login,
					Auth:         None,
					ResponseType: JSON,
				},
			},
			"/api/v1/createMovie": {
				http.MethodPost: Decl{
					HandlerFunc:  movieService.Create,
					Auth:         User,
					ResponseType: JSON,
				},
			},
			"/api/v1/getMovieByUser": {
				http.MethodPost: Decl{
					HandlerFunc:  movieService.GetMovieByUser,
					Auth:         User,
					ResponseType: JSON,
				},
			},
			"/api/v1/listMoviesByUser": {
				http.MethodPost: Decl{
					HandlerFunc:  movieService.ListMoviesByUser,
					Auth:         User,
					ResponseType: JSON,
				},
			},
			"/api/v1/listMovies": {
				http.MethodPost: Decl{
					HandlerFunc:  movieService.ListMovies,
					Auth:         None,
					ResponseType: JSON,
				},
			},
			"/login": {
				http.MethodGet: Decl{
					HandlerFunc:  userService.GetLoginPage,
					Auth:         None,
					ResponseType: HTML,
				},
			},
			"/register": {
				http.MethodGet: Decl{
					HandlerFunc:  userService.GetRegisterPage,
					Auth:         None,
					ResponseType: HTML,
				},
			},
			"/": {
				http.MethodGet: Decl{
					HandlerFunc:  userService.GetHomePage,
					Auth:         None,
					ResponseType: HTML,
				},
			},
			"/movies": {
				http.MethodGet: Decl{
					HandlerFunc:  movieService.GetCreateMoviePage,
					Auth:         None,
					ResponseType: HTML,
				},
			},
		},
	}
}

func (s *RemiService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "*")

	if req.Method == http.MethodOptions {
		resp.WriteHeader(http.StatusOK)
		return
	}

	handler, ok := s.acl[req.URL.Path]
	if !ok {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	decl, ok := handler[req.Method]
	if !ok {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	handlerFunc := decl.HandlerFunc
	auth := decl.Auth

	// authorization
	switch auth {
	case User:
		var ok bool
		req, ok = s.validToken(req)
		if !ok {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}
	case None:
		// no-op
	}

	switch decl.ResponseType {
	case JSON:
		// Call functions of services
		// Todo: check function (type of arguments, outputs)
		typ := reflect.TypeOf(handlerFunc)
		funcArgs := reflect.New(typ.In(1).Elem()).Interface()

		// parse body and request values
		switch req.Method {
		case http.MethodPost:
			err := json.NewDecoder(req.Body).Decode(funcArgs)
			defer req.Body.Close()
			if err != nil {
				log.Panicf("error %v", err)
				return
			}

		case http.MethodGet:
			funcArgsTyp := reflect.TypeOf(funcArgs).Elem()
			for i := 0; i < funcArgsTyp.NumField(); i++ {
				jsonTag := funcArgsTyp.Field(i).Tag.Get("json")
				reflect.ValueOf(funcArgs).Elem().Field(i).Set(reflect.ValueOf(req.FormValue(jsonTag)))
			}
		}

		ctx := req.Context()
		outs := reflect.ValueOf(handlerFunc).Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(funcArgs)})
		resultFunc := outs[0].Interface()
		errFunc := outs[1].Interface()

		resp.Header().Set("Content-Type", "application/json")
		if errFunc != nil {
			err := errFunc.(xerror.XError)
			resp.WriteHeader(err.HttpStatus())
			json.NewEncoder(resp).Encode(map[string]string{
				"error": err.Message,
			})
			return
		} else {
			json.NewEncoder(resp).Encode(map[string]interface{}{
				"data": resultFunc,
			})
		}

	case HTML:
		reflect.ValueOf(handlerFunc).Call([]reflect.Value{reflect.ValueOf(resp), reflect.ValueOf(req)})
	}
}

func (s *RemiService) validToken(req *http.Request) (*http.Request, bool) {
	token := req.Header.Get("Authorization")

	claims := make(jwt.MapClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		log.Println(err)
		return req, false
	}

	if !t.Valid {
		return req, false
	}

	id, ok := claims["id"].(string)
	if !ok {
		return req, false
	}

	req = req.WithContext(context.WithValue(req.Context(), userAuthKey(0), id))
	return req, true
}

type userAuthKey int8

func userIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(userAuthKey(0))
	id, ok := v.(string)
	return id, ok
}
