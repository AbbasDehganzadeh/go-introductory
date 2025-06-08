package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"abasd80/link_finder/utils"
)

type DB interface {
	Connect()
}

type DBStruct struct {
	args utils.Arguments
	repo *sqlx.DB
}

const schema = `
  CREATE TABLE IF NOT EXISTS urls (
	id INTEGER PRIMARY KEY,
	origin VARCHAR(50) NOT NULL,
	path VARCHAR(120) NOT NULL,
	status INT DEFAULT 100 CHECK(status IN (0, 1000)),
	reason TEXT DEFAULT "Ok",
	seen BOOLEAN DEFAULT 1 CHECK(seen IN (0, 1)),
	visit BOOLEAN DEFAULT 1 CHECK(visit IN (0, 1)),
	error Text
  );
  CREATE TABLE IF NOT EXISTS links (
	id INTEGER PRIMARY KEY,
	page VARCHAR(150) NOT NULL,
	href VARCHAR(100) NOT NULL,
	title VARCHAR(150) DEFAULT "",
	text TEXT,
	urlid INTEGER,
	FOREIGN KEY (urlid) REFERENCES urls(id)
  );
`

func NewDB(args utils.Arguments) *DBStruct {
	return &DBStruct{args: args}
}

func (db *DBStruct) Connect() {
	if db.repo == nil {
		if db.args.Verbose {
			log.Println("Connecting to Database")
		}
		db_, err := sqlx.Connect(("sqlite3"), "db.sqlite3")
		if err != nil {
			log.Fatalf("DBError: %v\n", err)
		}
		if db.args.Verbose {
			log.Println("Connected successfully!")
		}
		db.repo = db_
	}
	db.repo.MustExec(schema)
}
