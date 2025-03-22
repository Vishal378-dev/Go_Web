package Account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	utils "github.com/vishal/reservation_system/Handlers/Utils"
	Handlers "github.com/vishal/reservation_system/Handlers/WrongPath"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AccountHandler(accountCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := utils.Ctx(5)
		defer cancel()
		if r.Method == "GET" {
			var Account types.BankAccount
			id := r.URL.Query().Get("id")
			if id == "" {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("missing id"), http.StatusBadRequest))
				return
			}
			qid, err := bson.ObjectIDFromHex(id)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid id"), http.StatusBadRequest))
				return
			}
			filter := bson.M{"_id": qid}
			err = accountCollection.FindOne(ctx, filter).Decode(&Account)

			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("no record found"), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, Account)
			return
		} else if r.Method == "POST" {
			var BankAccountRequest *types.BankAccount
			if err := json.NewDecoder(r.Body).Decode(&BankAccountRequest); err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("invalid request -%s", err.Error()), http.StatusBadRequest))
				return
			}
			err := BankAccountRequest.ValidateRequest()
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(err, http.StatusBadRequest))
				return
			}
			BankAccountRequest.Created = time.Now()
			BankAccountRequest.Updated = time.Now()
			result, err := accountCollection.InsertOne(ctx, BankAccountRequest)
			if err != nil {
				utils.ResponseWriter(w, http.StatusBadRequest, utils.CommonError(fmt.Errorf("error while inserting the data"), http.StatusBadRequest))
				return
			}
			utils.ResponseWriter(w, http.StatusOK, result)
			return
		} else {
			Handlers.WrongPathTemplate(w, r)
		}
	}
}
