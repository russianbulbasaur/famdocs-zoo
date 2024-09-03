package main

import (
	"database/sql"
	"encoding/json"
	"famdocs-zoo/internal/handlers"
	"famdocs-zoo/internal/repositories"
	"famdocs-zoo/internal/services"
	"famdocs-zoo/pkg"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

var authHandler handlers.AuthHandler
var userHandler handlers.UserHandler
var familyHandler handlers.FamilyHandler
var folderHandler handlers.FolderHandler
var fileHandler handlers.FileHandler
var db *sql.DB

func main() {
	loadEnv()
	db = pkg.InitDBConnection()

	authRepo, userRepo, famRepo, folderRepo, fileRepo := initRepositories(db)
	authServ, userServ, famServ, folderServ, fileServ :=
		initServices(authRepo, userRepo, famRepo, folderRepo, fileRepo)
	initHandlers(authServ, userServ, famServ, folderServ, fileServ)
	startServer()
}

func initRepositories(db *sql.DB) (repositories.AuthRepository, repositories.UserRepository,
	repositories.FamilyRepository, repositories.FolderRepository, repositories.FileRepository) {
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	famRepo := repositories.NewFamilyRepository(db)
	folderRepo := repositories.NewFolderRepository(db)
	fileRepo := repositories.NewFileRepository(db)
	return authRepo, userRepo, famRepo, folderRepo, fileRepo
}

func initServices(authRepo repositories.AuthRepository, userRepo repositories.UserRepository,
	famRepo repositories.FamilyRepository,
	folderRepo repositories.FolderRepository, fileRepo repositories.FileRepository) (services.AuthService, services.UserService,
	services.FamilyService, services.FolderService, services.FileService) {
	authServ := services.NewAuthService(authRepo)
	userServ := services.NewUserService(userRepo)
	famServ := services.NewFamilyService(famRepo)
	folderServ := services.NewFolderService(folderRepo)
	fileServ := services.NewFileService(fileRepo)
	return authServ, userServ, famServ, folderServ, fileServ
}

func initHandlers(authServ services.AuthService, userServ services.UserService,
	famServ services.FamilyService,
	folderServ services.FolderService,
	fileServ services.FileService) {
	authHandler = handlers.NewAuthHandler(authServ)
	userHandler = handlers.NewUserHandler(userServ)
	familyHandler = handlers.NewFamilyHandler(famServ)
	folderHandler = handlers.NewFolderHandler(folderServ)
	fileHandler = handlers.NewFileHandler(fileServ)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env")
	}
}

func startServer() {
	port := os.Getenv("PORT")
	appRouter := chi.NewRouter()
	appRouter.Get("/auth/login", authHandler.Login)
	appRouter.Post("/auth/signup", authHandler.Signup)
	appRouter.Group(initProtectedRoutes)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), appRouter)
	if err != nil {
		log.Fatal(err)
	}
}

func initProtectedRoutes(appRouter chi.Router) {
	appRouter.Use(jwtAuthenticator)
	appRouter.Mount("/user", userRouter())
	appRouter.Mount("/family", familyRouter())
}

func userRouter() chi.Router {
	userRouter := chi.NewRouter()
	userRouter.Get("/get", userHandler.Get)
	userRouter.Get("/get-families", userHandler.GetUserFamilies)
	userRouter.Post("/update", userHandler.Update)
	return userRouter
}

func familyRouter() chi.Router {
	familyRouter := chi.NewRouter()
	familyRouter.Mount("/folder", folderRouter())
	familyRouter.Post("/join", familyHandler.JoinFamily)
	familyRouter.Post("/create", familyHandler.Create)
	familyRouter.Get("/get-root-folder", familyHandler.GetRootFolder)
	return familyRouter
}

func folderRouter() chi.Router {
	folderRouter := chi.NewRouter()
	folderRouter.Mount("/file", fileRouter())
	folderRouter.Post("/create", folderHandler.Create)
	folderRouter.Get("/get-contents", folderHandler.GetContents)
	return folderRouter
}

func fileRouter() chi.Router {
	fileRouter := chi.NewRouter()
	fileRouter.Post("/create", fileHandler.Create)
	fileRouter.Post("/delete", fileHandler.Delete)
	fileRouter.Post("/download", fileHandler.Download)
	return fileRouter
}

func jwtAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		data, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			key := []byte(os.Getenv("JWT_SECRET"))
			return key, nil
		})
		if err != nil || !data.Valid {
			http.Error(w, "token invalid", http.StatusUnauthorized)
			return
		}
		sub, err := data.Claims.GetSubject()
		if err != nil {
			log.Println(err)
			return
		}
		userRepo := repositories.NewUserRepository(db)
		userId, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		user, err := userRepo.GetUserFromId(userId)
		if err != nil {
			log.Println(err)
			return
		}
		userString, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return
		}
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}
		r.Form.Set("user", string(userString))
		next.ServeHTTP(w, r)
	})
}
