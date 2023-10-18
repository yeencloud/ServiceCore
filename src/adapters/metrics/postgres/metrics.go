package postgres

import "gorm.io/gorm"

type MetricHTTPRequest struct {
	gorm.Model

	Service string
	Method  string
	Code    int
	URL     string

	ElapsedTime float64
}

func (db *Database) PushNewMetric(metric MetricHTTPRequest) {
	db.engine.Create(&metric)
}
