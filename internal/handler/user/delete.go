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

// DeleteUserRequest is list request parameter for Delete Api
type DeleteUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DeleteUserHandler is func handler for Delete user
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[DeleteUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body DeleteUserRequest
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
	var ok bool
	token, ok = r.Context().Value("token").(string)
	if !ok {
		code = http.StatusInternalServerError
		err = fmt.Errorf("Internal Server Error")
		return
	}

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		err = h.service.DeleteUser(user.DeleteUserServiceRequest{
			TokenRequest: token,
			Username:     body.Username,
			Password:     body.Password,
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
			if err == user.ErrUserNameNotExists || err == user.ErrPasswordIsIncorrect {
				code = http.StatusBadRequest
			} else if err == user.ErrUnauthorized || err == user.ErrCannotDeleteOtherUser {
				code = http.StatusUnauthorized
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResponseDelete()
}

func mapResponseDelete() utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	return res
}
