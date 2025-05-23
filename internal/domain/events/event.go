package events

import (
	"time"

	"go-stock/internal/domain/models"
)

// EventType represents the type of a domain event.
type EventType string

const (
	// StockPriceChanged is triggered when a stock price changes.
	StockPriceChanged EventType = "STOCK_PRICE_CHANGED"

	// StockAdded is triggered when a new stock is added.
	StockAdded EventType = "STOCK_ADDED"

	// StockRemoved is triggered when a stock is removed.
	StockRemoved EventType = "STOCK_REMOVED"

	// AlertTriggered is triggered when a stock alert condition is met.
	AlertTriggered EventType = "ALERT_TRIGGERED"

	// SettingsChanged is triggered when user settings are updated.
	SettingsChanged EventType = "SETTINGS_CHANGED"

	// DataRefreshed is triggered when stock data is refreshed.
	DataRefreshed EventType = "DATA_REFRESHED"

	// MarketOpened is triggered when the market opens.
	MarketOpened EventType = "MARKET_OPENED"

	// MarketClosed is triggered when the market closes.
	MarketClosed EventType = "MARKET_CLOSED"
)

// Event is the base interface for all domain events.
type Event interface {
	// Type returns the type of the event.
	Type() EventType

	// Timestamp returns when the event occurred.
	Timestamp() time.Time

	// Payload returns the event data.
	Payload() interface{}
}

// BaseEvent provides common functionality for all events.
type BaseEvent struct {
	EventType  EventType
	OccurredAt time.Time
	Data       interface{}
}

// Type returns the type of the event.
func (e BaseEvent) Type() EventType {
	return e.EventType
}

// Timestamp returns when the event occurred.
func (e BaseEvent) Timestamp() time.Time {
	return e.OccurredAt
}

// Payload returns the event data.
func (e BaseEvent) Payload() interface{} {
	return e.Data
}

// NewEvent creates a new event with the given type and data.
func NewEvent(eventType EventType, data interface{}) Event {
	return BaseEvent{
		EventType:  eventType,
		OccurredAt: time.Now(),
		Data:       data,
	}
}

// StockPriceChangedEvent is triggered when a stock price changes.
type StockPriceChangedEvent struct {
	BaseEvent
	Stock    *models.StockInfo
	OldPrice float64
	NewPrice float64
}

// NewStockPriceChangedEvent creates a new stock price changed event.
func NewStockPriceChangedEvent(stock *models.StockInfo, oldPrice, newPrice float64) Event {
	event := StockPriceChangedEvent{
		Stock:    stock,
		OldPrice: oldPrice,
		NewPrice: newPrice,
	}
	event.EventType = StockPriceChanged
	event.OccurredAt = time.Now()
	event.Data = event // Self-reference for payload

	return event
}

// AlertTriggeredEvent is triggered when a stock alert condition is met.
type AlertTriggeredEvent struct {
	BaseEvent
	Alert *models.StockAlert
	Stock *models.StockInfo
}

// NewAlertTriggeredEvent creates a new alert triggered event.
func NewAlertTriggeredEvent(alert *models.StockAlert, stock *models.StockInfo) Event {
	event := AlertTriggeredEvent{
		Alert: alert,
		Stock: stock,
	}
	event.EventType = AlertTriggered
	event.OccurredAt = time.Now()
	event.Data = event // Self-reference for payload

	return event
}

// SettingsChangedEvent is triggered when user settings are updated.
type SettingsChangedEvent struct {
	BaseEvent
	OldSettings *models.Settings
	NewSettings *models.Settings
}

// NewSettingsChangedEvent creates a new settings changed event.
func NewSettingsChangedEvent(oldSettings, newSettings *models.Settings) Event {
	event := SettingsChangedEvent{
		OldSettings: oldSettings,
		NewSettings: newSettings,
	}
	event.EventType = SettingsChanged
	event.OccurredAt = time.Now()
	event.Data = event // Self-reference for payload

	return event
}

// MarketEvent represents market opening or closing events.
type MarketEvent struct {
	BaseEvent
	Exchange string
	IsOpen   bool
}

// NewMarketOpenedEvent creates a new market opened event.
func NewMarketOpenedEvent(exchange string) Event {
	event := MarketEvent{
		Exchange: exchange,
		IsOpen:   true,
	}
	event.EventType = MarketOpened
	event.OccurredAt = time.Now()
	event.Data = event // Self-reference for payload

	return event
}

// NewMarketClosedEvent creates a new market closed event.
func NewMarketClosedEvent(exchange string) Event {
	event := MarketEvent{
		Exchange: exchange,
		IsOpen:   false,
	}
	event.EventType = MarketClosed
	event.OccurredAt = time.Now()
	event.Data = event // Self-reference for payload

	return event
}
