// This file contains types that are used throughout the project
package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeTransferOut TransactionType = "TransferOut"
	TransactionTypeTopUp       TransactionType = "TopUp"
)

type TransactionStatus string

const (
	TransactionStatusSuccessful TransactionStatus = "Successful"
	TransactionStatusFailed     TransactionStatus = "Failed"
)

type User struct {
	ID          int64      `db:"id"`
	FullName    string     `db:"full_name"`
	PhoneNumber string     `db:"phone_number"`
	Balance     float32    `db:"balance"`
	Password    string     `db:"password"`
	CreatedTime time.Time  `db:"created_time"`
	UpdatedTime *time.Time `db:"updated_time"`
}

type UserFilter struct {
	UserID      int64  `db:"user_id"`
	PhoneNumber string `db:"phone_number"`
}

type Transaction struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	UserID      int64             `json:"user_id" db:"user_id"`
	Amount      float32           `json:"amount" db:"amount"`
	Type        TransactionType   `json:"type" db:"type"`
	RecipientID int64             `json:"recipient_id,omitempty" db:"recipient_id"` // Pointer to handle NULL values
	CreatedTime time.Time         `json:"created_time" db:"created_time"`
	UpdatedTime *time.Time        `json:"updated_time,omitempty" db:"updated_time"` // Pointer to handle NULL values
	Status      TransactionStatus `json:"status" db:"status"`
	Description string            `json:"description" db:"description"`
	Password    string
}
