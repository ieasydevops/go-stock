package repositories

import (
	"context"
	"time"

	"go-stock/internal/domain/models"
)

// StockRepository defines the interface for stock data access.
type StockRepository interface {
	// GetByID retrieves a stock by its ID.
	GetByID(ctx context.Context, id string) (*models.StockInfo, error)

	// GetByCode retrieves a stock by its code.
	GetByCode(ctx context.Context, code string) (*models.StockInfo, error)

	// Search searches for stocks based on query.
	Search(ctx context.Context, query string, limit int) ([]*models.StockInfo, error)

	// ListByExchange lists all stocks for a given exchange.
	ListByExchange(ctx context.Context, exchange string) ([]*models.StockInfo, error)

	// ListFollowed lists all stocks that the user is following.
	ListFollowed(ctx context.Context) ([]*models.StockInfo, error)

	// ListByGroup lists all stocks in a specific group.
	ListByGroup(ctx context.Context, groupID string) ([]*models.StockInfo, error)

	// Save saves a stock.
	Save(ctx context.Context, stock *models.StockInfo) error

	// Delete deletes a stock.
	Delete(ctx context.Context, id string) error

	// GetHistoricalPrices retrieves historical price data for a stock.
	GetHistoricalPrices(ctx context.Context, stockID string, startDate, endDate time.Time) ([]*models.HistoricalPrice, error)

	// SaveHistoricalPrice saves a historical price record.
	SaveHistoricalPrice(ctx context.Context, price *models.HistoricalPrice) error

	// BulkSaveHistoricalPrices saves multiple historical price records in a batch.
	BulkSaveHistoricalPrices(ctx context.Context, prices []*models.HistoricalPrice) error

	// GetAlerts retrieves all alerts for a stock.
	GetAlerts(ctx context.Context, stockID string) ([]*models.StockAlert, error)

	// SaveAlert saves an alert.
	SaveAlert(ctx context.Context, alert *models.StockAlert) error

	// DeleteAlert deletes an alert.
	DeleteAlert(ctx context.Context, alertID string) error

	// GetActiveAlerts retrieves all active alerts.
	GetActiveAlerts(ctx context.Context) ([]*models.StockAlert, error)

	// UpdateAlertStatus updates the status of an alert.
	UpdateAlertStatus(ctx context.Context, alertID string, isActive, triggered bool, lastTriggered time.Time) error
}

// FollowedStockRepository defines the interface for followed stock data access.
type FollowedStockRepository interface {
	// Add adds a stock to the followed list.
	Add(ctx context.Context, followedStock *models.FollowedStock) error

	// Remove removes a stock from the followed list.
	Remove(ctx context.Context, stockID string) error

	// ListAll lists all followed stocks.
	ListAll(ctx context.Context) ([]*models.FollowedStock, error)

	// GetByStockID retrieves a followed stock by stock ID.
	GetByStockID(ctx context.Context, stockID string) (*models.FollowedStock, error)

	// UpdateNote updates the note for a followed stock.
	UpdateNote(ctx context.Context, stockID, note string) error

	// UpdateWatchingStatus updates the watching status for a followed stock.
	UpdateWatchingStatus(ctx context.Context, stockID string, isWatching bool) error

	// AssignToGroup assigns a followed stock to a group.
	AssignToGroup(ctx context.Context, stockID, groupID string) error
}

// StockGroupRepository defines the interface for stock group data access.
type StockGroupRepository interface {
	// Create creates a new stock group.
	Create(ctx context.Context, group *models.StockGroup) error

	// Update updates a stock group.
	Update(ctx context.Context, group *models.StockGroup) error

	// Delete deletes a stock group.
	Delete(ctx context.Context, groupID string) error

	// GetByID retrieves a stock group by its ID.
	GetByID(ctx context.Context, groupID string) (*models.StockGroup, error)

	// ListAll lists all stock groups.
	ListAll(ctx context.Context) ([]*models.StockGroup, error)
}
