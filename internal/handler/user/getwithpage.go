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

// GetWithPageUserRequest is list request parameter for GetWithPage Api
type GetWithPageUserRequest struct {
	Size   int `json:"size"`
	Cursor int `json:"cursor"`
}

// GetWithPageUserResponse is list response parameter for GetWithPage Api
type GetWithPageUserResponse struct {
	User   []UserInfo `json:"users"`
	Cursor *int       `json:"next_cursor,omitempty"`
}

// GetWithPageUserHandler is func handler for GetWithPage user
func (h *UserHandler) GetWithPageUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[GetWithPageUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body GetWithPageUserRequest
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

	var token string
	var ok bool
	token, ok = r.Context().Value("token").(string)
	if !ok {
		code = http.StatusInternalServerError
		err = fmt.Errorf("Internal Server Error")
		return
	}

	errChan := make(chan error, 1)
	var listUser user.GetAllUserWithPaggingServiceResponse
	go func(ctx context.Context) {
		listUser, err = h.service.GetAllUserWithPagging(user.GetAllUserWithPaggingServiceRequest{
			TokenRequest: token,
			Size:         body.Size,
			Cursor:       body.Cursor,
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
			if err == user.ErrUnauthorized {
				code = http.StatusUnauthorized
			} else if err == user.ErrDataNotFound {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResponseGetWithPage(listUser)
}

func mapResponseGetWithPage(result user.GetAllUserWithPaggingServiceResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	var data GetWithPageUserResponse
	var users []UserInfo

	for _, user := range result.UserList {
		users = append(users, UserInfo{
			UserID:      user.UserId,
			Username:    user.Username,
			Email:       user.Email,
			Fullname:    user.Fullname,
			CreatedDate: user.CreatedDate,
		})
	}

	data.User = users
	if result.NextCursor > 0 {
		data.Cursor = &result.NextCursor
	}

	res.Data = data

	return res
}
