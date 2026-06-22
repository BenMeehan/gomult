package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benmeehan/gomult/internal/ai"
	"github.com/benmeehan/gomult/internal/config"
	"github.com/benmeehan/gomult/internal/sandbox"
	"github.com/benmeehan/gomult/internal/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	aiKey := flag.String("ai-key", "", "Deepseek API key for AI analysis (or set DEEPSEEK_API_KEY env var)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("gomult starting")

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	engine, err := sandbox.NewEngine(cfg)
	if err != nil {
		log.Fatalf("failed to create execution engine: %v", err)
	}

	key := *aiKey
	if key == "" {
		key = os.Getenv("DEEPSEEK_API_KEY")
	}
	var aiClient *ai.Client
	if key != "" {
		aiClient = ai.NewClient(key)
		log.Print("AI analysis enabled (deepseek)")
	}

	srv := server.New(cfg, engine, aiClient)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Print("shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("listening on :%d", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
