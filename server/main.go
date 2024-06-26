package main

import (
	"log" ////////////////////// SWAGGER
	_ "project/docs"
	"project/internal/database/datasets"
	"project/pkg/middleware"
	"project/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title                    ozinshe-api-test
// @description              This is a test-backend-RESTful-API for managing movies, seasons, episodes and etc.
// @contact.name             Adilhan
// @contact.email            bigxxby@yandex.ru
// @host                     localhost:8080
// @BasePath                 /
// @schemes                  http
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err.Error())
		return
	}
	//init main
	main, err := routes.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// test data
	err = datasets.InitDatasets(main.DB)
	if err != nil {
		log.Println(err.Error())
		return
	}
	router := gin.Default()

	// add swagger
	router.GET("/docs/*any", swagger.WrapHandler(swaggerFiles.Handler))

	// API
	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware())

	{

		// auth
		auth := apiRoutes.Group("/")
		{
			auth.POST("/check-auth", main.AuthRoute.POST_CheckAuth) //returns CURRENT users role
			auth.POST("/signUp", main.AuthRoute.POST_SignUp)        //{email, password, role } required bindings
			auth.POST("/signIn", main.AuthRoute.POST_SignIn)        //{email, password } required bindings

			//OLD

			//htmlRoutes.GET("/send-code", main.GET_HTML_SendRestoreCode) //-------->
			// auth.POST("/send-code", main.AuthRoute.POST_SendCode) // sends restore code to the users email
			// auth.POST("/verify", main.AuthRoute.POST_VerifyCode)  // gets code from user an verifies it, if valid sends to the email link for resetting the pass
			// htmlRoutes.GET("/change-password", main.GET_ChangePassword) //------->
			// auth.POST("/reset-password", main.AuthRoute.POST_ResetPassword) // changes password of the user

			//NEW
			auth.POST("/send-code", main.AuthRoute.POST_SendCode)             // sends verifiication code to the email
			auth.POST("/verify-code", main.AuthRoute.POST_VerifyCode)         // verifies code and changes users pass to temporary
			auth.POST("/change-password", main.AuthRoute.POST_ChangePassword) // changes password of the user

		}
		// movies
		movies := apiRoutes.Group("/movies")
		{
			movies.GET("/", main.MoviesRoute.GET_Movies)                                  //?limit=<limitInt>
			movies.GET("/:id", main.MoviesRoute.GET_Movie)                                // returns movie by id of movie
			movies.GET("/by-genre/", main.MoviesRoute.GET_EveryMovieByGenre)              // get all movies for every genre
			movies.GET("/by-category/", main.MoviesRoute.GET_EveryMovieByCategory)        // get all movies for every category
			movies.GET("/by-age-category/", main.MoviesRoute.GET_EveryMovieByAgeCategory) // get all movies for every age category
			movies.POST("/", main.MoviesRoute.POST_Movie)                                 // admin only create movie
			movies.DELETE("/:id", main.MoviesRoute.DELETE_Movie)                          // admin only deletes movie
			movies.GET("/search", main.MoviesRoute.GET_Search)                            // ?query=<searchQuery>
			movies.POST("watch/:id", main.MoviesRoute.POST_Watch)                         // +1 movie watch count ONLY if user authenticated
			movies.PUT("category/:id", main.MoviesRoute.PUT_MovieCategory)                // admin only change category
			movies.PUT("/data/:id", main.MoviesRoute.PUT_MovieData)                       // admin only change data of movie (not related to other tables)
			movies.PUT("/age-category/:id", main.MoviesRoute.PUT_MovieAgeCategory)        // admin only change age category of the movie
			movies.PUT("/genres/:id", main.MoviesRoute.PUT_MovieGenres)                   // admin only change genres of the movie
		}
		//seasons
		seasons := apiRoutes.Group("/seasons")
		{
			seasons.GET("/:id", main.SeasonsRoute.GET_Season)                                 // returns season by id
			seasons.GET("/movie/:id", main.SeasonsRoute.GET_AllSeasonsOfMovie)                // returns all seasons by movieId
			seasons.POST("/:id", main.SeasonsRoute.POST_CreateSeason)                         // adds season to the movie Id (admin)
			seasons.POST("/:id/multiple", main.SeasonsRoute.POST_CreateSeasons)               // adds season to the movie Id (admin)
			seasons.PUT("/:id", main.SeasonsRoute.PUT_Season)                                 // updates season by its id
			seasons.DELETE("/:id", main.SeasonsRoute.DELETE_Season)                           // deletes season by id  (admin)
			seasons.DELETE("/movie/:id/:seasonNumber", main.SeasonsRoute.DELETE_SeasonNumber) // deletes *ALL season numers of current movie id  (admin)
			seasons.DELETE("/movie/:id/clear", main.SeasonsRoute.DELETE_AllSeasonsOfMovie)    // deletes all seasons of current movie id  (admin)
		}
		//episodes (order not structured can add same episodes over and over)
		episodes := apiRoutes.Group("/episodes")
		{
			episodes.GET("/:id", main.EpisodesRoute.GET_Episode)                                            // returns all episodes by seasonId
			episodes.POST("/:id", main.EpisodesRoute.POST_Episode)                                          // adds episode to the season
			episodes.POST("/:id/multiple", main.EpisodesRoute.POST_Episodes)                                // adds multiple episodes to the season
			episodes.PUT("/:id", main.EpisodesRoute.PUT_Episode)                                            // updates episode by its ids
			episodes.DELETE("/:id", main.EpisodesRoute.DELETE_Episode)                                      // deletes episode by its id
			episodes.DELETE("/season/:id/:episodeNumber", main.EpisodesRoute.DELETE_EpisodeOfCurrentSeason) // deletes *ALL episodes by episode_number of selected season
			episodes.DELETE("/season/:id/clear", main.EpisodesRoute.DELETE_AllEpisodesOfSeason)             // deletes all episodes of selected season
		}

		// profile
		profile := apiRoutes.Group("/profile")
		{
			profile.GET("/", main.UsersRoute.GET_Profile) //get profile of current user
			profile.PUT("/", main.UsersRoute.PUT_Profile) //update profile of current user (dob , name , phone)
		}

		// favorites
		favorites := apiRoutes.Group("/favorites")
		{
			favorites.GET("/", main.FavoritesRoute.GET_Favorites)             //returns fav movies of CURRENT USER
			favorites.POST("/:id", main.FavoritesRoute.POST_Favorite)         //returns fav movie by id of CURRENT USER
			favorites.DELETE("/:id", main.FavoritesRoute.DELETE_Favorite)     //delete fav movie of CURRENT USER
			favorites.DELETE("/clear/", main.FavoritesRoute.DELETE_Favorites) //delete all fav movies of CURRENT USER
		}
		// categories
		categories := apiRoutes.Group("/categories")
		{
			categories.GET("/", main.CategoriesRoute.GET_Categories)        // return all categories
			categories.GET("/:id", main.CategoriesRoute.GET_Category)       // returns category by id
			categories.POST("/", main.CategoriesRoute.POST_Category)        //admin only
			categories.PUT("/:id", main.CategoriesRoute.PUT_Category)       //admin only
			categories.DELETE("/:id", main.CategoriesRoute.DELETE_Category) //admin only
		}
		// genres
		genres := apiRoutes.Group("/genres")
		{
			genres.GET("/", main.GenreRoute.GET_Genres)         //return all genres
			genres.GET("/:id", main.GenreRoute.GET_Genre)       //return genre by id
			genres.POST("/", main.GenreRoute.POST_Genre)        ////admin
			genres.DELETE("/:id", main.GenreRoute.DELETE_Genre) ////admin
			genres.PUT("/:id", main.GenreRoute.PUT_Genre)       ////admin

		}
		//age
		age := apiRoutes.Group("/age-categories")
		{
			age.GET("/", main.AgeRoute.GET_AgeCategories)        // return all age categories
			age.GET("/:id", main.AgeRoute.GET_AgeCategory)       //return category by id
			age.POST("/", main.AgeRoute.POST_AgeCategory)        //admin only
			age.DELETE("/:id", main.AgeRoute.DELETE_AgeCategory) //admin only
			age.PUT("/:id", main.AgeRoute.PUT_AgeCategory)       //admin only
		}
		//posters
		posters := apiRoutes.Group("/posters")
		{
			posters.GET("/:id", main.PosterRoute.GET_PostersOfMoive)           //get posters of movie
			posters.POST("/:id", main.PosterRoute.POST_PostersOfMoive)         //create posters of movie (admin)
			posters.DELETE("/:id", main.PosterRoute.DELETE_Posters)            //delete posters (admin )
			posters.DELETE("/movie/:id", main.PosterRoute.DELETE_PostersMovie) //delete posters of movie id (admin )
		}

	}
	router.Run(":8080")
}
