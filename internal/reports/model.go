// Package reports produce reportes agregados de ventas (solo lectura), acotados
// por negocio y rango de fechas.
package reports

type PaymentBreakdown struct {
	Method     string `json:"method"`
	Count      int    `json:"count"`
	TotalCents int    `json:"totalCents"`
}

type CategoryBreakdown struct {
	CategoryName string `json:"categoryName"`
	Quantity     int    `json:"quantity"`
	TotalCents   int    `json:"totalCents"`
}

type HourBreakdown struct {
	Hour       int `json:"hour"`
	Count      int `json:"count"`
	TotalCents int `json:"totalCents"`
}

type SalesReport struct {
	TotalCents      int                 `json:"totalCents"`
	SalesCount      int                 `json:"salesCount"`
	ByPaymentMethod []PaymentBreakdown  `json:"byPaymentMethod"`
	ByCategory      []CategoryBreakdown `json:"byCategory"`
	ByHour          []HourBreakdown     `json:"byHour"`
}
