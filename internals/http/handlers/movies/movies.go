package movies

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/MahfujulSagor/movies_crud/internals/db"
	"github/MahfujulSagor/movies_crud/internals/logger"
	"github/MahfujulSagor/movies_crud/internals/types"
	"github/MahfujulSagor/movies_crud/internals/utils/response"
	"io"
	"net/http"

	"github.com/go-playground/validator"
)

func New(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Root handler has been called")

		//? Decode JSON into Movie struct
		var movie types.Movie
		err := json.NewDecoder(r.Body).Decode(&movie)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			logger.Error.Println("Empty body:", err)
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			logger.Error.Println("Error decoding movie:", err)
			return
		}
		defer r.Body.Close()

		//? Request validation
		if err := validator.New().Struct(movie); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			logger.Error.Println("Validation error:", err)
			return
		}

		//* Create movie in database
		id, err := db.CreateMovie(movie.Title, movie.Rating, movie.Director, movie.Cast)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Failed to create movie:", err)
			return
		}

		logger.Info.Println("Movie created with ID:", id)

		//? Send response
		response.WriteJson(w, http.StatusCreated, map[string]string{
			"success": "OK",
			"message": fmt.Sprintf("Movie created with ID: %d", id),
		})
	}
}
