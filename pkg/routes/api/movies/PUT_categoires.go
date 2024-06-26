package movies

import (
	"database/sql"
	"log"
	"net/http"
	"project/internal/utils"

	"github.com/gin-gonic/gin"
)

// @Tags			movies
// @Summary		Update movie category
// @Description	Updates the category of a movie with the specified ID
// @Produce		json
// @Param			id	path	string	true	"Movie ID"
// @Security		ApiKeyAuth
// @Param			categoryName	body		string							true	"Category Name"
// @Success		200				{object}	routes.DefaultMessageResponse	"Movie category updated"
// @Failure		400				{object}	routes.DefaultMessageResponse	"Bad request"
// @Failure		401				{object}	routes.DefaultMessageResponse	"Unauthorized"
// @Failure		404				{object}	routes.DefaultMessageResponse	"Category not found"
// @Failure		404				{object}	routes.DefaultMessageResponse	"Movie not found"
// @Failure		500				{object}	routes.DefaultMessageResponse	"Internal server error"
// @Router			/api/movies/category/{id} [put]
func (m *MoviesRoute) PUT_MovieCategory(c *gin.Context) {
	movieId := c.Param("id")
	userRole := c.GetString("role")
	userId := c.GetInt("userId")

	if userRole != "admin" || userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	valid, movieIdNum := utils.IsValidNum(movieId)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Bad request",
		})
		return
	}
	data := struct {
		CategoryName string `json:"categoryName" binding:"required"`
	}{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}
	category, err := m.DB.CategoriesRepository.GetCategoryByName(data.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{
				"message": "Category not found",
			})
			return
		}
		log.Println(err.Error())
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}

	err = m.DB.MovieRepository.UpdateMovieCategory(movieIdNum, category.ID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(404, gin.H{
			"message": "Movie not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Movie category updated",
	})

}
