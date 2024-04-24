package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type PasswordResetDB struct {
	databaseUrl string
	dbx         *sqlx.DB
}

type ResetCode struct {
	UserId     int       `db:"user_id"`
	HashedCode string    `db:"hashed_code"`
	Expiry     time.Time `db:"expiry"`
}

type CreateResetCodeParams struct {
	UserId     int
	HashedCode string
}

func NewPasswordResetDB(databaseUrl string) *PasswordResetDB {
	return &PasswordResetDB{
		databaseUrl: databaseUrl,
	}
}

func (db *PasswordResetDB) connect() error {
	dbx, err := sqlx.Connect("postgres", db.databaseUrl)
	if err != nil {
		return err
	}

	db.dbx = dbx

	fmt.Println("Connected to PasswordReset database")

	return nil
}

func (db *PasswordResetDB) close() error {
	return db.dbx.Close()
}

func (db *PasswordResetDB) create(params CreateResetCodeParams) error {
	err := db.connect()
	if err != nil {
		return err
	}
	defer db.close()

	reset_code := ResetCode{
		UserId:     params.UserId,
		HashedCode: params.HashedCode,
		Expiry:     time.Now().Add(time.Hour),
	}

	if _, err := db.dbx.NamedExec(
		`INSERT INTO resetcode
                (user_id, hashed_code, expiry)
        VALUES 
        (:user_id, :hashed_code, :expiry)`,
		reset_code); err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return fmt.Errorf("Duplicate key error: %v", 1)
		}

		return err
	}

	return nil
}

func (db *PasswordResetDB) getOneById() (ResetCode, error) {
	err := db.connect()
	if err != nil {
		return ResetCode{}, err
	}
	defer db.close()

	var reset_code ResetCode
	if err := db.dbx.Get(
		&reset_code,
		`SELECT
            UserId, HashedCode
        FROM resetcode 
        WHERE UserId = $1`,
		1); err != nil {
		if err != sql.ErrNoRows {
			return ResetCode{}, err
		}
		return ResetCode{}, errors.New("Record not found")
	}

	return reset_code, err
}

// gets all the reset codes associated with the user id
func (db *PasswordResetDB) getAllById(id int) ([]ResetCode, error) {
	err := db.connect()
	if err != nil {
		return []ResetCode{}, err
	}
	defer db.close()

	var resetCodes []ResetCode
	if err := db.dbx.Select(
		&resetCodes,
		`SELECT * FROM resetcode WHERE user_id = $1`,
		id); err != nil {
		return nil, err
	}

	return resetCodes, nil
}

// deletes all reset codes associated with the user id
func (db *PasswordResetDB) deleteAllById(id int) error {
	err := db.connect()
	if err != nil {
		return err
	}
	defer db.close()

	if _, err := db.dbx.Exec(
		`DELETE FROM resetcode WHERE user_id = $1`, id); err != nil {
		return err
	}

	return nil
}

func (db *PasswordResetDB) getByToken(token string) (ResetCode, error) {
	err := db.connect()
	if err != nil {
		return ResetCode{}, err
	}
	defer db.close()

	var resetCode ResetCode
	if err = db.dbx.Get(&resetCode, `SELECT * FROM resetcode where hashed_code = $1 LIMIT 1`, token); err != nil {
		return ResetCode{}, err
	}

	return resetCode, nil
}

func (db *PasswordResetDB) CreateTable() error {
	err := db.connect()
	if err != nil {
		return err
	}
	defer db.close()
	var schema = `
    CREATE TABLE IF NOT EXISTS resetcode (
        user_id integer,
        hashed_code text,
        expiry timestamp
    )
    `
	// execute a query on the server. MustExec panics on error
	db.dbx.MustExec(schema)

	return nil
}
