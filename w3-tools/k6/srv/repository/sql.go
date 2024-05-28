package repository

import (
	"context"
	"database/sql"
	"math/rand/v2"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) PostUsedITAccount(ctx context.Context, usr *entity.UsedUser) error {
	_, err := r.db.Exec("INSERT INTO used_it_account (nip, password) VALUES (?, ?)", usr.Nip, usr.Password)
	return err
}
func (r *Repository) PostUsedNurseAccount(ctx context.Context, usr *entity.UsedUser) error {
	_, err := r.db.Exec("INSERT INTO used_nurse_account (nip, password) VALUES (?, ?)", usr.Nip, usr.Password)
	return err
}
func (r *Repository) GetUsedITAccount(ctx context.Context) (*entity.UsedUser, error) {
	row := r.db.QueryRow("SELECT value FROM meta_data WHERE key = 'itIndex'")

	var itIndex int
	row.Scan(&itIndex)

	if itIndex == 0 {
		return nil, nil
	}

	usedUser := &entity.UsedUser{}
	row = r.db.QueryRow("SELECT nip, password FROM used_it_account WHERE id = ?", rand.IntN(itIndex)+1)
	row.Scan(&usedUser.Nip, &usedUser.Password)

	return usedUser, nil
}

func (r *Repository) GetUsedNurseAccount(ctx context.Context) (*entity.UsedUser, error) {
	row := r.db.QueryRow("SELECT value FROM meta_data WHERE key = 'nurseIndex'")

	var itIndex int
	row.Scan(&itIndex)
	if itIndex == 0 {
		return nil, nil
	}

	usedUser := &entity.UsedUser{}
	row = r.db.QueryRow("SELECT nip, password FROM used_nurse_account WHERE id = ?", rand.IntN(itIndex)+1)
	row.Scan(&usedUser.Nip, &usedUser.Password)

	return usedUser, nil
}

func (r *Repository) Reset(ctx context.Context) error {
	_, err := r.db.Exec(`DROP TABLE IF EXISTS used_it_account`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DROP TABLE IF EXISTS used_nurse_account`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DROP TABLE IF EXISTS meta_data`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS used_it_account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS used_nurse_account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS meta_data (
		key TEXT PRIMARY KEY,
		value INT
	);`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT OR IGNORE INTO meta_data (key, value) VALUES ('itIndex', 0)`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT OR IGNORE INTO meta_data (key, value) VALUES ('nurseIndex', 0)`)
	if err != nil {
		return err
	}

	return nil
}
