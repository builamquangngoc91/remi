package up

import (
	"strings"

	"remi/pkg/xerror"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (r *RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "username can't be null")
	}
	if strings.TrimSpace(r.Password) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "password can't be null")
	}
	if strings.TrimSpace(r.Name) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "name can't be null")
	}

	return nil
}

type RegisterResponse struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "username can't be null")
	}
	if strings.TrimSpace(r.Password) == "" {
		return xerror.ErrorM(xerror.InvalidArgument, nil, "password can't be null")
	}
	return nil
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}
