package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blackopsrepl/go-rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// ENVIRONMENT //
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	dbQueries := database.New(db)

	apiConfig := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ROUTERS //
	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReady)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiConfig.handlerUsersCreate)
	v1Router.Get("/users", apiConfig.requireUserAuth(apiConfig.handlerUsersGet))

	v1Router.Post("/feeds", apiConfig.requireUserAuth(apiConfig.handlerFeedCreate))
	v1Router.Get("/feeds", apiConfig.handlerFeedsGet)

	v1Router.Post("/follows", apiConfig.requireUserAuth(apiConfig.handlerFollowCreate))
	v1Router.Get("/follows", apiConfig.requireUserAuth(apiConfig.handlerFollowsGet))
	v1Router.Delete("/follows/{FollowID}", apiConfig.requireUserAuth(apiConfig.handlerFollowDelete))

	v1Router.Get("/posts", apiConfig.requireUserAuth(apiConfig.handlerPostsGet))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// SCRAPING //
	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	// START SERVER //
	log.Printf("Server starting on port %v", portString)
	log.Fatal(srv.ListenAndServe())

}
