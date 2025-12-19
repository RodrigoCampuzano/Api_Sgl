package inventory

import (
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// GetStockUseCase implementa HU-05: Monitor de stock en tiempo real
type GetStockUseCase struct {
	inventoryRepo domain.InventoryRepository
	productRepo   domain.ProductRepository
}

func NewGetStockUseCase(
	inventoryRepo domain.InventoryRepository,
	productRepo domain.ProductRepository,
) *GetStockUseCase {
	return &GetStockUseCase{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

type StockItem struct {
	ProductID         uuid.UUID           `json:"product_id"`
	ProductName       string              `json:"product_name"`
	SKU               string              `json:"sku"`
	Brand             domain.Brand        `json:"brand"`
	Category          string              `json:"category"`
	TotalStock        int                 `json:"total_stock"`
	AvailableStock    int                 `json:"available_stock"`
	ReservedStock     int                 `json:"reserved_stock"`
	LowStockAlert     bool                `json:"low_stock_alert"`
	ExpirationWarning bool                `json:"expiration_warning"`
	Lots              []*domain.Inventory `json:"lots,omitempty"`
}

func (uc *GetStockUseCase) Execute(brand *domain.Brand, category *string) ([]*StockItem, error) {
	// Obtener productos filtrados
	filters := make(map[string]interface{})
	if brand != nil {
		filters["brand"] = *brand
	}
	if category != nil {
		filters["category"] = *category
	}

	products, err := uc.productRepo.List(filters, 100, 0)
	if err != nil {
		return nil, err
	}

	var stockItems []*StockItem
	for _, product := range products {
		// Obtener inventario del producto
		lots, _ := uc.inventoryRepo.FindByProduct(product.ID)

		totalStock := 0
		availableStock := 0
		reservedStock := 0
		expirationWarning := false

		for _, lot := range lots {
			totalStock += lot.Quantity

			if lot.Status == domain.StockDisponible {
				availableStock += lot.Quantity
			} else if lot.Status == domain.StockReservado {
				reservedStock += lot.Quantity
			}

			// Verificar caducidad < 30 días (HU-06)
			if lot.ExpirationDate != nil {
				daysUntilExp := lot.DaysUntilExpiration()
				if daysUntilExp != nil && *daysUntilExp < 30 {
					expirationWarning = true
				}
			}
		}

		// Indicador de punto de reorden (stock bajo) - ejemplo: < 10 unidades
		lowStockAlert := availableStock < 10 && availableStock > 0

		stockItems = append(stockItems, &StockItem{
			ProductID:         product.ID,
			ProductName:       product.Name,
			SKU:               product.SKU,
			Brand:             product.Brand,
			Category:          product.Category,
			TotalStock:        totalStock,
			AvailableStock:    availableStock,
			ReservedStock:     reservedStock,
			LowStockAlert:     lowStockAlert,
			ExpirationWarning: expirationWarning,
			Lots:              lots,
		})
	}

	return stockItems, nil
}
