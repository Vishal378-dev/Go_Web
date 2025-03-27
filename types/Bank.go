package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	SpendIn        string          `json:"spendin" bson:"spendin"`
	SpendingItemId []bson.ObjectID `json:"spendingitemid" bson:"spendingitemid"`
	Amount         int             `json:"amount" bson:"amount"`
	Created        time.Time       `json:"created_at" bson:"created_at"`
	Updated        time.Time       `json:"updated_at" bson:"updated_at"`
}

type BankAccount struct {
	ID                  bson.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	BankName            string         `json:"bankname" bson:"bankname"`
	AccountNumber       int            `json:"accountnumber" bson:"accountnumber"`
	BankIfsc            string         `json:"bankifsc" bson:"bankifsc"`
	BankHolderFirstName string         `json:"bankholderfirstname" bson:"bankholderfirstname"`
	BankHolderLastName  string         `json:"bankholderlastname" bson:"bankholderlastname"`
	Balance             float32        `json:"balance" bson:"balance"`
	TransactionHistory  []Transaction  `json:"transactionhistory,omitempty" bson:"transactionhistory,omitempty"`
	UserId              *bson.ObjectID `json:"userid" bson:"userid"`
	Created             time.Time      `json:"created_at" bson:"created_at"`
	Updated             time.Time      `json:"updated_at" bson:"updated_at"`
}

func (ba *BankAccount) ValidateRequest() error {
	if len(ba.BankName) < 3 {
		return fmt.Errorf("invalid bankname")
	}
	if ba.AccountNumber <= 9999 {
		return fmt.Errorf("invalid account number")
	}
	if ba.BankIfsc == "" {
		return fmt.Errorf("invalid ifsc")
	}
	if len(ba.BankHolderFirstName) < 3 {
		return fmt.Errorf("invalid first name")
	}
	if len(ba.BankHolderLastName) < 3 {
		return fmt.Errorf("invalid last name")
	}
	if ba.UserId == nil {
		return fmt.Errorf("invalid user id")
	}
	return nil
}
