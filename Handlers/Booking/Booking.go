package Booking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	constants "github.com/vishal/reservation_system/Constants"
	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/WrongPath"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bookings(bookingCollection, userCollection, roomCollection, accountCollecton *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := utils.Ctx(5)
		defer cancel()
		if r.Method == "POST" {
			var bookingRequest types.BookingRequest
			if err := json.NewDecoder(r.Body).Decode(&bookingRequest); err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid request -%s", err.Error()), http.StatusBadRequest))
				return
			}
			currentUser, ok := r.Context().Value("authorizeduser").(*types.User)
			if !ok {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid token"), http.StatusBadRequest))
				return
			}
			var Booking types.Booking
			var Account types.BankAccount
			if userId, ok := currentUser.ID.(bson.ObjectID); ok {
				Booking.UserId = userId
				// fetch bank
				accountCollecton.FindOne(ctx, bson.M{"userid": userId}).Decode(&Account)
			} else {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid user id"), http.StatusBadRequest))
				return
			}
			parsedBookingStartDate, err := utils.ParseDate(bookingRequest.StartDate)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}

			parsedBookingEndDate, err := utils.ParseDate(bookingRequest.EndDate)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}

			Booking.StartDate = *parsedBookingStartDate
			Booking.EndDate = *parsedBookingEndDate
			Booking.RoomId = bookingRequest.RoomId

			// fetch room
			var Room types.Room
			result := roomCollection.FindOne(ctx, bson.M{"_id": bookingRequest.RoomId}).Decode(&Room)
			if result != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf(constants.InvalidRoomId), http.StatusBadRequest))
				return
			}
			Booking.AmountPaid = float32(Room.Price)
			fmt.Printf("%+v\n", Booking)

			// insert booking
			insertResult, err := bookingCollection.InsertOne(ctx, Booking)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while inserting booking"), http.StatusBadRequest))
				return
			}
			// update room

			updateResult, err := roomCollection.UpdateOne(ctx, bson.M{"_id": bookingRequest.RoomId}, bson.M{
				"$set": bson.M{
					"isbooked":             true,
					"unavailablestartdate": parsedBookingStartDate,
					"unavailableenddate":   parsedBookingEndDate,
				},
			})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while updating room"), http.StatusBadRequest))
				return
			}
			newBalance := Account.Balance - Room.Price
			fmt.Println(Account.ID)
			// update bank user
			updateBankResult, err := accountCollecton.UpdateOne(ctx, bson.M{"_id": Account.ID}, bson.M{
				"$set": bson.M{
					"balance": newBalance,
					"updated": time.Now(),
				},
				"$push": bson.M{
					"transactionhistory": bson.M{
						"spendin":        "room",
						"spendingitemid": bookingRequest.RoomId,
						"amount":         Room.Price,
						"created":        time.Now(),
						"updated":        time.Now(),
					},
				},
			},
			)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while updating Account"), http.StatusBadRequest))
				return
			}
			// update admin account
			adminId, err := bson.ObjectIDFromHex("67e58fd9b260e000f4637ac0")
			if err != nil {
				panic(err)
			}
			// addMoney := Account.Balance + Room.Price
			// fetch the account admin current balance
			accountCollecton.UpdateOne(ctx, bson.M{"_id": adminId}, bson.M{
				"$set": bson.M{
					// "balance": addMoney,
					"updated": time.Now(),
				},
				"$push": bson.M{
					"transactionhistory": bson.M{
						"spendin":        "room",
						"spendingitemid": bookingRequest.RoomId,
						"amount":         Room.Price,
						"created":        time.Now(),
						"updated":        time.Now(),
					},
				},
			},
			)

			fmt.Println(updateResult, insertResult, updateBankResult)
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}
