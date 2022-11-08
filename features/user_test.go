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

func TestUserService_Register_Success(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("user-" + idutil.NewID()),
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
}

func TestUserService_Register_Error(t *testing.T) {
	registerReq := &up.RegisterRequest{}
	registerReqMarshalled, err := json.Marshal(registerReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal RegisterRequest: %w", err))
	}

	registerRespMarshalled, err := http.Post("http://localhost:8080/api/v1/register", "application/json", bytes.NewBuffer(registerReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api register: %w", err))
	}

	registerRespMarshalledBody, err := ioutil.ReadAll(registerRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("can't ioutil.ReadAll registerRespMarshalled.Body %w", err))
	}

	var errorResponse up.ErrorResponse
	if err := json.Unmarshal(registerRespMarshalledBody, &errorResponse); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal registerRespMarshalledBody %w", err))
	}
	assert.Equal(t, registerRespMarshalled.StatusCode, http.StatusBadRequest)
	assert.NotEmpty(t, errorResponse.Error)
}

func TestUserService_Login_Success(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("user-" + idutil.NewID()),
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
		log.Fatal(fmt.Errorf("can't marshal LoginRequest: %w", err))
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
}

func TestUserService_Login_Error(t *testing.T) {
	registerReq := &up.RegisterRequest{
		Name:     fmt.Sprintf("user-" + idutil.NewID()),
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
		Username: registerReq.Username + "-failed",
		Password: registerReq.Password + "-failed",
	}
	loginReqMarshalled, err := json.Marshal(loginReq)
	if err != nil {
		log.Fatal(fmt.Errorf("can't marshal LoginRequest: %w", err))
	}

	loginRespMarshalled, err := http.Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewBuffer(loginReqMarshalled))
	if err != nil {
		log.Fatal(fmt.Errorf("error when call api login: %w", err))
	}
	defer loginRespMarshalled.Body.Close()

	loginRespMarshalledBody, err := ioutil.ReadAll(loginRespMarshalled.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("can't ioutil.ReadAll loginRespMarshalled.Body %w", err))
	}

	var errorResponse up.ErrorResponse
	if err := json.Unmarshal(loginRespMarshalledBody, &errorResponse); err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal loginRespMarshalledBody %w", err))
	}

	assert.Equal(t, http.StatusUnauthorized, loginRespMarshalled.StatusCode)
	assert.NotEmpty(t, errorResponse.Error)
}
