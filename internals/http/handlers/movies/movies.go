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
	"strconv"

	"github.com/go-playground/validator"
)

func New(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Root handler has been called")

		//? Decode JSON into Movie struct
		var movie types.Movie
		err := json.NewDecoder(r.Body).Decode(&movie)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
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
		id, err := db.CreateMovie(&movie)
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

func GetByID(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Get movie by ID handler called")

		//? Get ID string from pathvalue
		idStr := r.PathValue("id")
		if idStr == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing movie ID in URL")))
			logger.Error.Println("Missing movie ID in URL:")
			return
		}

		//? Parse idStr into int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid ID")))
			logger.Error.Println("Error parsing ID into int64:", err)
			return
		}

		//* Retrieve movie from database
		movie, err := db.GetMovieByID(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error retrieving movie:", err)
			return
		}

		if movie == nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("movie not found")))
			logger.Error.Println("Movie not found")
			return
		}

		//? Send response
		response.WriteJson(w, http.StatusOK, movie)
	}
}

func GetList(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Get movie list handler called")

		//? Get limit and offset from URL
		query := r.URL.Query()
		limitStr := query.Get("limit")
		offsetStr := query.Get("offset")

		//? Set default values if not provided
		if limitStr == "" {
			limitStr = "10"
		}
		if offsetStr == "" {
			offsetStr = "0"
		}

		//? Convert limit and offset to integers
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid limit value")))
			logger.Error.Println("Invalid limit value:", err)
			return
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid offset value")))
			logger.Error.Println("Invalid offset value:", err)
			return
		}

		//? Enforce a hard upper bound on limit
		const maxLimit int = 50
		if limit > maxLimit {
			limit = maxLimit
		}

		//* Retrieve movie list from database
		movies, err := db.GetMovieList(limit, offset)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error retrieving students:", err)
			return
		}

		//? Handle empty results gracefully
		if len(movies) == 0 {
			response.WriteJson(w, http.StatusOK, []types.Movie{})
			logger.Info.Println("Movie list is empty")
			return
		}

		//? Send response
		response.WriteJson(w, http.StatusOK, movies)
	}
}

func Update(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Update movie handler called")

		//? Get id from URL
		idStr := r.PathValue("id")
		if idStr == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing ID")))
			logger.Error.Println("Missing ID in URL")
			return
		}

		//? Parse idStr into int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid ID")))
			logger.Error.Println("Error parsing ID:", err)
			return
		}

		//? Check if movie exists
		m, err := db.GetMovieByID(id)
		if err != nil || m == nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("movie does not exist")))
			logger.Error.Println("Movie doen not exist", err)
			return
		}

		//? Decode JSON
		var movie types.Movie
		err = json.NewDecoder(r.Body).Decode(&movie)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
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

		//* Update movie
		updated_movie_id, err := db.UpdateMovie(id, &movie)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Failed to update movie:", err)
			return
		}

		if updated_movie_id == 0 {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("movie not found")))
			logger.Error.Println("Movie to update not found")
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{
			"success": "OK",
			"message": fmt.Sprintf("Movie updated with ID %d", updated_movie_id),
		})
	}
}

func DeleteByID(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Delete movie by ID handler called")

		//? Get id string from URL
		idStr := r.PathValue("id")
		if idStr == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing ID")))
			logger.Error.Println("Missing ID in URL")
			return
		}

		//? Parse idStr into int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid ID")))
			logger.Error.Println("Invalid ID:", err)
			return
		}

		//* Delete movie from database
		deleted_movie_id, err := db.DeleteMovieByID(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Failed to delete movie:", err)
		}

		if deleted_movie_id == 0 {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("movie not found")))
			logger.Error.Println("Movie to delete not found")
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{
			"success": "OK",
			"message": fmt.Sprintf("Movie deleted with ID %d", deleted_movie_id),
		})
	}
}
