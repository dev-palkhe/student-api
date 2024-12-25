package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-palkhe/student-api/internal/config"
	"github.com/dev-palkhe/student-api/internal/handlers"
	"github.com/dev-palkhe/student-api/internal/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	instanceName := os.Getenv("INSTANCE_NAME")
	if instanceName == "" {
		instanceName = "default" // Set a default name if not provided
	}

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	studentRepo := repository.NewStudentRepository(db)
	studentHandler := handlers.NewStudentHandler(studentRepo)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/healthcheck", healthcheck)
		v1.POST("/students", studentHandler.CreateStudent)
		v1.GET("/students", studentHandler.GetAllStudents)
		v1.GET("/students/:id", studentHandler.GetStudentByID)
		v1.PUT("/students/:id", studentHandler.UpdateStudent)
		v1.DELETE("/students/:id", studentHandler.DeleteStudent)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
