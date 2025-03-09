package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/compose-spec/compose-go/v2/dotenv"
	db "github.com/vishal/reservation_system/DB"
	"github.com/vishal/reservation_system/Handlers/Middleware"
	"github.com/vishal/reservation_system/Handlers/Users"
	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/dummy"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	err := dotenv.Load(".dev.env")
	if err != nil {
		slog.Error("No Environment File Found", slog.String("error: - ", err.Error()))
	}
	PORT := os.Getenv("port")
	MONGODB_URI := os.Getenv("MONGODB_URI")
	mongoClient := db.DB_Connection(MONGODB_URI)

	userCollection := mongoClient.Database("reservation").Collection("Users")

	utils.CreateEmailIndex(userCollection)
	if PORT == "" {
		PORT = ":8080"
	}
	r := http.NewServeMux()

	r.HandleFunc("/", Handlers.WrongPathTemplate)
	// user Collection
	r.HandleFunc("/user", UserAuthenticate(Users.Users(userCollection), userCollection))
	r.HandleFunc("/user/{id}", UserAuthenticate(Users.UserDeleteOrUpdate(userCollection), userCollection))
	r.HandleFunc("/user/login", Users.Login(userCollection))
	//checkingMiddleware
	r.HandleFunc("/middle", UserAuthenticate(Middleware.Authorize(userCollection), userCollection))

	server := http.Server{
		Addr:    PORT,
		Handler: r,
	}
	fmt.Println("********************************")
	fmt.Printf("Server is listening on %s- \n", PORT)
	fmt.Println("********************************")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			mongoClient.Disconnect(context.Background())
			log.Fatalf("Closing the Server - %v\n", err.Error())
		}
	}()
	<-c

	slog.Info("Server is Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down the Server", slog.String("error: - ", err.Error()))
	}

	slog.Info("Server ShutDown Successfully")

}

func UserAuthenticate(next http.HandlerFunc, userCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "user") {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("token")
		if token == "" {
			utils.ResponseWriter(w, http.StatusUnauthorized, utils.CommonError(fmt.Errorf("No Token Found"), http.StatusUnauthorized))
			return
		}
		claims, err := utils.ParseToken(token, os.Getenv("TOKEN_SECRET"))
		if err != nil {
			utils.ResponseWriter(w, http.StatusUnauthorized, utils.CommonError(fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized))
			return
		}
		var authorizeduser *types.User
		res := userCollection.FindOne(context.TODO(), bson.M{"email": claims.Email}).Decode(&authorizeduser)
		if res != nil {
			utils.ResponseWriter(w, http.StatusUnauthorized, utils.CommonError(fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized))
			return
		}
		ctx := context.WithValue(r.Context(), "authorizeduser", authorizeduser)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		return
	})
}
