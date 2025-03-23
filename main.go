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
	"github.com/vishal/reservation_system/Handlers/Account"
	"github.com/vishal/reservation_system/Handlers/Booking"
	"github.com/vishal/reservation_system/Handlers/Hotels"
	Rooms "github.com/vishal/reservation_system/Handlers/Rooms"
	"github.com/vishal/reservation_system/Handlers/Users"
	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/WrongPath"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/time/rate"
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
	hotelCollection := mongoClient.Database("reservation").Collection("Hotels")
	roomCollection := mongoClient.Database("reservation").Collection("Rooms")
	accountCollecton := mongoClient.Database("reservation").Collection("Accounts")
	bookingCollection := mongoClient.Database("reservation").Collection("Bookings")

	utils.CreateEmailIndex(userCollection)
	if PORT == "" {
		PORT = ":8080"
	}
	r := http.NewServeMux()

	r.HandleFunc("/", Handlers.WrongPathTemplate)
	// user
	r.HandleFunc("/user", UserAuthenticate(Users.Users(userCollection), userCollection))
	r.HandleFunc("/user/{id}", UserAuthenticate(Users.UserDeleteOrUpdate(userCollection), userCollection))
	r.HandleFunc("/user/login", Users.Login(userCollection))

	// Hotel
	// rate limiter for hotel
	hotelApiRateLimiter := RateLimiter{
		Rate:  5,
		Burst: 2,
	}
	r.HandleFunc("/hotels", UserAuthenticate(Hotels.Hotel(hotelCollection), userCollection))
	r.HandleFunc("/allhotels", hotelApiRateLimiter.ApiRateLimiter(Hotels.Hotel(hotelCollection)))
	r.HandleFunc("/hotel/{id}", hotelApiRateLimiter.ApiRateLimiter(Hotels.HotelById(hotelCollection)))
	r.HandleFunc("/hotels/{id}", UserAuthenticate(Hotels.UpdateHotelById(hotelCollection), userCollection))

	// rooms
	r.HandleFunc("/room", UserAuthenticate(Rooms.POSTRooms(roomCollection, hotelCollection), userCollection))

	// Accounts
	r.HandleFunc("/account", UserAuthenticate(Account.AccountHandler(accountCollecton), userCollection))

	//
	r.HandleFunc("/booking", UserAuthenticate(Booking.Bookings(bookingCollection, userCollection, roomCollection, accountCollecton), userCollection))
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
			utils.ResponseWriter(w, http.StatusUnauthorized, utils.CommonError(fmt.Errorf("no token found"), http.StatusUnauthorized))
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

type RateLimiter struct {
	Rate  int
	Burst int
}

func (rl *RateLimiter) ApiRateLimiter(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second*time.Duration(rl.Rate)), rl.Burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			errMsg := struct {
				Msg  string `json:"msg"`
				Body string `json:"body"`
			}{
				Msg:  "too many Request",
				Body: "Please Try after few seconds",
			}
			utils.ResponseWriter(w, http.StatusTooManyRequests, utils.CommonError(fmt.Errorf("+%v", errMsg), http.StatusTooManyRequests))
			return
		} else {
			next(w, r)
		}
	})
}
