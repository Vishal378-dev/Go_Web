package rooms

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/dummy"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GETRooms(roomCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := utils.Ctx(5)
		defer cancel()
		if r.Method == "GET" {
			cursor, err := roomCollection.Find(ctx, bson.M{})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found - %s", err.Error()), http.StatusBadRequest))
				return
			}
			var rooms *types.Room
			result := cursor.All(ctx, &rooms)
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{
				"msg":  "successfully fetched the data",
				"data": result,
			})
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}

func POSTRooms(roomCollection, hotelCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := utils.Ctx(5)
		defer cancel()
		if r.Method == "POST" {
			var Room *types.Room
			currentUser, ok := r.Context().Value("authorizeduser").(*types.User)
			if !ok {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing the data"), http.StatusBadRequest))
				return
			}
			if currentUser.Role != "Admin" {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("only admin is allowed to perform this action"), http.StatusBadRequest))
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&Room); err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid request -%s", err.Error()), http.StatusBadRequest))
				return
			}
			err := Room.RequestValidation()
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			id, err := bson.ObjectIDFromHex(Room.HotelID.Hex())
			if err != nil {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			var FetchedRoom types.Room
			roomCollection.FindOne(ctx, bson.M{"roomnumber": Room.RoomNumber}).Decode(&FetchedRoom)
			if FetchedRoom.RoomNumber == Room.RoomNumber {
				utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("room number is already present"), http.StatusBadRequest))
				return
			}
			result, err := roomCollection.InsertOne(ctx, Room)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while inserting room -%s", err.Error()), http.StatusBadRequest))
				return
			}
			result2, err := hotelCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.M{"$push": bson.M{"rooms": result.InsertedID}})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while inserting room -%s", err.Error()), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{
				"msg":   "successfully fetched the data",
				"data":  result,
				"data2": result2,
			})
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}
