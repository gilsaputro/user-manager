package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gilsaputro/user-manager/internal/handler/middleware"
	user_handler "gilsaputro/user-manager/internal/handler/user"
	user_service "gilsaputro/user-manager/internal/service/user"
	user_store "gilsaputro/user-manager/internal/store/user"
	"gilsaputro/user-manager/pkg/hash"
	"gilsaputro/user-manager/pkg/postgres"
	"gilsaputro/user-manager/pkg/token"
)

func main() {
	pg, err := postgres.NewPostgresClient("host=localhost port=5492 user=user_binary dbname=management_user password=postgres_p4ssW0Rd sslmode=disable")
	if err != nil {
		fmt.Println("error pg:", err)
		return
	}
	store := user_store.NewUserStore(pg)
	tokenMethod := token.NewTokenMethod("abc", 1)
	hashMethod := hash.NewHashMethod(10)
	service := user_service.NewUserService(store, tokenMethod, hashMethod)
	handler := user_handler.NewUserHandler(service, user_handler.WithTimeoutOptions(5))

	midlewareService := middleware.NewMiddleware(tokenMethod, hashMethod)

	r := mux.NewRouter()
	r.HandleFunc("/login", handler.LoginUserHandler).Methods("POST")
	r.HandleFunc("/register", midlewareService.MiddlewareVerifyToken(handler.RegisterUserHandler)).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
