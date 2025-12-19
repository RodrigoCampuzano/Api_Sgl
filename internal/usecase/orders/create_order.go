package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// CreateOrderUseCase implementa HU-07, HU-08, HU-09
type CreateOrderUseCase struct {
	orderRepo     domain.OrderRepository
	orderLineRepo domain.OrderLineRepository
	customerRepo  domain.CustomerRepository
	productRepo   domain.ProductRepository
	inventoryRepo domain.InventoryRepository
	auditRepo     domain.AuditRepository
}

func NewCreateOrderUseCase(
	orderRepo domain.OrderRepository,
	orderLineRepo domain.OrderLineRepository,
	customerRepo domain.CustomerRepository,
	productRepo domain.ProductRepository,
	inventoryRepo domain.InventoryRepository,
	auditRepo domain.AuditRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:     orderRepo,
		orderLineRepo: orderLineRepo,
		customerRepo:  customerRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		auditRepo:     auditRepo,
	}
}

type OrderLineInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CreateOrderInput struct {
	CustomerID uuid.UUID        `json:"customer_id"`
	Lines      []OrderLineInput `json:"lines"`
	UserID     uuid.UUID        `json:"-"`
}

type CreateOrderOutput struct {
	Order             *domain.Order       `json:"order"`
	Lines             []*domain.OrderLine `json:"lines"`
	SuggestedVehicle  domain.VehicleType  `json:"suggested_vehicle"`
	LoadingAlert      string              `json:"loading_alert,omitempty"`
	LoadingEfficiency float64             `json:"loading_efficiency"` // HU-18
	BrandsMixed       []domain.Brand      `json:"brands_mixed"`       // HU-07
}

func (uc *CreateOrderUseCase) Execute(input CreateOrderInput) (*CreateOrderOutput, error) {
	// 1. Verificar cliente
	customer, err := uc.customerRepo.FindByID(input.CustomerID)
	if err != nil {
		return nil, errors.New("cliente no encontrado")
	}

	if !customer.IsActive {
		return nil, errors.New("cliente inactivo")
	}

	// 2. Generar número de pedido
	orderNumber := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	// 3. Procesar líneas y calcular métricas
	var orderLines []*domain.OrderLine
	var totalCost float64
	var totalWeightKg float64
	var totalVolumeM3 float64
	hasFragile := false
	hasHeavy := false
	brandMap := make(map[domain.Brand]bool)

	for _, lineInput := range input.Lines {
		// Obtener producto
		product, err := uc.productRepo.FindByID(lineInput.ProductID)
		if err != nil {
			return nil, fmt.Errorf("producto %s no encontrado", lineInput.ProductID)
		}

		// HU-07: Detectar mezcla de marcas
		brandMap[product.Brand] = true

		// Calcular pesos y volúmenes
		lineWeight := product.WeightKg * float64(lineInput.Quantity)
		lineVolume := product.CalculateVolume() * float64(lineInput.Quantity)
		subtotal := product.UnitPrice * float64(lineInput.Quantity)

		totalWeightKg += lineWeight
		totalVolumeM3 += lineVolume
		totalCost += subtotal

		// HU-09: Detectar productos frágiles y pesados
		if product.IsFragile {
			hasFragile = true
		}
		if product.WeightKg > 10 { // > 10kg considerado pesado
			hasHeavy = true
		}

		// Reservar inventario FEFO
		lots, _ := uc.inventoryRepo.FindByProductFEFO(product.ID)
		var inventoryID *uuid.UUID
		remainingQty := lineInput.Quantity

		for _, lot := range lots {
			if lot.Quantity >= remainingQty {
				inventoryID = &lot.ID
				break
			}
		}

		orderLine := &domain.OrderLine{
			ProductID:   product.ID,
			InventoryID: inventoryID,
			Quantity:    lineInput.Quantity,
			UnitPrice:   product.UnitPrice,
			Subtotal:    subtotal,
		}
		orderLines = append(orderLines, orderLine)
	}

	// HU-07: Verificar si es pedido multi-marca
	var brands []domain.Brand
	for brand := range brandMap {
		brands = append(brands, brand)
	}

	// 4. Crear orden
	order := &domain.Order{
		OrderNumber:     orderNumber,
		CustomerID:      input.CustomerID,
		Status:          domain.OrderBorrador,
		TotalWeightKg:   totalWeightKg,
		TotalVolumeM3:   totalVolumeM3,
		TotalCost:       totalCost,
		HasFragileItems: hasFragile,
		HasHeavyItems:   hasHeavy,
		CreatedBy:       input.UserID,
	}

	// HU-08: Sugerencia automática de vehículo
	suggestedVehicle := order.SuggestVehicle()
	order.SuggestedVehicle = &suggestedVehicle

	// HU-09: Generar alerta de estiba
	loadingAlert := order.GenerateLoadingAlert()
	order.LoadingAlert = loadingAlert

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// 5. Crear líneas
	for _, line := range orderLines {
		line.OrderID = order.ID
	}
	if err := uc.orderLineRepo.CreateBatch(orderLines); err != nil {
		return nil, err
	}

	// HU-18: Calcular eficiencia de carga
	var vehicleCapacity float64
	switch suggestedVehicle {
	case domain.VehicleVan:
		vehicleCapacity = 10.0
	case domain.VehicleCamion35:
		vehicleCapacity = 20.0
	case domain.VehicleTorton:
		vehicleCapacity = 40.0
	default:
		vehicleCapacity = 10.0
	}

	loadingEfficiency := (totalVolumeM3 / vehicleCapacity) * 100
	if loadingEfficiency > 100 {
		loadingEfficiency = 100
	}

	// 6. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "CREATE_ORDER",
		EntityType: "ORDER",
		EntityID:   &order.ID,
		NewValues: map[string]interface{}{
			"order_number":       orderNumber,
			"total_volume_m3":    totalVolumeM3,
			"suggested_vehicle":  suggestedVehicle,
			"loading_efficiency": loadingEfficiency,
			"brands_mixed":       len(brands) > 1,
		},
	})

	return &CreateOrderOutput{
		Order:             order,
		Lines:             orderLines,
		SuggestedVehicle:  suggestedVehicle,
		LoadingAlert:      loadingAlert,
		LoadingEfficiency: loadingEfficiency,
		BrandsMixed:       brands,
	}, nil
}
