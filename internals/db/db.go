package db

import "github/MahfujulSagor/movies_crud/internals/types"

type DB interface {
	CreateMovie(movie *types.Movie) (int64, error)
	GetMovieByID(id int64) (*types.Movie, error)
	GetMovieList(limit int, offset int) ([]*types.Movie, error)
	UpdateMovie(movie *types.Movie) (int64, error)
	DeleteMovieByID(id int64) (int64, error)
}
