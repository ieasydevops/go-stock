package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go-stock/internal/domain/events"
	"go-stock/internal/domain/models"
	"go-stock/internal/domain/repositories"
)

// StockMonitorService provides stock monitoring and alerting functionality.
type StockMonitorService struct {
	stockRepo       repositories.StockRepository
	settingsRepo    repositories.SettingsRepository
	eventDispatcher events.EventDispatcher
	monitorTicker   *time.Ticker
	stopChan        chan struct{}
	running         bool
	mutex           sync.RWMutex
	notifyFunc      func(alert *models.StockAlert, stock *models.StockInfo) // Function to call for notifications
}

// NewStockMonitorService creates a new stock monitor service.
func NewStockMonitorService(
	stockRepo repositories.StockRepository,
	settingsRepo repositories.SettingsRepository,
	eventDispatcher events.EventDispatcher,
	notifyFunc func(alert *models.StockAlert, stock *models.StockInfo),
) *StockMonitorService {
	return &StockMonitorService{
		stockRepo:       stockRepo,
		settingsRepo:    settingsRepo,
		eventDispatcher: eventDispatcher,
		stopChan:        make(chan struct{}),
		notifyFunc:      notifyFunc,
	}
}

// Start begins monitoring stocks for alerts.
func (s *StockMonitorService) Start(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return errors.New("monitor is already running")
	}

	// Get refresh interval from settings
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Set default interval if not configured
	interval := time.Duration(settings.AutoRefreshInterval) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second // Default to 60 seconds
	}

	s.monitorTicker = time.NewTicker(interval)
	s.running = true

	// Start the monitoring loop in a separate goroutine
	go s.monitorLoop(ctx)

	return nil
}

// Stop stops the monitoring loop.
func (s *StockMonitorService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	s.monitorTicker.Stop()
	close(s.stopChan)
	s.running = false
}

// monitorLoop is the main monitoring loop.
func (s *StockMonitorService) monitorLoop(ctx context.Context) {
	// Check alerts immediately on start
	if err := s.checkAlerts(ctx); err != nil {
		log.Printf("Error checking alerts: %v", err)
	}

	for {
		select {
		case <-s.monitorTicker.C:
			if err := s.checkAlerts(ctx); err != nil {
				log.Printf("Error checking alerts: %v", err)
			}
		case <-s.stopChan:
			return
		case <-ctx.Done():
			s.Stop()
			return
		}
	}
}

// checkAlerts checks all active alerts for trigger conditions.
func (s *StockMonitorService) checkAlerts(ctx context.Context) error {
	// Get all active alerts
	alerts, err := s.stockRepo.GetActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active alerts: %w", err)
	}

	// Process each alert
	for _, alert := range alerts {
		stock, err := s.stockRepo.GetByID(ctx, alert.StockID)
		if err != nil {
			log.Printf("Failed to get stock for alert %s: %v", alert.ID, err)
			continue
		}

		if stock == nil {
			log.Printf("Stock not found for alert %s", alert.ID)
			continue
		}

		// Check if alert should be triggered
		triggered := s.isAlertTriggered(alert, stock)

		if triggered && !alert.Triggered {
			// Alert has been triggered
			alert.Triggered = true
			alert.LastTriggered = time.Now()

			// Update alert status
			if err := s.stockRepo.UpdateAlertStatus(ctx, alert.ID, alert.IsActive, alert.Triggered, alert.LastTriggered); err != nil {
				log.Printf("Failed to update alert status: %v", err)
				continue
			}

			// Dispatch alert triggered event
			event := events.NewAlertTriggeredEvent(alert, stock)
			s.eventDispatcher.Dispatch(event)

			// Call notification function if provided
			if s.notifyFunc != nil {
				s.notifyFunc(alert, stock)
			}
		} else if !triggered && alert.Triggered {
			// Alert is no longer triggered, reset it
			alert.Triggered = false

			// Update alert status
			if err := s.stockRepo.UpdateAlertStatus(ctx, alert.ID, alert.IsActive, alert.Triggered, alert.LastTriggered); err != nil {
				log.Printf("Failed to update alert status: %v", err)
			}
		}
	}

	return nil
}

// isAlertTriggered checks if an alert should be triggered based on current stock data.
func (s *StockMonitorService) isAlertTriggered(alert *models.StockAlert, stock *models.StockInfo) bool {
	switch alert.AlertType {
	case models.PriceAbove:
		return stock.Price > alert.Threshold
	case models.PriceBelow:
		return stock.Price < alert.Threshold
	case models.ChangeRateAbove:
		return stock.ChangeRate > alert.Threshold
	case models.ChangeRateBelow:
		return stock.ChangeRate < alert.Threshold
	default:
		return false
	}
}

// CreateAlert creates a new stock alert.
func (s *StockMonitorService) CreateAlert(ctx context.Context, stockCode string, alertType models.AlertType, threshold float64) error {
	if stockCode == "" {
		return fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	// Validate alert type
	switch alertType {
	case models.PriceAbove, models.PriceBelow, models.ChangeRateAbove, models.ChangeRateBelow:
		// Valid
	default:
		return fmt.Errorf("%w: invalid alert type: %s", ErrInvalidInput, alertType)
	}

	// Get the stock
	stock, err := s.stockRepo.GetByCode(ctx, stockCode)
	if err != nil {
		return fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, stockCode)
	}

	// Create the alert
	alert := &models.StockAlert{
		StockID:   stock.ID,
		AlertType: alertType,
		Threshold: threshold,
		IsActive:  true,
		Triggered: false,
		CreatedAt: time.Now(),
	}

	// Save the alert
	if err := s.stockRepo.SaveAlert(ctx, alert); err != nil {
		return fmt.Errorf("failed to save alert: %w", err)
	}

	return nil
}

// DeleteAlert deletes a stock alert.
func (s *StockMonitorService) DeleteAlert(ctx context.Context, alertID string) error {
	if alertID == "" {
		return fmt.Errorf("%w: alert ID is required", ErrInvalidInput)
	}

	// Delete the alert
	if err := s.stockRepo.DeleteAlert(ctx, alertID); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}

// GetAlerts retrieves all alerts for a stock.
func (s *StockMonitorService) GetAlerts(ctx context.Context, stockCode string) ([]*models.StockAlert, error) {
	if stockCode == "" {
		return nil, fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	// Get the stock
	stock, err := s.stockRepo.GetByCode(ctx, stockCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return nil, fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, stockCode)
	}

	// Get alerts for the stock
	alerts, err := s.stockRepo.GetAlerts(ctx, stock.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	return alerts, nil
}

// EnableAlert enables a stock alert.
func (s *StockMonitorService) EnableAlert(ctx context.Context, alertID string) error {
	if alertID == "" {
		return fmt.Errorf("%w: alert ID is required", ErrInvalidInput)
	}

	// Get alerts for the stock
	alerts, err := s.stockRepo.GetActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get alerts: %w", err)
	}

	// Find the alert
	var found bool
	var alert *models.StockAlert
	for _, a := range alerts {
		if a.ID == alertID {
			alert = a
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%w: alert with ID %s not found", ErrStockNotFound, alertID)
	}

	// Update alert status
	if err := s.stockRepo.UpdateAlertStatus(ctx, alertID, true, alert.Triggered, alert.LastTriggered); err != nil {
		return fmt.Errorf("failed to update alert status: %w", err)
	}

	return nil
}

// DisableAlert disables a stock alert.
func (s *StockMonitorService) DisableAlert(ctx context.Context, alertID string) error {
	if alertID == "" {
		return fmt.Errorf("%w: alert ID is required", ErrInvalidInput)
	}

	// Get alerts for the stock
	alerts, err := s.stockRepo.GetActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get alerts: %w", err)
	}

	// Find the alert
	var found bool
	var alert *models.StockAlert
	for _, a := range alerts {
		if a.ID == alertID {
			alert = a
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%w: alert with ID %s not found", ErrStockNotFound, alertID)
	}

	// Update alert status
	if err := s.stockRepo.UpdateAlertStatus(ctx, alertID, false, alert.Triggered, alert.LastTriggered); err != nil {
		return fmt.Errorf("failed to update alert status: %w", err)
	}

	return nil
}
