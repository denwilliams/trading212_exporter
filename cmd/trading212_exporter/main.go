package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	openPositionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "trading212_open_positions_total",
		Help: "Total number of open positions in the Trading 212 account",
	})
	positionQtyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "trading212_position_quantity",
		Help: "Quantity of individual open positions",
	}, []string{"ticker"})
	positionValueGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "trading212_position_value",
		Help: "Value of individual open positions",
	}, []string{"ticker"})
	positionCostGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "trading212_position_cost",
		Help: "Cost of individual open positions",
	}, []string{"ticker"})
	positionPplGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "trading212_position_ppl",
		Help: "PPL of individual open positions",
	}, []string{"ticker"})
	positionFxPplGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "trading212_position_fxppl",
		Help: "FX PPL of individual open positions",
	}, []string{"ticker"})
)

// Position represents the structure for the given data.
type Position struct {
	AveragePrice    float64   `json:"averagePrice"`
	CurrentPrice    float64   `json:"currentPrice"`
	Frontend        string    `json:"frontend"`
	FxPpl           float64   `json:"fxPpl"`
	InitialFillDate time.Time `json:"initialFillDate"`
	MaxBuy          float64   `json:"maxBuy"`
	MaxSell         float64   `json:"maxSell"`
	PieQuantity     float64   `json:"pieQuantity"`
	Ppl             float64   `json:"ppl"`
	Quantity        float64   `json:"quantity"`
	Ticker          string    `json:"ticker"`
}

// FetchOpenPositions retrieves the list of open positions from the Trading 212 API
func FetchOpenPositions(apiKey string) ([]Position, error) {
	url := "https://live.trading212.com/api/v0/equity/portfolio"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiKey)

	// client := &http.Client{Timeout: 10 * time.Second}
	// resp, err := client.Do(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

// UpdateMetrics updates Prometheus metrics based on the fetched open positions
func UpdateMetrics(apiKey string) {
	positions, err := FetchOpenPositions(apiKey)
	if err != nil {
		log.Printf("Error fetching open positions: %v", err)
		return
	}

	// Update the total number of open positions
	openPositionsGauge.Set(float64(len(positions)))

	// Reset and update individual position values
	for _, pos := range positions {
		positionQtyGauge.WithLabelValues(pos.Ticker).Set(pos.Quantity)
		positionValueGauge.WithLabelValues(pos.Ticker).Set(pos.Quantity * pos.CurrentPrice)
		positionCostGauge.WithLabelValues(pos.Ticker).Set(pos.Quantity * pos.AveragePrice)
		positionPplGauge.WithLabelValues(pos.Ticker).Set(pos.Ppl)
		positionFxPplGauge.WithLabelValues(pos.Ticker).Set(pos.FxPpl)
	}

	log.Printf("Updated metrics for %d open positions", len(positions))
}

func main() {
	apiKey := os.Getenv("TRADING212_API_KEY")
	if apiKey == "" {
		log.Fatal("TRADING212_API_KEY environment variable is not set")
	}

	// Register Prometheus metrics
	prometheus.MustRegister(openPositionsGauge)
	prometheus.MustRegister(positionQtyGauge)
	prometheus.MustRegister(positionValueGauge)
	prometheus.MustRegister(positionCostGauge)
	prometheus.MustRegister(positionPplGauge)
	prometheus.MustRegister(positionFxPplGauge)

	// Start HTTP server for Prometheus to scrape
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting server on :9977")
		if err := http.ListenAndServe(":9977", nil); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// Periodically fetch data and update metrics
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	UpdateMetrics(apiKey)

	for range ticker.C {
		UpdateMetrics(apiKey)
	}
}
