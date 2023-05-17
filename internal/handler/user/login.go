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

// LoginUserRequest is list request parameter for Login Api
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserResponse is list response parameter for Login Api
type LoginUserResponse struct {
	Token string `json:"token"`
}

// LoginUserHandler is func handler for login
func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[LoginUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body LoginUserRequest
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

	errChan := make(chan error, 1)
	var token string
	go func(ctx context.Context) {
		token, err = h.service.LoginUser(
			user.LoginUserServiceRequest{
				Username: body.Username,
				Password: body.Password,
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
			if err == user.ErrUserNameNotExists {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseLogin(token)
}

func mapResonseLogin(r string) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := LoginUserResponse{
		Token: r,
	}
	res.Data = data
	return res
}
