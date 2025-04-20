package main

import (
	"GoPattern/config"
	"GoPattern/db"
	"GoPattern/handlers"
	"GoPattern/internal/logger"
	"GoPattern/internal/redisdb"
	base "GoPattern/internal/shared"
	"GoPattern/middleware"
	"GoPattern/repository"
	authRouter "GoPattern/router/auth"
	MDRouter "GoPattern/router/masterdata"
	userRouter "GoPattern/router/user"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	cfg := config.LoadConfig()
	db.InitDB(cfg)
	redisClient := redisdb.NewRedisClient("localhost:6379", "", 0)

	userRepo := repository.NewUserRepository(db.GetDB())
	userHandler := handlers.NewUserHandler(userRepo, redisClient)
	baseHandler := base.NewBaseHandler(redisClient)

	redisHandler := handlers.NewAuthHandler(baseHandler)
	dbRepoMd := repository.NewProvinceRepository(db.GetDB())
	mdHandler := handlers.NewMasterdataHandler(baseHandler, dbRepoMd)
	// Tạo module
	userModule := userRouter.NewModule(userHandler)
	authModule := authRouter.NewModule(redisHandler)
	mdModule := MDRouter.NewModule(mdHandler)

	// Gắn module vào router
	r := chi.NewRouter()
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.RequestLogger)

	r.Mount("/users", userModule.Router())    // users/
	r.Mount("/auth", authModule.Router())     // auth/login, auth/refresh
	r.Mount("/masterdata", mdModule.Router()) // masterdata/countries

	fmt.Println("✅ Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
