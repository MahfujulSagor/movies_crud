package sqlite

import (
	"database/sql"
	"github/MahfujulSagor/movies_crud/internals/config"
	"github/MahfujulSagor/movies_crud/internals/types"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS directors(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS casts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		actor TEXT NOT NULL,
		actress TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		rating INTEGER NOT NULL,
		director_id INTEGER,
		cast_id INTEGER,
		FOREIGN KEY (director_id) REFERENCES directors(id),
		FOREIGN KEY (cast_id) REFERENCES casts(id)
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLite{
		DB: db,
	}, nil
}

func (s *SQLite) CreateMovie(title string, rating int, director *types.Director, cast *types.Cast) (int64, error) {
	//? Transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	//* Directors
	//? Execute statement
	dir_res, err := tx.Exec("INSERT INTO directors(name, age) VALUES (?, ?)", director.Name, director.Age)
	if err != nil {
		return 0, err
	}

	//? Get the last inserted ID
	director_id, err := dir_res.LastInsertId()
	if err != nil {
		return 0, err
	}

	//* Casts
	//? Execute statement
	cast_res, err := tx.Exec("INSERT INTO casts(actor, actress) VALUES (?, ?)", cast.Actor, cast.Actress)
	if err != nil {
		return 0, err
	}

	//? Get the last inserted ID
	cast_id, err := cast_res.LastInsertId()
	if err != nil {
		return 0, err
	}

	//* Movies
	//? Execute statement
	movie_res, err := tx.Exec("INSERT INTO movies(title, rating, director_id, cast_id) VALUES (?, ?, ?, ?)", title, rating, director_id, cast_id)
	if err != nil {
		return 0, err
	}

	//? Get the last inserted ID
	movie_id, err := movie_res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil
	}

	return movie_id, nil
}
