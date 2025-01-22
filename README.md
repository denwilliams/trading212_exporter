# Trading 212 Prometheus Metrics Exporter

Barebones Prometheus exporter for Trading 212. It fetches your open positions from [the Trading 212 API](https://t212public-api-docs.redoc.ly/#tag/Personal-Portfolio) and exposes metrics in a format that can be scraped by Prometheus.

Current metrics:

```
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
```

Run it with your API key as an environment variable:

TRADING212_API_KEY=your_api_key_here go run main.go
