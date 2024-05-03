package models

import (
	"database/sql"
	"errors"
	"gin-mvc/internal"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	MatchStatusPending  = "pending"
	MatchStatusRejected = "rejected"
	MatchStatusApproved = "approved"

	ErrSameGender                 = errors.New("can't match cats with same gender")
	ErrCantMatchAlreadyMatchedCat = errors.New("cats already matched")
	ErrCantMatchOwnedCats         = errors.New("can't match owned cats")
	ErrMatchNotFound              = errors.New("match not found")
	ErrMatchIsNoLongerValid       = errors.New("match is no longer valid")
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

type MatchCatDetail struct {
	ID          uuid.UUID      `db:"id" json:"id"`                   // UUID primary key
	Name        string         `db:"name" json:"name"`               // VARCHAR(30)
	Sex         string         `db:"sex" json:"sex"`                 // VARCHAR(10)
	AgeInMonth  int            `db:"age_in_month" json:"ageInMonth"` // INT
	Description string         `db:"description" json:"description"` // VARCHAR(20)
	ImageURLs   pq.StringArray `db:"image_urls" json:"imageUrls"`    // TEXT[], array of strings in Go
	Race        string         `db:"race" json:"race"`               // VARCHAR(50)
	HasMatched  bool           `db:"has_matched" json:"hasMatched"`  // BOOLEAN with default false
	CreatedAt   time.Time      `db:"created_at" json:"createdAt"`    // timestamp with time zone
}

type MatchInfo struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	IssuedBy       UserMatch      `json:"issuedBy"`
	MatchCatDetail MatchCatDetail `json:"matchCatDetail"`
	UserCatDetail  MatchCatDetail `json:"userCatDetail"`
	Message        string         `db:"message" json:"message"`      // VARCHAR(120)
	Status         string         `db:"status" json:"status"`        // VARCHAR(20)
	CreatedAt      time.Time      `db:"created_at" json:"createdAt"` // timestamp with time zone
}

type MatchCreateIn struct {
	IssuerCatID   string
	ReceiverCatID string
	Message       string
}

type MatchAnswerIn struct {
	MatchID string
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

func MatchCreate(userID string, data MatchCreateIn) error {
	db := internal.GetDB()

	issuerCat, err := GetCayById(data.IssuerCatID, db)
	if err != nil {
		return err
	}

	if err := CheckIfCatIsOwnedByUser(issuerCat.OwnerID.String(), userID); err != nil {
		return err
	}

	receiverCat, err := GetCayById(data.ReceiverCatID, db)
	if err != nil {
		return err
	}

	if issuerCat.OwnerID.String() == receiverCat.OwnerID.String() {
		return MatchError{Message: ErrCantMatchOwnedCats.Error(), StatusCode: http.StatusBadRequest}
	}

	err = matchCheckIfSameSex(issuerCat.Sex, receiverCat.Sex)
	if err != nil {
		return err
	}

	if issuerCat.HasMatched || receiverCat.HasMatched {
		return MatchError{Message: ErrCantMatchAlreadyMatchedCat.Error(), StatusCode: http.StatusBadRequest}
	}

	err = matchAlreadyExists(data.IssuerCatID, data.ReceiverCatID, db)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO matches(id, issuer_id, receiver_id, issuer_cat_id, receiver_cat_id, message, status, updated_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, NOW())
	`,
		issuerCat.OwnerID,
		receiverCat.OwnerID,
		issuerCat.ID,
		receiverCat.ID,
		data.Message,
		MatchStatusPending,
	)

	if err != nil {
		return err
	}

	return nil
}

func MatchApprove(userID string, data MatchAnswerIn) error {
	db := internal.GetDB()

	match, err := matchGetByIdAndReceiverId(data.MatchID, userID, db)
	if err != nil {
		return err
	}

	err = matchUpdateStatus(match.ID.String(), MatchStatusApproved, db)
	if err != nil {
		return err
	}

	err = UpdateHasMatchedCat(
		[]string{match.IssuerCatID.String(), match.ReceiverCatID.String()},
		[]string{match.IssuerID.String(), match.ReceiverID.String()},
		db,
	)

	if err != nil {
		return err
	}

	err = matchRemoveRelatedMatch(match.ID.String(), match.ReceiverCatID.String(), db)

	return err
}

func MatchReject(userId string, data MatchAnswerIn) error {
	db := internal.GetDB()

	match, err := matchGetByIdAndReceiverId(data.MatchID, userId, db)
	if err != nil {
		return err
	}

	if match.Status == MatchStatusApproved {
		return MatchError{Message: ErrMatchIsNoLongerValid.Error(), StatusCode: http.StatusBadRequest}
	}

	err = matchUpdateStatus(match.ID.String(), MatchStatusRejected, db)
	if err != nil {
		return err
	}

	return nil
}

func matchUpdateStatus(matchId string, status string, db *sqlx.DB) error {
	query := `UPDATE matches SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, matchId)
	return err
}

func matchRemoveRelatedMatch(matchId string, receiverCatId string, db *sqlx.DB) error {
	query := `UPDATE matches SET status = $1, updated_at = NOW(), deleted_at = NOW() WHERE id != $2 AND receiver_cat_id = $3`
	_, err := db.Exec(query, MatchStatusRejected, matchId, receiverCatId)
	return err
}

func matchGetByIdAndReceiverId(matchId string, receiverId string, db *sqlx.DB) (Match, error) {
	var match Match

	query := `
		SELECT id, issuer_id, issuer_cat_id, receiver_id, receiver_cat_id, message, status, deleted_at
		FROM matches
		WHERE id = $1 AND receiver_id = $2 AND deleted_at is null
	`

	err := db.QueryRowx(query, matchId, receiverId).Scan(
		&match.ID,
		&match.IssuerID,
		&match.IssuerCatID,
		&match.ReceiverID,
		&match.ReceiverCatID,
		&match.Message,
		&match.Status,
		&match.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return match, MatchError{Message: ErrMatchNotFound.Error(), StatusCode: http.StatusNotFound}
		}

		return match, err
	}

	if match.Status == MatchStatusRejected {
		return match, MatchError{Message: ErrMatchIsNoLongerValid.Error(), StatusCode: http.StatusBadRequest}
	}

	return match, nil
}

func matchCheckIfSameSex(issuerSex string, receiverSex string) error {
	if issuerSex == receiverSex {
		return MatchError{Message: ErrSameGender.Error(), StatusCode: http.StatusBadRequest}
	}

	return nil
}

func matchAlreadyExists(issuerCatId string, receiverCatId string, db *sqlx.DB) error {
	var exist bool
	_ = db.QueryRow("SELECT CASE WHEN EXISTS(SELECT 1 FROM matches WHERE issuer_cat_id = $1 AND receiver_cat_id = $2) THEN true ELSE false END", issuerCatId, receiverCatId).Scan(&exist)

	if exist {
		return MatchError{Message: ErrCantMatchAlreadyMatchedCat.Error(), StatusCode: http.StatusBadRequest}
	}

	return nil
}
