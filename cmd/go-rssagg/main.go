package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blackopsrepl/go-rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {

	// ENVIRONMENT //

	// parse flags and load anvironment
	envFile := flag.String("env", "", "Path to .env file")
	flag.Parse()

	setEnv(*envFile)

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not found in the environment")
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		log.Fatal("OLLAMA_URL is not found in the environment")
	}

	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		log.Fatal("OLLAMA_MODEL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	dbQueries := database.New(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	apiConfig := apiConfig{
		DB:  dbQueries,
		RDB: rdb,
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
	v1Router.Get("/posts/{PostDate}", apiConfig.requireUserAuth(apiConfig.handlerPostsGetByDate))
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// SCRAPING //
	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	go startCrawling(dbQueries, rdb, collectionConcurrency, collectionInterval)

	// START SERVER //
	log.Printf("Server starting on port %v", portString)
	log.Fatal(srv.ListenAndServe())

}

func setEnv(envFile string) {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	} else if os.Getenv("DB_URL") == "" || os.Getenv("PORT") == "" || os.Getenv("OLLAMA_URL") == "" || os.Getenv("OLLAMA_MODEL") == "" || os.Getenv("REDIS_URL") == "" {
		log.Fatalf("DB_URL, PORT, OLLAMA_URL, OLLAMA_MODEL, and REDIS_URL must be set as environment variables!\n\nLoaded DB_URL: %s, PORT: %s, OLLAMA_URL: %s, OLLAMA_MODEL: %s, REDIS_URL: %s", os.Getenv("DB_URL"), os.Getenv("PORT"), os.Getenv("OLLAMA_URL"), os.Getenv("OLLAMA_MODEL"), os.Getenv("REDIS_URL"))
	}
}
