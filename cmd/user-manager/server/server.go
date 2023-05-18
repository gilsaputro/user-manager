package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"gilsaputro/user-manager/cmd/user-manager/config"
	"gilsaputro/user-manager/internal/handler/middleware"
	user_handler "gilsaputro/user-manager/internal/handler/user"
	user_service "gilsaputro/user-manager/internal/service/user"
	user_store "gilsaputro/user-manager/internal/store/user"
	"gilsaputro/user-manager/pkg/hash"
	"gilsaputro/user-manager/pkg/postgres"
	"gilsaputro/user-manager/pkg/token"
	"gilsaputro/user-manager/pkg/vault"
)

// Servcer is list configuration to run Server
type Server struct {
	cfg         config.Config
	vault       vault.VaultMethod
	hashMethod  hash.HashMethod
	tokenMethod token.TokenMethod
	postgres    postgres.PostgresMethod
	middleware  middleware.Middleware
	userStore   user_store.UserStoreMethod
	userService user_service.UserServiceMethod
	userHandler user_handler.UserHandler
	httpServer  *http.Server
}

// NewServer is func to create server with all configuration
func NewServer() (*Server, error) {
	s := &Server{}

	// ======== Init Dependencies Related ========
	// Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return s, err
	}

	// Init Vault
	{
		token := os.Getenv("VAULT_TOKEN")
		if len(token) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_TOKEN")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_TOKEN")
		}

		host := os.Getenv("VAULT_HOST")
		if len(host) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_HOST")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_HOST")
		}

		vaultMethod, err := vault.NewVaultClient(token, host)
		if err != nil {
			fmt.Print("[Got Error]-Vault :", err)
		}
		s.vault = vaultMethod

		log.Println("Init-Vault")
	}

	// Get Config from yaml and replace by secret
	{
		secret, err := s.vault.GetConfig()
		if err != nil {
			fmt.Print("[Got Error]-Load Secret :", err)
		}
		cfg, err := config.GetConfig(secret)
		if err != nil {
			fmt.Print("[Got Error]-Load Config :", err)
			return s, err
		}
		s.cfg = cfg

		log.Println("LOAD-Config")
	}

	// Init Postgres
	{
		postgresMethod, err := postgres.NewPostgresClient(s.cfg.Postgres.Config)
		if err != nil {
			fmt.Print("[Got Error]-Postgres :", err)
			return s, err
		}

		s.postgres = postgresMethod

		log.Println("Init-Postgres")
	}

	// Init Hash Package
	{
		hashMethod := hash.NewHashMethod(s.cfg.Hash.Cost)
		s.hashMethod = hashMethod
		log.Println("Init-Hash Package")
	}

	// Init Token Package
	{
		tokenMethod := token.NewTokenMethod(s.cfg.Token.Secret, s.cfg.Token.ExpInHour)
		s.tokenMethod = tokenMethod
		log.Println("Init-Token Package")
	}

	// ======== Init Dependencies Store ========
	// Init User Store
	{
		userStore := user_store.NewUserStore(s.postgres)
		s.userStore = userStore
		log.Println("Init-User Store")
	}

	// ======== Init Dependencies Service ========
	// Init User Service
	{
		userService := user_service.NewUserService(s.userStore, s.tokenMethod, s.hashMethod)
		s.userService = userService
		log.Println("Init-User Service")
	}

	// ======== Init Dependencies Handler ========
	// Init Middleware
	{
		midlewareService := middleware.NewMiddleware(s.tokenMethod)
		s.middleware = midlewareService
		log.Println("Init-Middleware")
	}

	// Init User Handler
	{
		var opts []user_handler.Option
		opts = append(opts, user_handler.WithTimeoutOptions(s.cfg.UserHandler.TimeoutInSec))
		userHandler := user_handler.NewUserHandler(s.userService, opts...)
		s.userHandler = *userHandler
		log.Println("Init-User Handler")
	}

	// Init Router
	{
		r := mux.NewRouter()
		// Init Guest Path
		r.HandleFunc("/login", s.userHandler.LoginUserHandler).Methods("POST")
		r.HandleFunc("/register", s.middleware.MiddlewareParseToken(s.userHandler.RegisterUserHandler)).Methods("POST")

		port := ":" + s.cfg.Port
		log.Println("running on port ", port)

		server := &http.Server{
			Addr:    port,
			Handler: r,
		}

		s.httpServer = server
	}
	return s, nil
}

func (s *Server) Start() int {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for a signal to shut down the application
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout to allow the server to cleanly shut down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	log.Println("complete, shutting down.")
	return 0
}

// Run is func to create server and invoke Start()
func Run() int {
	s, err := NewServer()
	if err != nil {
		return 1
	}

	return s.Start()
}
