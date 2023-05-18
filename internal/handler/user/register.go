package user

import (
	"context"
	"encoding/json"
	"fmt"
	"gilsaputro/user-manager/internal/handler/utilhttp"
	"gilsaputro/user-manager/internal/service/user"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RegisterUserRequest is list request parameter for Register Api
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// RegisterUserResponse is list response parameter for Register Api
type RegisterUserResponse struct {
	Token string `json:"token"`
}

// RegisterUserHandler is func handler for Register user
func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response utilhttp.StandardResponse
	var code int = http.StatusOK

	defer func() {
		response.Code = code
		if err == nil {
			response.Message = "success"
		} else {
			response.Message = err.Error()
		}

		data, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			log.Println("[RegisterUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body RegisterUserRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	// checking valid body
	if len(body.Username) < 1 || len(body.Password) < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	var token string

	token, _ = r.Context().Value("token").(string)

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		err = h.service.RegisterUser(user.RegisterUserServiceRequest{
			TokenRequest: token,
			Username:     body.Username,
			Password:     body.Password,
			Fullname:     body.Fullname,
			Email:        body.Email,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		code = http.StatusGatewayTimeout
		err = fmt.Errorf("Timeout")
		return
	case err = <-errChan:
		if err != nil {
			if err == user.ErrUserNameAlreadyExists {
				code = http.StatusConflict
			} else if err == user.ErrNotGuest {
				code = http.StatusUnauthorized
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseRegister(token)
}

func mapResonseRegister(r string) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	return res
}
