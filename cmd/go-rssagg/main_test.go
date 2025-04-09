package main

import (
	"os"
	"testing"
)

// func setEnv() {
// 	envFile := flag.String("env", ".env", "Path to .env file")

// 	flag.Parse()

// 	if *envFile != "" {
// 		err := godotenv.Load(*envFile)
// 		if err != nil {
// 			log.Fatalf("Error loading .env file: %v", err)
// 		}
// 	} else if os.Getenv("DB_URL") == "" || os.Getenv("PORT") == "" {
// 		log.Fatal("DB_URL and PORT must be set as environment variables!")
// 	}
// }

func TestSetEnv(t *testing.T) {

	testCases := []struct {
		name       string
		envFile    string
		preloadEnv bool
		wantErr    bool
	}{
		{
			name:       "No flag",
			envFile:    "",
			preloadEnv: true,
			wantErr:    false,
		},
		{
			name:       "No flag and no environment variables",
			envFile:    "",
			preloadEnv: false,
			wantErr:    true,
		},
		{
			name:       "Valid env file",
			envFile:    "../../.env.template",
			preloadEnv: false,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.preloadEnv {
				os.Setenv("DB_URL", "db_url")
				os.Setenv("PORT", "port")
			}

			// Call the function
			setEnv(tc.envFile)

			// Check the result
			if tc.envFile != "" {
				if os.Getenv("DB_URL") == "" || os.Getenv("PORT") == "" {
					t.Errorf("Expected DB_URL and PORT to be set, but they were not")
				}
			} else if tc.envFile == "" && os.Getenv("DB_URL") == "" && os.Getenv("PORT") == "" {
				t.Errorf("Expected DB_URL and PORT to be empty, but they were not")
			}
		})
	}

}
