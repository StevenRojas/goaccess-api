package codec

import (
	"context"
	"encoding/json"
	"net/http"

	e "github.com/StevenRojas/goaccess-api/pkg/errors"
)

// DecodeLoginRequest decode login request
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	emailReq := struct {
		Email string `json:"email"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		return nil, err
	}
	if emailReq.Email == "" {
		return nil, e.HTTPBadRequestFromString("Email is missing")
	}
	return emailReq.Email, nil
}
