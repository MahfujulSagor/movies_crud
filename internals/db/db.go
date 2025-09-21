package db

import "github/MahfujulSagor/movies_crud/internals/types"

type DB interface {
	CreateMovie(title string, rating int, director *types.Director, cast *types.Cast) (int64, error)
	GetMovieByID(id int64) (*types.Movie, error)
}
