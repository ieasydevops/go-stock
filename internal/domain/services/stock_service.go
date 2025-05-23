package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-stock/internal/domain/events"
	"go-stock/internal/domain/models"
	"go-stock/internal/domain/repositories"
)

// Common errors
var (
	ErrStockNotFound  = errors.New("stock not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDuplicateStock = errors.New("stock already exists")
)

// StockService provides operations related to stocks.
type StockService struct {
	stockRepo       repositories.StockRepository
	followedRepo    repositories.FollowedStockRepository
	eventDispatcher events.EventDispatcher
}

// NewStockService creates a new stock service.
func NewStockService(
	stockRepo repositories.StockRepository,
	followedRepo repositories.FollowedStockRepository,
	eventDispatcher events.EventDispatcher,
) *StockService {
	return &StockService{
		stockRepo:       stockRepo,
		followedRepo:    followedRepo,
		eventDispatcher: eventDispatcher,
	}
}

// GetStockInfo retrieves a stock by its code.
func (s *StockService) GetStockInfo(ctx context.Context, code string) (*models.StockInfo, error) {
	if code == "" {
		return nil, fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	stock, err := s.stockRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return nil, fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, code)
	}

	return stock, nil
}

// SearchStocks searches for stocks based on a query.
func (s *StockService) SearchStocks(ctx context.Context, query string, limit int) ([]*models.StockInfo, error) {
	if query == "" {
		return nil, fmt.Errorf("%w: search query is required", ErrInvalidInput)
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	stocks, err := s.stockRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search stocks: %w", err)
	}

	return stocks, nil
}

// UpdatePrice updates the price of a stock and triggers appropriate events.
func (s *StockService) UpdatePrice(ctx context.Context, code string, price float64) error {
	if code == "" {
		return fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	if price < 0 {
		return fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}

	stock, err := s.stockRepo.GetByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, code)
	}

	oldPrice := stock.Price

	// Calculate change and change rate
	change := price - oldPrice
	var changeRate float64
	if oldPrice != 0 {
		changeRate = (change / oldPrice) * 100
	}

	// Update stock info
	stock.Price = price
	stock.Change = change
	stock.ChangeRate = changeRate
	stock.UpdateTime = time.Now()

	// Save updated stock
	if err := s.stockRepo.Save(ctx, stock); err != nil {
		return fmt.Errorf("failed to save updated stock: %w", err)
	}

	// Dispatch price changed event if price actually changed
	if oldPrice != price {
		event := events.NewStockPriceChangedEvent(stock, oldPrice, price)
		s.eventDispatcher.Dispatch(event)
	}

	return nil
}

// GetHistoricalPrices retrieves historical price data for a stock.
func (s *StockService) GetHistoricalPrices(ctx context.Context, code string, startDate, endDate time.Time) ([]*models.HistoricalPrice, error) {
	if code == "" {
		return nil, fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	stock, err := s.stockRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return nil, fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, code)
	}

	// If no date range is specified, use the last 7 days
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -7)
	}

	if endDate.IsZero() {
		endDate = time.Now()
	}

	prices, err := s.stockRepo.GetHistoricalPrices(ctx, stock.ID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical prices: %w", err)
	}

	return prices, nil
}

// FollowStock adds a stock to the user's followed list.
func (s *StockService) FollowStock(ctx context.Context, code, groupID, note string) error {
	if code == "" {
		return fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	stock, err := s.stockRepo.GetByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, code)
	}

	// Check if stock is already followed
	existingFollow, err := s.followedRepo.GetByStockID(ctx, stock.ID)
	if err != nil && !errors.Is(err, ErrStockNotFound) {
		return fmt.Errorf("failed to check if stock is already followed: %w", err)
	}

	if existingFollow != nil {
		return fmt.Errorf("%w: stock with code %s is already followed", ErrDuplicateStock, code)
	}

	// Create followed stock
	followedStock := &models.FollowedStock{
		StockID:    stock.ID,
		GroupID:    groupID,
		Note:       note,
		AddedAt:    time.Now(),
		IsWatching: true,
	}

	// Save followed stock
	if err := s.followedRepo.Add(ctx, followedStock); err != nil {
		return fmt.Errorf("failed to add stock to followed list: %w", err)
	}

	return nil
}

// UnfollowStock removes a stock from the user's followed list.
func (s *StockService) UnfollowStock(ctx context.Context, code string) error {
	if code == "" {
		return fmt.Errorf("%w: stock code is required", ErrInvalidInput)
	}

	stock, err := s.stockRepo.GetByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get stock by code: %w", err)
	}

	if stock == nil {
		return fmt.Errorf("%w: stock with code %s not found", ErrStockNotFound, code)
	}

	// Remove from followed stocks
	if err := s.followedRepo.Remove(ctx, stock.ID); err != nil {
		return fmt.Errorf("failed to remove stock from followed list: %w", err)
	}

	return nil
}

// GetFollowedStocks retrieves all stocks that the user is following.
func (s *StockService) GetFollowedStocks(ctx context.Context) ([]*models.StockInfo, error) {
	stocks, err := s.stockRepo.ListFollowed(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list followed stocks: %w", err)
	}

	return stocks, nil
}
