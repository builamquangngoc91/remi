package features

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"remi/pkg/golibs/idutil"
	"remi/up"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovieService_Create_Success(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("name-" + idutil.NewID()),
		Username: fmt.Sprintf("user-" + idutil.NewID()),
		Password: "password",
	}
	registerReqMarshalled, err := json.Marshal(registerReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal RegisterRequest: %w", err))
	}

	registerRespMarshalled, err := http.Post("http://localhost:8080/api/v1/register", "application/json", bytes.NewBuffer(registerReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api register: %w", err))
	}
	assert.Equal(t, registerRespMarshalled.StatusCode, http.StatusOK)

	loginReq := &up.LoginRequest{
		Username: registerReq.Username,
		Password: registerReq.Password,
	}
	loginReqMarshalled, err := json.Marshal(loginReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal loginRequest: %w", err))
	}

	loginRespMarshalled, err := http.Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewBuffer(loginReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api login: %w", err))
	}
	assert.Equal(t, loginRespMarshalled.StatusCode, http.StatusOK)

	defer loginRespMarshalled.Body.Close()

	loginRespMarshalledBody, err := ioutil.ReadAll(loginRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("can't ioutil.ReadAll loginRespMarshalled.Body %w", err))
	}

	var loginResp up.LoginResponse
	if err := json.Unmarshal(loginRespMarshalledBody, &loginResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal loginRespMarshalledBody %w", err))
	}

	assert.Equal(t, registerReq.Username, loginResp.Username)
	assert.Equal(t, registerReq.Name, loginResp.Name)
	assert.NotEmpty(t, loginResp.Token)

	createMovieReq := &up.CreateMovieRequest{
		Name:        "movie-" + idutil.NewID(),
		Description: "description-" + idutil.NewID(),
		Link:        "https://www.youtube.com/watch?v=" + idutil.NewID(),
	}

	createMovieReqMarshalled, err := json.Marshal(createMovieReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal createMovieReq: %w", err))
	}

	client := &http.Client{}
	createMovieHttpReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/createMovie", bytes.NewBuffer(createMovieReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - http.NewRequest: %w", err))
	}

	createMovieHttpReq.Header.Set("Content-Type", "application/json")
	createMovieHttpReq.Header.Set("authorization", loginResp.Token)

	createMovieRespMarshalled, err := client.Do(createMovieHttpReq)
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - client.Do: %w", err))
	}

	assert.Equal(t, http.StatusOK, createMovieRespMarshalled.StatusCode)

	defer loginRespMarshalled.Body.Close()

	createMovieRespMarshalledBody, err := ioutil.ReadAll(createMovieRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - ioutil.ReadAll %w", err))
	}

	var createMovieResp up.CreateMovieResponse
	if err := json.Unmarshal(createMovieRespMarshalledBody, &createMovieResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal createMovieRespMarshalledBody %w", err))
	}

	assert.Equal(t, http.StatusOK, createMovieRespMarshalled.StatusCode)
	assert.NotEmpty(t, createMovieResp.ID)
}

func TestMovieService_Create_Error(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("name-" + idutil.NewID()),
		Username: fmt.Sprintf("user-" + idutil.NewID()),
		Password: "password",
	}
	registerReqMarshalled, err := json.Marshal(registerReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal RegisterRequest: %w", err))
	}

	registerRespMarshalled, err := http.Post("http://localhost:8080/api/v1/register", "application/json", bytes.NewBuffer(registerReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api register: %w", err))
	}
	assert.Equal(t, registerRespMarshalled.StatusCode, http.StatusOK)

	loginReq := &up.LoginRequest{
		Username: registerReq.Username,
		Password: registerReq.Password,
	}
	loginReqMarshalled, err := json.Marshal(loginReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal loginRequest: %w", err))
	}

	loginRespMarshalled, err := http.Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewBuffer(loginReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api login: %w", err))
	}
	assert.Equal(t, loginRespMarshalled.StatusCode, http.StatusOK)

	defer loginRespMarshalled.Body.Close()

	loginRespMarshalledBody, err := ioutil.ReadAll(loginRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("can't ioutil.ReadAll loginRespMarshalled.Body %w", err))
	}

	var loginResp up.LoginResponse
	if err := json.Unmarshal(loginRespMarshalledBody, &loginResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal loginRespMarshalledBody %w", err))
	}

	assert.Equal(t, registerReq.Username, loginResp.Username)
	assert.Equal(t, registerReq.Name, loginResp.Name)
	assert.NotEmpty(t, loginResp.Token)

	createMovieReq := &up.CreateMovieRequest{
		Name:        "movie-" + idutil.NewID(),
		Description: "description-" + idutil.NewID(),
		Link:        "link-" + idutil.NewID(),
	}

	createMovieReqMarshalled, err := json.Marshal(createMovieReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal createMovieReq: %w", err))
	}

	client := &http.Client{}
	createMovieHttpReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/createMovie", bytes.NewBuffer(createMovieReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - http.NewRequest: %w", err))
	}

	createMovieHttpReq.Header.Set("Content-Type", "application/json")
	createMovieHttpReq.Header.Set("authorization", loginResp.Token)

	createMovieRespMarshalled, err := client.Do(createMovieHttpReq)
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - client.Do: %w", err))
	}
	defer loginRespMarshalled.Body.Close()

	createMovieRespMarshalledBody, err := ioutil.ReadAll(createMovieRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("createMovie - ioutil.ReadAll %w", err))
	}

	var errorResp up.ErrorResponse
	if err := json.Unmarshal(createMovieRespMarshalledBody, &errorResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal createMovieRespMarshalledBody %w", err))
	}

	assert.Equal(t, http.StatusBadRequest, createMovieRespMarshalled.StatusCode)
	assert.NotEmpty(t, errorResp.Error)
}

func TestMovieService_ListMovies_Success(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("name-" + idutil.NewID()),
		Username: fmt.Sprintf("user-" + idutil.NewID()),
		Password: "password",
	}
	registerReqMarshalled, err := json.Marshal(registerReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal RegisterRequest: %w", err))
	}

	registerRespMarshalled, err := http.Post("http://localhost:8080/api/v1/register", "application/json", bytes.NewBuffer(registerReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api register: %w", err))
	}
	assert.Equal(t, registerRespMarshalled.StatusCode, http.StatusOK)

	loginReq := &up.LoginRequest{
		Username: registerReq.Username,
		Password: registerReq.Password,
	}
	loginReqMarshalled, err := json.Marshal(loginReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal loginRequest: %w", err))
	}

	loginRespMarshalled, err := http.Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewBuffer(loginReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api login: %w", err))
	}
	assert.Equal(t, loginRespMarshalled.StatusCode, http.StatusOK)

	defer loginRespMarshalled.Body.Close()

	loginRespMarshalledBody, err := ioutil.ReadAll(loginRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("can't ioutil.ReadAll loginRespMarshalled.Body %w", err))
	}

	var loginResp up.LoginResponse
	if err := json.Unmarshal(loginRespMarshalledBody, &loginResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal loginRespMarshalledBody %w", err))
	}

	assert.Equal(t, registerReq.Username, loginResp.Username)
	assert.Equal(t, registerReq.Name, loginResp.Name)
	assert.NotEmpty(t, loginResp.Token)

	createMovieRepMap := make(map[string]*up.CreateMovieRequest)
	for i := 0; i < 5; i++ {
		createMovieReq := &up.CreateMovieRequest{
			Name:        "movie-" + idutil.NewID(),
			Description: "description-" + idutil.NewID(),
			Link:        "https://www.youtube.com/watch?v=" + idutil.NewID(),
		}

		createMovieReqMarshalled, err := json.Marshal(createMovieReq)
		if err != nil {
			log.Fatal(fmt.Errorf("can't marshal createMovieReq: %w", err))
		}

		client := &http.Client{}
		createMovieHttpReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/createMovie", bytes.NewBuffer(createMovieReqMarshalled))
		if err != nil {
			log.Fatal(fmt.Errorf("createMovie - http.NewRequest: %w", err))
		}

		createMovieHttpReq.Header.Set("Content-Type", "application/json")
		createMovieHttpReq.Header.Set("authorization", loginResp.Token)

		createMovieRespMarshalled, err := client.Do(createMovieHttpReq)
		if err != nil {
			log.Fatal(fmt.Errorf("createMovie - client.Do: %w", err))
		}

		assert.Equal(t, http.StatusOK, createMovieRespMarshalled.StatusCode)

		defer loginRespMarshalled.Body.Close()

		createMovieRespMarshalledBody, err := ioutil.ReadAll(createMovieRespMarshalled.Body)
		if err != nil {
			log.Fatal(fmt.Errorf("createMovie - ioutil.ReadAll %w", err))
		}

		var createMovieResp up.CreateMovieResponse
		if err := json.Unmarshal(createMovieRespMarshalledBody, &createMovieResp); err != nil {
			log.Fatal(fmt.Errorf("can't unmarshal createMovieRespMarshalledBody %w", err))
		}

		assert.Equal(t, http.StatusOK, createMovieRespMarshalled.StatusCode)
		assert.NotEmpty(t, createMovieResp.ID)

		createMovieRepMap[createMovieResp.ID] = createMovieReq
	}

	offset := 0
	limit := 5
	listMoviesReq := &up.ListMoviesRequest{
		Offset: &offset,
		Limit:  &limit,
	}

	listMoviesReqMarshalled, err := json.Marshal(listMoviesReq)
	if err != nil {
		log.Fatal(fmt.Errorf("listMovies - json.Marshal: %w", err))
	}

	client := &http.Client{}
	listMoviesHttpReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/listMovies", bytes.NewBuffer(listMoviesReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("listMovies - http.NewRequest: %w", err))
	}

	listMoviesHttpReq.Header.Set("Content-Type", "application/json")
	listMoviesHttpReq.Header.Set("authorization", loginResp.Token)

	listMoviesRespMarshalled, err := client.Do(listMoviesHttpReq)
	if err != nil {
		log.Fatal(fmt.Errorf("listMovies - client.Do: %w", err))
	}

	assert.Equal(t, http.StatusOK, listMoviesRespMarshalled.StatusCode)

	defer listMoviesRespMarshalled.Body.Close()

	listMoviesRespMarshalledBody, err := ioutil.ReadAll(listMoviesRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("listMovies - ioutil.ReadAll %w", err))
	}

	var listMoviesResp up.ListMoviesResponse
	if err := json.Unmarshal(listMoviesRespMarshalledBody, &listMoviesResp); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal createMovieRespMarshalledBody %w", err))
	}

	assert.Equal(t, http.StatusOK, listMoviesRespMarshalled.StatusCode)

	for _, movie := range listMoviesResp.Movies {
		assert.NotNil(t, createMovieRepMap[movie.ID])
	}
}
