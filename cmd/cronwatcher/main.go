package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/alert"
	"github.com/cronwatcher/cronwatcher/internal/api"
	"github.com/cronwatcher/cronwatcher/internal/config"
	"github.com/cronwatcher/cronwatcher/internal/notify"
	"github.com/cronwatcher/cronwatcher/internal/scheduler"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build alerter chain
	alerter := buildAlerter(cfg)

	// Initialise watcher with registered jobs
	w := watcher.New(cfg)

	// Set up scheduler for cron-based miss detection
	sched := scheduler.New(cfg, w)
	for _, job := range cfg.Jobs {
		if err := sched.Register(job); err != nil {
			log.Fatalf("failed to register job %q: %v", job.Name, err)
		}
	}

	// Root context — cancelled on OS signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start notification runner (missed + long-running checks)
	notifier := notify.New(cfg, w, alerter)
	runner := notify.NewRunner(cfg, notifier)
	go runner.Run(ctx)

	// Start history pruner
	pruner := notify.NewPruner(cfg, w)
	go pruner.Run(ctx)

	// Start HTTP API server
	router := http.NewServeMux()
	srv := api.New(cfg, w, router)
	api.RegisterRoutes(router, cfg, w)

	httpServer := &http.Server{
		Addr:         cfg.API.ListenAddr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API server listening on %s", cfg.API.ListenAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down…")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("cronwatcher stopped")
}

// buildAlerter constructs the alerter chain from configuration.
func buildAlerter(cfg *config.Config) alert.Alerter {
	var alerters []alert.Alerter

	// Always include the structured logger alerter
	alerters = append(alerters, alert.NewLogger(log.Default()))

	if cfg.Alerts.Slack.WebhookURL != "" {
		alerters = append(alerters, alert.NewSlack(cfg.Alerts.Slack.WebhookURL))
	}

	if cfg.Alerts.PagerDuty.RoutingKey != "" {
		alerters = append(alerters, alert.NewPagerDuty(cfg.Alerts.PagerDuty.RoutingKey))
	}

	if cfg.Alerts.Webhook.URL != "" {
		alerters = append(alerters, alert.NewWebhook(cfg.Alerts.Webhook.URL))
	}

	return alert.NewMulti(alerters...)
}
