package models

import (
	"time"
)

// StockInfo represents basic information about a stock.
type StockInfo struct {
	ID         string    `json:"id"`         // Unique identifier
	Code       string    `json:"code"`       // Stock code/symbol (e.g., "AAPL")
	Name       string    `json:"name"`       // Stock name (e.g., "Apple Inc.")
	Exchange   string    `json:"exchange"`   // Exchange (e.g., "NASDAQ")
	Industry   string    `json:"industry"`   // Industry
	Price      float64   `json:"price"`      // Current price
	Change     float64   `json:"change"`     // Change in price
	ChangeRate float64   `json:"changeRate"` // Change rate in percentage
	Volume     int64     `json:"volume"`     // Trading volume
	UpdateTime time.Time `json:"updateTime"` // Last update time
}

// HistoricalPrice represents a historical price record of a stock.
type HistoricalPrice struct {
	ID        string    `json:"id"`        // Unique identifier
	StockID   string    `json:"stockId"`   // Reference to the stock
	Date      time.Time `json:"date"`      // Date of the price
	Open      float64   `json:"open"`      // Opening price
	Close     float64   `json:"close"`     // Closing price
	High      float64   `json:"high"`      // Highest price
	Low       float64   `json:"low"`       // Lowest price
	Volume    int64     `json:"volume"`    // Trading volume
	Turnover  float64   `json:"turnover"`  // Turnover
	UpdatedAt time.Time `json:"updatedAt"` // When this record was updated
}

// StockAlert represents an alert configuration for a stock price.
type StockAlert struct {
	ID            string    `json:"id"`            // Unique identifier
	StockID       string    `json:"stockId"`       // Reference to the stock
	AlertType     AlertType `json:"alertType"`     // Type of alert
	Threshold     float64   `json:"threshold"`     // Price threshold
	IsActive      bool      `json:"isActive"`      // Whether the alert is active
	Triggered     bool      `json:"triggered"`     // Whether the alert has been triggered
	LastTriggered time.Time `json:"lastTriggered"` // When the alert was last triggered
	CreatedAt     time.Time `json:"createdAt"`     // When the alert was created
}

// AlertType represents the type of a stock alert.
type AlertType string

const (
	// PriceAbove is triggered when the stock price goes above a threshold
	PriceAbove AlertType = "PRICE_ABOVE"
	// PriceBelow is triggered when the stock price goes below a threshold
	PriceBelow AlertType = "PRICE_BELOW"
	// ChangeRateAbove is triggered when the change rate goes above a threshold
	ChangeRateAbove AlertType = "CHANGE_RATE_ABOVE"
	// ChangeRateBelow is triggered when the change rate goes below a threshold
	ChangeRateBelow AlertType = "CHANGE_RATE_BELOW"
)

// FollowedStock represents a stock that the user is following.
type FollowedStock struct {
	ID         string    `json:"id"`         // Unique identifier
	StockID    string    `json:"stockId"`    // Reference to the stock
	GroupID    string    `json:"groupId"`    // Group this stock belongs to
	AddedAt    time.Time `json:"addedAt"`    // When the stock was added to the followed list
	Note       string    `json:"note"`       // User notes
	IsWatching bool      `json:"isWatching"` // Whether the user is actively watching this stock
}

// StockGroup represents a group of stocks for organization.
type StockGroup struct {
	ID          string    `json:"id"`          // Unique identifier
	Name        string    `json:"name"`        // Group name
	Description string    `json:"description"` // Group description
	CreatedAt   time.Time `json:"createdAt"`   // When the group was created
}
