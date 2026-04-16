package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting Reconciliation Service...")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Running reconciliation routine at", t)
				RunReconciliation()
				GenerateAuditReport()
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Reconciliation Service...")
	done <- true
}

func RunReconciliation() {
	// Compare data from discovery with the source of truth in Neo4j
	log.Println("Reconciling discovered data against Graph DB inventory...")
}

func GenerateAuditReport() {
	log.Println("Generating reconciliation audit report...")
}
