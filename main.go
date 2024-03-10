package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Create MinIO operations instance
    minioOps, err := NewMinIOOperations(
        os.Getenv("MINIO_ENDPOINT"),
        os.Getenv("MINIO_ACCESS_KEY"),
        os.Getenv("MINIO_SECRET_KEY"),
        true,
    )
    if err != nil {
        log.Fatalf("Failed to create MinIO operations: %v", err)
    }

    // Create Docker operations instance
    dockerOps, err := NewDockerOperations()
    if err != nil {
        log.Fatalf("Failed to create Docker operations: %v", err)
    }

    // Create HTTP server
    srv := &http.Server{
        Addr: ":8080",
        Handler: &Handler{
            MinIOOps:  minioOps,
            DockerOps: dockerOps,
        },
    }

    log.Println("Server started on port 8080")
    log.Fatal(srv.ListenAndServe())
}