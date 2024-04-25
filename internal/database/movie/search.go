package movie

import (
	"strings"
)

func (db *MovieRepository) SearchMovie(query string) ([]Movie, error) {
	query = strings.ToLower(query)

	sqlQuery := `
		SELECT * FROM movies
		WHERE LOWER(name) LIKE '%' || $1 || '%'
		OR LOWER(keywords) LIKE '%' || $1 || '%'
	`

	rows, err := db.Database.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.Id,
			&movie.UserId,
			&movie.ImageUrl,
			&movie.Name,
			&movie.Year,
			&movie.CategoryId,
			&movie.AgeCategoryId,
			&movie.WatchCount,
			&movie.DurationMinutes,
			&movie.Keywords,
			&movie.Description,
			&movie.Director,
			&movie.Producer,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
