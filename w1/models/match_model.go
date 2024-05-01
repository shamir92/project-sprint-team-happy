package models

import (
	"errors"
	"gin-mvc/internal"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	MatchStatusPending  = "pending"
	MatchStatusRejected = "rejected"
	MatchStatusApproved = "approved"

	ErrNeitherCatNotFound   = errors.New("neither match/own cat is not found")
	ErrSameGender           = errors.New("can't match cat with same gender")
	ErrCatAlreadyMatched    = errors.New("cat already matched")
	ErrCantMatchYourOwnCats = errors.New("can't match your own cat(s)")
	ErrCantMatchOtherCats   = errors.New("can't match other cat(s)")
)

type MatchError struct {
	Message    string
	StatusCode int
}

func (m MatchError) Error() string {
	return m.Message
}

func (m MatchError) HTTPStatusCode() int {
	return m.StatusCode
}

type Match struct {
	ID            uuid.UUID  `db:"id" json:"id"`                           // UUID primary key
	IssuerID      uuid.UUID  `db:"issuer_id" json:"issuer_id"`             // UUID foreign key to User
	IssuerCatID   uuid.UUID  `db:"issuer_cat_id" json:"issuer_cat_id"`     // UUID foreign key to Cat
	ReceiverID    uuid.UUID  `db:"receiver_id" json:"receiver_id"`         // UUID foreign key to User
	ReceiverCatID uuid.UUID  `db:"receiver_cat_id" json:"receiver_cat_id"` // UUID foreign key to Cat
	Message       string     `db:"message" json:"message"`                 // VARCHAR(120)
	Status        string     `db:"status" json:"status"`                   // VARCHAR(20)
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`           // timestamp with time zone
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`           // timestamp with time zone, nullable
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`           // TIMESTAMP WITH TIME ZONE, nullable
}

type MatchInfo struct {
	ID             uuid.UUID        `db:"id" json:"id"`
	IssuedBy       UserMatch        `json:"issuedBy"`
	MatchCatDetail CatMatchResponse `json:"matchCatDetail"`
	UserCatDetail  CatMatchResponse `json:"userCatDetail"`
	Message        string           `db:"message" json:"message"`      // VARCHAR(120)
	Status         string           `db:"status" json:"status"`        // VARCHAR(20)
	CreatedAt      time.Time        `db:"created_at" json:"createdAt"` // timestamp with time zone
}

type MatchSubmit struct {
	IssuerID      string
	IssuerCatID   string `json:"userCatId"`
	ReceiverID    string
	ReceiverCatID string `json:"matchCatId"`
	Message       string `json:"message" validate:"required, min=5, max=120"`
}

func MatchAll(userID string) ([]MatchInfo, error) {
	db := internal.GetDB()

	var matches []MatchInfo
	rows, err := db.Queryx(`
			SELECT
				m.id,

				i.fullname,
				i.email,
				i.created_at as issuer_created_at,

				receiver_cat_id,
				cr.name as receiver_cat_name,
				cr.race as receiver_cat_race,
				cr.sex as receiver_cat_sex,
				cr.description as receiver_cat_description,
				cr.age_in_month as receiver_cat_age_in_month,
				cr.image_urls as receiver_cat_image_urls,
				cr.has_matched as receiver_cat_has_matched,
				cr.created_at as receiver_cat_created_at,

				issuer_cat_id,
				ci.name as issuer_cat_name,
				ci.race as issuer_cat_race,
				ci.sex as issuer_cat_sex,
				ci.description as issuer_cat_description,
				ci.age_in_month as issuer_cat_age_in_month,
				ci.image_urls as issuer_cat_image_urls,
				ci.has_matched as issuer_cat_has_matched,
				ci.created_at as issuer_cat_created_at,

				message,
				status,
				m.created_at
			FROM matches as m
			LEFT JOIN users as i ON i.id = issuer_id
			LEFT JOIN cats as ci ON ci.id = issuer_cat_id
			LEFT JOIN cats as cr ON cr.id = receiver_cat_id
			WHERE
				(issuer_id = $1 OR receiver_id = $1) AND m.deleted_at IS NULL
	`, userID)

	if err != nil {
		log.Fatalln(err)
		return matches, err
	}

	for rows.Next() {
		var match MatchInfo
		err = rows.Scan(
			&match.ID,
			&match.IssuedBy.FullName,
			&match.IssuedBy.Email,
			&match.IssuedBy.CreatedAt,

			&match.MatchCatDetail.ID,
			&match.MatchCatDetail.Name,
			&match.MatchCatDetail.Race,
			&match.MatchCatDetail.Sex,
			&match.MatchCatDetail.Description,
			&match.MatchCatDetail.AgeInMonth,
			&match.MatchCatDetail.ImageURLs,
			&match.MatchCatDetail.HasMatched,
			&match.MatchCatDetail.CreatedAt,

			&match.UserCatDetail.ID,
			&match.UserCatDetail.Name,
			&match.UserCatDetail.Race,
			&match.UserCatDetail.Sex,
			&match.UserCatDetail.Description,
			&match.UserCatDetail.AgeInMonth,
			&match.UserCatDetail.ImageURLs,
			&match.UserCatDetail.HasMatched,
			&match.UserCatDetail.CreatedAt,

			&match.Message,
			&match.Status,
			&match.CreatedAt,
		)
		if err != nil {
			log.Fatalln(err)
			return []MatchInfo{}, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func MatchCreate(userID string, data MatchSubmit) error {
	db := internal.GetDB()

	err := matchAlreadyExists(data.IssuerCatID, data.ReceiverCatID, db)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO matches(id, issuer_id, receiver_id, issuer_cat_id, receiver_cat_id, message, status, updated_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, NOW())
	`,
		data.IssuerID,
		data.ReceiverID,
		data.IssuerCatID,
		data.ReceiverCatID,
		data.Message,
		MatchStatusPending,
	)

	if err != nil {
		return err
	}

	return nil
}

func matchAlreadyExists(issuerCatId string, receiverCatId string, db *sqlx.DB) error {
	var exist bool
	_ = db.QueryRow("SELECT CASE WHEN EXISTS(SELECT 1 FROM matches WHERE issuer_cat_id = $1 AND receiver_cat_id = $2) THEN true ELSE false END", issuerCatId, receiverCatId).Scan(&exist)

	if exist {
		return MatchError{Message: ErrCatAlreadyMatched.Error(), StatusCode: http.StatusBadRequest}
	}

	return nil
}
