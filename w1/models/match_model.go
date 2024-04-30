package models

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID            uuid.UUID  `db:"id" json:"id"`                           // UUID primary key
	IssuerCatID   uuid.UUID  `db:"issuer_cat_id" json:"issuer_cat_id"`     // UUID foreign key to Cat
	ReceiverCatID uuid.UUID  `db:"receiver_cat_id" json:"receiver_cat_id"` // UUID foreign key to Cat
	IssuerID      uuid.UUID  `db:"issuer_id" json:"issuer_id"`             // UUID foreign key to User
	ReceiverID    uuid.UUID  `db:"receiver_id" json:"receiver_id"`         // UUID foreign key to User
	Message       string     `db:"message" json:"message"`                 // VARCHAR(120)
	Status        string     `db:"status" json:"status"`                   // VARCHAR(20)
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`           // timestamp with time zone
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`           // timestamp with time zone, nullable
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`           // TIMESTAMP WITH TIME ZONE, nullable
}
