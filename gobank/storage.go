package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank host=54.250.166.42 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {

	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {

	query := `create table  if no exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		create_at timestamp
		 )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(*Account) error {

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {

	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {

	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	return nil, nil
}
