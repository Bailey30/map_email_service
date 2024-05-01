package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"strings"

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

func (db *PasswordResetDB) connect() error {
	var (
		dbUser                 = os.Getenv("DB_USER") // e.g. 'service-account-name@project-id.iam'
		dbPwd                  = os.Getenv("DB_PASS")
		dbName                 = os.Getenv("DB_NAME") // e.g. 'my-database'
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	)

	// Creating a new Cloud SQL dialer.
	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	// Creating a Data Source Name (DSN) for PostgreSQL connection.
	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)

	// Parsing the DSN into a PostgreSQL configuration.
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}

	var opts []cloudsqlconn.DialOption

	// Customizing the Dial function of the configuration to use the Cloud SQL dialer.
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}

	// Registering the connection configuration with the standard library.
	dbURI := stdlib.RegisterConnConfig(config)

	// Opening a new database connection pool using pgx driver.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	// Assigning the connection pool to a struct field.
	db.dbx = dbPool

	// Returning nil, indicating successful setup of the database connection.
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

	if _, err := db.dbx.Exec(
		`INSERT INTO resetcode (user_id, hashed_code, expiry) VALUES ($1, $2, $3)`,
		params.UserId, params.HashedCode, time.Now().Add(time.Hour),
	); err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return fmt.Errorf("Duplicate key error: %v", 1)
		}

		return fmt.Errorf("Error storing resetcode: %v", err)
	}

	return nil
}

// UNUSED
func (db *PasswordResetDB) getOneById() (ResetCode, error) {
	err := db.connect()
	if err != nil {
		return ResetCode{}, err
	}
	defer db.close()
	// defer cleanup()

	var resetCode ResetCode
	// Prepare the SQL query
	query := `
        SELECT
            UserId, HashedCode
        FROM resetcode
        WHERE UserId = $1
    `

	// Execute the query
	row := db.dbx.QueryRow(query, 1)
	err = row.Scan(&resetCode.UserId, &resetCode.HashedCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return ResetCode{}, errors.New("Record not found")
		}
		return ResetCode{}, err
	}
	return resetCode, err
}

// gets all the reset codes associated with the user id
func (db *PasswordResetDB) getAllById(id int) ([]ResetCode, error) {
	err := db.connect()
	if err != nil {
		return []ResetCode{}, err
	}
	defer db.close()

	var resetCodes []ResetCode

	rows, err := db.dbx.Query("SELECT * FROM resetcode WHERE user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Error querying for resetcodes: %v", err)
	}

	for rows.Next() {
		var resetCode ResetCode
		if err := rows.Scan(&resetCode.UserId, &resetCode.HashedCode, &resetCode.Expiry); err != nil {
			return nil, fmt.Errorf("Error scanning rows when getting all by id: %v", err)
		}
		resetCodes = append(resetCodes, resetCode)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error scanning rows when getting all by id: %v", err)
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

	row := db.dbx.QueryRow("SELECT * FROM resetcode where hashed_code = $1 LIMIT 1", token)
	if err := row.Scan(&resetCode.UserId, &resetCode.HashedCode, &resetCode.Expiry); err != nil {
		if err == sql.ErrNoRows {
			return ResetCode{}, fmt.Errorf("Invalid password reset token")
		}
		return ResetCode{}, fmt.Errorf("Error scanning rows")
	}

	return resetCode, nil
}

func (db *PasswordResetDB) CreateTable() error {
	fmt.Println("create table")
	err := db.connect()
	if err != nil {
		return fmt.Errorf("Error connecting to db: %v", err)
	}
	defer db.close()

	var schema = `
	   CREATE TABLE IF NOT EXISTS resetcode (
	       user_id integer,
	       hashed_code text,
	       expiry timestamp
	   )
	`
	_, err = db.dbx.Exec(schema)
	if err != nil {
		return fmt.Errorf("Error creating table: %v", err)
	}

	return nil
}
