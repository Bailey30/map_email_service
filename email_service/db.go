package main

import (
	"context"
	"database/sql"
	"net"
	"os"

	// "errors"
	"fmt"
	// "strings"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	// "cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	// "github.com/jmoiron/sqlx"
)

type PasswordResetDB struct {
	databaseUrl string
	dbx         *sql.DB
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

func (db *PasswordResetDB) connect() (*sql.DB, error) {
	var (
		dbUser                 = os.Getenv("DB_USER") // e.g. 'service-account-name@project-id.iam'
		dbPwd                  = os.Getenv("DB_PASS")
		dbName                 = os.Getenv("DB_NAME") // e.g. 'my-database'
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	)

	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.DialOption
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil

	// cleanup, err := pgxv4.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
	// if err != nil {
	// 	return nil, fmt.Errorf("Error on pgxv4.RegisterDriver: %v", err)
	// }
	//
	// dbx, err := sql.Open("cloudsql-postgres", db.databaseUrl)
	// if err != nil {
	// 	return nil, fmt.Errorf("error connecting to db: sqlx.Connect: %v", err)
	// }
	//
	// fmt.Println("db:", dbx)
	// db.dbx = dbx
	//
	// fmt.Println("Connected to PasswordReset database")
	//
	// return cleanup, err
}

func (db *PasswordResetDB) close() error {
	return db.dbx.Close()
}

func (db *PasswordResetDB) create(params CreateResetCodeParams) error {
	cleanup, err := db.connect()
	if err != nil {
		return err
	}
	defer cleanup.Close()
	// cleanup.Close()
	// defer cleanup()

	// reset_code := ResetCode{
	// 	UserId:     params.UserId,
	// 	HashedCode: params.HashedCode,
	// 	Expiry:     time.Now().Add(time.Hour),
	// }

	// if _, err := db.dbx.NamedExec(
	// 	`INSERT INTO resetcode
	//                (user_id, hashed_code, expiry)
	//        VALUES
	//        (:user_id, :hashed_code, :expiry)`,
	// 	reset_code); err != nil {
	// 	if strings.Contains(err.Error(), "Error 1062") {
	// 		return fmt.Errorf("Duplicate key error: %v", 1)
	// 	}
	//
	// 	return err
	// }

	return nil
}

func (db *PasswordResetDB) getOneById() (ResetCode, error) {
	cleanup, err := db.connect()
	if err != nil {
		return ResetCode{}, err
	}
	defer cleanup.Close()
	// defer db.close()
	// defer cleanup()

	var reset_code ResetCode
	// if err := db.dbx.Get(
	// 	&reset_code,
	// 	`SELECT
	//            UserId, HashedCode
	//        FROM resetcode
	//        WHERE UserId = $1`,
	// 	1); err != nil {
	// 	if err != sql.ErrNoRows {
	// 		return ResetCode{}, err
	// 	}
	// 	return ResetCode{}, errors.New("Record not found")
	// }

	return reset_code, err
}

// gets all the reset codes associated with the user id
func (db *PasswordResetDB) getAllById(id int) ([]ResetCode, error) {
	cleanup, err := db.connect()
	if err != nil {
		return []ResetCode{}, err
	}
	defer cleanup.Close()
	// defer db.close()
	// defer cleanup()

	var resetCodes []ResetCode
	// if err := db.dbx.Select(
	// 	&resetCodes,
	// 	`SELECT * FROM resetcode WHERE user_id = $1`,
	// 	id); err != nil {
	// 	return nil, err
	// }

	return resetCodes, nil
}

// deletes all reset codes associated with the user id
func (db *PasswordResetDB) deleteAllById(id int) error {
	cleanup, err := db.connect()
	if err != nil {
		return err
	}
	defer cleanup.Close()
	// defer db.close()
	// defer cleanup()

	if _, err := db.dbx.Exec(
		`DELETE FROM resetcode WHERE user_id = $1`, id); err != nil {
		return err
	}

	return nil
}

func (db *PasswordResetDB) getByToken(token string) (ResetCode, error) {
	cleanup, err := db.connect()
	if err != nil {
		return ResetCode{}, err
	}
	defer cleanup.Close()
	// defer db.close()
	// defer cleanup()

	var resetCode ResetCode
	// if err = db.dbx.Get(&resetCode, `SELECT * FROM resetcode where hashed_code = $1 LIMIT 1`, token); err != nil {
	// 	return ResetCode{}, err
	// }

	return resetCode, nil
}

func (db *PasswordResetDB) CreateTable() error {
	fmt.Println("create table")
	cleanup, err := db.connect()
	if err != nil {
		return fmt.Errorf("Error connecting to db: %v", err)
	}
	defer cleanup.Close()
	// defer db.close()
	// defer cleanup()

	var schema = `
	   CREATE TABLE IF NOT EXISTS resetcode (
	       user_id integer,
	       hashed_code text,
	       expiry timestamp
	   )
	`
	// execute a query on the server. MustExec panics on error
	_, err = cleanup.Exec(schema)
	if err != nil {
		return fmt.Errorf("Error creating table: %v", err)
	}

	return nil
}
