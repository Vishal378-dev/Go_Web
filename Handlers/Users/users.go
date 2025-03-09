package Users

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/dummy"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Users(userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var users []types.User
			cursor, err := userCollection.Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.D{{
				Key:   "password",
				Value: 0,
			}}))
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
				return
			}
			err = cursor.All(context.TODO(), &users)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing entries"), http.StatusBadRequest))
				return
			}

			Response := types.SuccessResponse{
				Status: http.StatusOK,
				Data:   users,
			}
			utils.ResponseWriter(w, http.StatusOK, Response)
			return
		} else if r.Method == "POST" {
			var user types.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			err = user.ValidateRequest()
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			var getUser types.User
			err = userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&getUser)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
					return
				}
			}
			if getUser.Email != "" {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("email already present"), http.StatusBadRequest))
				return
			}
			hashedPassword, _ := utils.HashPassword(user.Password)
			user.Password = hashedPassword
			user.Created = time.Now()
			user.Updated = time.Now()
			result, err := userCollection.InsertOne(context.TODO(), user)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			token, err := utils.NewAccessToken(types.UserClaims{Name: user.Name, Email: user.Email, Phone: user.Phone, Role: user.Role, RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now().UTC()), ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * 50))}})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while generating the token"), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusCreated, map[string]any{"user": result, "token": token})
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}

func UserDeleteOrUpdate(userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qid := strings.TrimPrefix(r.URL.Path, "/user/")
		if qid == "" {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("invalid id"), http.StatusBadRequest))
			return
		}
		id, err := bson.ObjectIDFromHex(qid)
		if err != nil {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
			return
		}
		currentUser, ok := r.Context().Value("authorizeduser").(*types.User)
		if !ok {
			utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing the data"), http.StatusBadRequest))
			return
		}
		if r.Method == "PUT" {
			var user types.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			updateFields := bson.M{}
			var userById types.User
			err = userCollection.FindOne(context.TODO(), bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{
				{
					Key:   "password",
					Value: 0,
				},
			})).Decode(&userById)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
					return
				}
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error fetching user: %v", err), http.StatusInternalServerError))
				return
			}
			if userById.Email != currentUser.Email {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("you can only update your entry"), http.StatusBadRequest))
				return
			}
			if user.Email != "" {
				updateFields["email"] = user.Email
			}
			if user.Name != "" {
				updateFields["name"] = user.Name
			}
			if user.Phone != "" {
				updateFields["phone"] = user.Phone
			}
			if len(updateFields) == 0 {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("no valid input to be updated"), http.StatusBadRequest))
				return
			}

			updateFields["updated_at"] = time.Now()
			result, err := userCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, bson.M{"$set": updateFields})
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("err while updating"), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{"msg": result})
			return
		} else if r.Method == "DELETE" {
			filter := bson.D{{Key: "_id", Value: id}}
			// result := userCollection.FindOne(context.TODO(), filter)
			// if result.Err() == mongo.ErrNoDocuments {
			// 	utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf(result.Err().Error()), http.StatusBadRequest))
			// 	return
			// }
			var userById types.User
			err = userCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.D{
				{
					Key:   "password",
					Value: 0,
				},
			})).Decode(&userById)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
					return
				}
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error fetching user: %v", err), http.StatusInternalServerError))
				return
			}
			if userById.Email != currentUser.Email {
				if currentUser.Role != "admin" {
					utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("you can only update your entry"), http.StatusBadRequest))
					return
				}
			}

			_, err := userCollection.DeleteOne(context.TODO(), filter)
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("err while deleting"), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{"msg": fmt.Sprintf("Successfully deleted - %v", id), "status": http.StatusOK})
			return
		} else if r.Method == "GET" {
			var userById types.User
			err := userCollection.FindOne(context.TODO(), bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{
				{
					Key:   "password",
					Value: 0,
				},
			})).Decode(&userById)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
					return
				}
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error fetching user: %v", err), http.StatusInternalServerError))
				return
			}
			if currentUser.Email == userById.Email {
				utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{"msg": userById})
				return
			} else {
				utils.ResponseWriter(w, http.StatusUnauthorized, utils.CommonError(fmt.Errorf("you are not allowed to access the record"), http.StatusUnauthorized))
				return
			}
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}

func Login(userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var requestBody types.UserRequestSignUp
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			err = requestBody.ValidateUserRequestSignup()
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			var userByEmail types.User
			err = userCollection.FindOne(context.TODO(), bson.M{"email": requestBody.Email}).Decode(&userByEmail)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					slog.Error("Signup : -User Entry Not found")
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("incorrect credentials email or password"), http.StatusBadRequest))
					return
				}
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("incorrect credentials email or password"), http.StatusInternalServerError))
				return
			}
			isCorrectPassword := utils.ComparePassword(requestBody.Password, userByEmail.Password)
			if !isCorrectPassword {
				slog.Error("Signup : -Password Incorrect")
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("incorrect credentials email or password"), http.StatusInternalServerError))
				return
			}
			token, err := utils.NewAccessToken(types.UserClaims{Name: userByEmail.Name, Email: userByEmail.Email, Phone: userByEmail.Phone, Role: userByEmail.Role, RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now().UTC()), ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * 5))}})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while generating the token"), http.StatusBadRequest))
				return
			}
			userByEmail.Password = ""
			utils.ResponseWriter(w, http.StatusCreated, map[string]any{"user": userByEmail, "token": token})
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}
