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
		UNIQUE(title, director_id, cast_id),
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
	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	//? ----------- DIRECTOR: check if exists -----------
	var director_id int64
	dir_row := tx.QueryRow("SELECT id FROM directors WHERE name = ?", director.Name)
	err = dir_row.Scan(&director_id)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err := tx.Exec("INSERT INTO directors(name, age) VALUES (?, ?)", director.Name, director.Age)
			if err != nil {
				return 0, err
			}
			director_id, err = res.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	//? ----------- Cast: check if exists -----------
	var cast_id int64
	cast_row := tx.QueryRow("SELECT id FROM casts WHERE actor = ? AND actress = ?", cast.Actor, cast.Actress)
	err = cast_row.Scan(&cast_id)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err := tx.Exec("INSERT INTO casts(actor, actress) VALUES (?, ?)", cast.Actor, cast.Actress)
			if err != nil {
				return 0, err
			}
			cast_id, err = res.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	//? ----------- Movie: check if exists -----------
	var movie_id int64
	movie_row := tx.QueryRow("SELECT id FROM movies WHERE title = ?", title)
	err = movie_row.Scan(&movie_id)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err := tx.Exec("INSERT INTO movies(title, rating, director_id, cast_id) VALUES (?, ?, ?, ?)", title, rating, director_id, cast_id)
			if err != nil {
				return 0, err
			}
			movie_id, err = res.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			return 0, nil
		}

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return movie_id, nil
}

func (s *SQLite) GetMovieByID(id int64) (*types.Movie, error) {
	row := s.DB.QueryRow(`
		SELECT
			m.id, m.title, m.rating,
			d.id, d.name, d.age,
			c.id, c.actor, c.actress
		FROM movies m
		LEFT JOIN directors d ON m.director_id = d.id
		LEFT JOIN casts c ON m.cast_id = c.id
		WHERE m.id = ?
	`, id)

	var movie types.Movie
	movie.Director = &types.Director{}
	movie.Cast = &types.Cast{}

	err := row.Scan(
		&movie.ID, &movie.Title, &movie.Rating,
		&movie.Director.ID, &movie.Director.Name, &movie.Director.Age,
		&movie.Cast.ID, &movie.Cast.Actor, &movie.Cast.Actress,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, nil
	}

	return &movie, nil
}

func (s *SQLite) GetMovieList(limit int, offset int) ([]*types.Movie, error) {
	rows, err := s.DB.Query(`
		SELECT
			m.id, m.title, m.rating,
			d.id, d.name, d.age,
			c.id, c.actor, c.actress
		FROM movies m
		LEFT JOIN directors d ON m.director_id = d.id
		LEFT JOIN casts c ON m.cast_id = c.id
		ORDER BY m.id
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*types.Movie

	for rows.Next() {
		movie := &types.Movie{
			Director: &types.Director{},
			Cast:     &types.Cast{},
		}

		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Rating,
			&movie.Director.ID, &movie.Director.Name, &movie.Director.Age,
			&movie.Cast.ID, &movie.Cast.Actor, &movie.Cast.Actress,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
