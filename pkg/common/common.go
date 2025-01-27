package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DefaultHeartbeat = 5 * time.Second
)

func GetDBConnectionString() string {

	// set the default values here
	requiredEnvVars := map[string]string{
		"POSTGRES_USER":     "",
		"POSTGRES_PASSWORD": "",
		"POSTGRES_DB":       "",
		"POSTGRES_HOST":     "localhost",
	}

	var missingEnvVars []string
	// set the value from the environement here
	for envVarName, _ := range requiredEnvVars {
		envVarValue := os.Getenv(envVarName)
		if envVarValue == "" {
			missingEnvVars = append(missingEnvVars, envVarName)
		}
		requiredEnvVars[envVarName] = envVarValue
	}

	if len(missingEnvVars) > 0 {
		log.Fatalf("The following required environment variables are not set: %s",
			strings.Join(missingEnvVars, ", "))
	}

	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		requiredEnvVars["POSTGRES_USER"],
		requiredEnvVars["POSTGRES_PASSWORD"],
		requiredEnvVars["POSTGRES_HOST"],
		requiredEnvVars["POSTGRES_DB"])

}

func ConnectToDatabase(ctx context.Context, dbConnectionString string) (*pgxpool.Pool, error) {
	var dbPool *pgxpool.Pool
	var err error
	retryCount := 0
	for retryCount < 5 {
		dbPool, err = pgxpool.Connect(ctx, dbConnectionString)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to the database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		retryCount++
	}

	if err != nil {
		log.Printf("Ran out of retries to connect to database (5)")
		return nil, err
	}

	log.Printf("Connected to the database.")
	return dbPool, nil
}
