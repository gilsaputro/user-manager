package user

import (
	"context"
	"encoding/json"
	"fmt"
	"gilsaputro/user-manager/internal/handler/utilhttp"
	"gilsaputro/user-manager/internal/service/user"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GetByIDUserResponse is list response parameter for GetByID Api
type GetByIDUserResponse struct {
	User UserInfo `json:"user"`
}

// GetByIDUserHandler is func handler for GetByID user
func (h *UserHandler) GetByIDUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[GetByIDUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	// checking valid body
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
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
	var userInfo user.UserServiceInfo
	go func(ctx context.Context) {
		userInfo, err = h.service.GetUserByID(user.GetByIDServiceRequest{
			TokenRequest: token,
			UserId:       int64(id),
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
			if err == user.ErrUnauthorized || err == user.ErrCannotGetOtherUser {
				code = http.StatusUnauthorized
			} else if err == user.ErrDataNotFound {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResponseGetByID(userInfo)
}

func mapResponseGetByID(result user.UserServiceInfo) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse

	res.Data = GetByIDUserResponse{
		User: UserInfo{
			UserID:      result.UserId,
			Username:    result.Username,
			Email:       result.Email,
			Fullname:    result.Fullname,
			CreatedDate: result.CreatedDate,
		},
	}

	return res
}
