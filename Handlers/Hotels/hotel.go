package Hotels

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/dummy"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Hotel(hotelCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			var hotels []types.Hotel
			response, err := hotelCollection.Find(context.TODO(), bson.D{})
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
				return
			}
			err = response.All(context.TODO(), &hotels)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing entries"), http.StatusBadRequest))
				return
			}

			Response := types.SuccessResponse{
				Status: http.StatusOK,
				Data:   hotels,
			}
			utils.ResponseWriter(w, http.StatusOK, Response)
			return
		} else if r.Method == "POST" {
			currentUser, ok := r.Context().Value("authorizeduser").(*types.User)
			if !ok {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing the data"), http.StatusBadRequest))
				return
			}
			if currentUser.Role != "Admin" {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("only admin is allowed to perform this action"), http.StatusBadRequest))
				return
			}
			ctx, cancel := utils.Ctx(5)
			defer cancel()
			var hotelBody types.Hotel
			err := json.NewDecoder(r.Body).Decode(&hotelBody)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			fmt.Println(hotelBody.Address)
			res, err := hotelCollection.InsertOne(ctx, hotelBody)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusCreated, map[string]any{"hotel": res})
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}

func HotelById(hotelCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qid := strings.TrimPrefix(r.URL.Path, "/hotel/")
		if qid == "" {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("invalid id"), http.StatusBadRequest))
			return
		}
		id, err := bson.ObjectIDFromHex(qid)
		if err != nil {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
			return
		}
		if r.Method == "GET" {
			isRoom := r.URL.Query().Get("isRoom")
			var hotels types.Hotel
			if isRoom == "true" {
				pipeline := bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
					bson.M{
						"$lookup": bson.M{
							"from":         "Rooms",
							"localField":   "rooms",
							"foreignField": "_id",
							"as":           "roomDetails",
						},
					},
					// bson.M{      // uncomment if you want to have each room detail as an object
					// 	"$unwind": bson.M{
					// 		"path":                       "$roomDetails",
					// 		"preserveNullAndEmptyArrays": true,
					// 	},
					// },
					bson.M{
						"$project": bson.M{
							"Address":      1,
							"address":      1,
							"amenities":    1,
							"description":  1,
							"hotelid":      1,
							"name":         1,
							"review":       1,
							"star":         1,
							"typesofrooms": 1,
							"roomDetails":  1,
						},
					},
				}
				cursor, err := hotelCollection.Aggregate(context.TODO(), pipeline)
				if err != nil {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
					return
				}
				defer cursor.Close(context.TODO())

				var results []bson.M
				if err = cursor.All(context.TODO(), &results); err != nil {
					log.Fatal(err)
				}
				Response := types.SuccessResponse{
					Status: http.StatusOK,
					Data:   results,
				}
				utils.ResponseWriter(w, http.StatusOK, Response)
				return

			} else {
				err := hotelCollection.FindOne(context.TODO(),
					bson.M{"_id": id}).Decode(&hotels)

				if err != nil {
					utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
					return
				}
			}

			Response := types.SuccessResponse{
				Status: http.StatusOK,
				Data:   hotels,
			}
			utils.ResponseWriter(w, http.StatusOK, Response)
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}
func UpdateHotelById(hotelCollection *mongo.Collection) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := r.Context().Value("authorizeduser").(*types.User)
		if !ok {
			utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while parsing the data"), http.StatusBadRequest))
			return
		}
		if currentUser.Role != "Admin" {
			utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("only admin is allowed to perform this action"), http.StatusBadRequest))
			return
		}
		qid := strings.TrimPrefix(r.URL.Path, "/hotels/")
		if qid == "" {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("invalid id"), http.StatusBadRequest))
			return
		}
		id, err := bson.ObjectIDFromHex(qid)
		if err != nil {
			utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(err, http.StatusBadRequest))
			return
		}
		ctx, cancel := utils.Ctx(5)
		defer cancel()
		if r.Method == "PUT" {
			var hotel types.Hotel
			err := json.NewDecoder(r.Body).Decode(&hotel)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest))
				return
			}
			var existingHotel types.Hotel
			err = hotelCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingHotel)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					utils.ResponseWriter(w, http.StatusNotFound, utils.CommonError(fmt.Errorf("no record found with id: %v", id), http.StatusNotFound))
					return
				}
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error fetching hotel: %v", err), http.StatusInternalServerError))
				return
			}
			updateFields := bson.M{}
			if hotel.Name != "" {
				updateFields["name"] = hotel.Name
			}
			if hotel.Description != "" {
				updateFields["description"] = hotel.Description
			}
			if hotel.Star != 0 {
				updateFields["star"] = hotel.Star
			}
			if hotel.Review != nil {
				updateFields["review"] = hotel.Review
			}
			if hotel.Amenities != nil {
				updateFields["amenities"] = hotel.Amenities
			}
			if hotel.AdditionalInfo1 != "" {
				updateFields["additionalinfo1"] = hotel.AdditionalInfo1
			}
			if hotel.AdditionalInfo2 != "" {
				updateFields["additionalinfo2"] = hotel.AdditionalInfo2
			}
			if hotel.AdditionalInfo3 != nil {
				updateFields["additionalinfo3"] = hotel.AdditionalInfo3
			}
			if hotel.TypesOfRooms != nil {
				updateFields["typesofrooms"] = hotel.TypesOfRooms
			}
			if hotel.Address.LandMark != "" {
				updateFields["address.landmark"] = hotel.Address.LandMark
			}
			if hotel.Address.City != "" {
				updateFields["address.city"] = hotel.Address.City
			}
			if hotel.Address.State != "" {
				updateFields["address.state"] = hotel.Address.State
			}
			if hotel.Address.Street != "" {
				updateFields["address.street"] = hotel.Address.Street
			}
			if hotel.Address.Pincode != 0 {
				updateFields["address.pincode"] = hotel.Address.Pincode
			}
			if hotel.Address.Coordinates.Latitude != 0 && hotel.Address.Coordinates.Longitude != 0 {
				updateFields["address.coordinates"] = bson.M{
					"latitude":  hotel.Address.Coordinates.Latitude,
					"longitude": hotel.Address.Coordinates.Longitude,
				}
			} else if hotel.Address.Coordinates.Latitude != 0 {
				updateFields["address.coordinates"] = bson.M{
					"latitude":  hotel.Address.Coordinates.Latitude,
					"longitude": existingHotel.Address.Coordinates.Longitude,
				}
			} else if hotel.Address.Coordinates.Longitude != 0 {
				updateFields["address.coordinates"] = bson.M{
					"longitude": hotel.Address.Coordinates.Longitude,
					"latitude":  existingHotel.Address.Coordinates.Latitude,
				}
			}
			res, err := hotelCollection.UpdateOne(
				ctx,
				bson.M{"_id": id},
				bson.M{"$set": updateFields},
			)
			if err != nil {
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error updating hotel: %v", err), http.StatusInternalServerError))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{
				"msg":    "Hotel updated successfully",
				"result": res,
			})
		} else if r.Method == "DELETE" {
			filter := bson.D{{Key: "_id", Value: id}}
			res, err := hotelCollection.DeleteOne(ctx, filter)
			if res.DeletedCount == 0 {
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("no record available to delete"), http.StatusInternalServerError))
				return
			}
			if err != nil {
				utils.ResponseWriter(w, http.StatusInternalServerError, utils.CommonError(fmt.Errorf("error deleting hotel: %v", err), http.StatusInternalServerError))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, map[string]interface{}{
				"msg": "Successfully Deleted the Hotel",
			})
		} else {
			Handlers.WrongPathTemplate(w, r)
			return
		}
	}
}
