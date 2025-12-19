package inventory

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// RegisterDamageUseCase implementa HU-13: Registro de Mermas con Foto
type RegisterDamageUseCase struct {
	inventoryRepo domain.InventoryRepository
	movementRepo  domain.InventoryMovementRepository
	auditRepo     domain.AuditRepository
}

func NewRegisterDamageUseCase(
	inventoryRepo domain.InventoryRepository,
	movementRepo domain.InventoryMovementRepository,
	auditRepo domain.AuditRepository,
) *RegisterDamageUseCase {
	return &RegisterDamageUseCase{
		inventoryRepo: inventoryRepo,
		movementRepo:  movementRepo,
		auditRepo:     auditRepo,
	}
}

type RegisterDamageInput struct {
	InventoryID      uuid.UUID `json:"inventory_id" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,min=1"`
	Reason           string    `json:"reason" binding:"required"`
	EvidencePhotoURL string    `json:"evidence_photo_url" binding:"required"` // HU-13: OBLIGATORIO
	UserID           uuid.UUID `json:"-"`
}

type RegisterDamageOutput struct {
	Message           string     `json:"message"`
	AdjustedInventory *uuid.UUID `json:"adjusted_inventory_id"`
	MovementID        uuid.UUID  `json:"movement_id"`
}

func (uc *RegisterDamageUseCase) Execute(input RegisterDamageInput) (*RegisterDamageOutput, error) {
	// HU-13: Validar foto obligatoria
	if input.EvidencePhotoURL == "" {
		return nil, errors.New("evidence_photo_url es obligatorio (HU-13)")
	}

	// Validar inventario
	inventory, err := uc.inventoryRepo.FindByID(input.InventoryID)
	if err != nil {
		return nil, errors.New("inventario no encontrado")
	}

	// Validar cantidad
	if inventory.Quantity < input.Quantity {
		return nil, errors.New("cantidad insuficiente")
	}

	// Descontar
	inventory.Quantity -= input.Quantity
	if err := uc.inventoryRepo.Update(inventory); err != nil {
		return nil, err
	}

	// Registrar movimiento (foto en metadata/reason)
	movement := &domain.InventoryMovement{
		InventoryID:  input.InventoryID,
		MovementType: "MERMA",
		Quantity:     -input.Quantity,
		Reason:       input.Reason + " | Foto: " + input.EvidencePhotoURL,
		PerformedBy:  input.UserID,
	}

	if err := uc.movementRepo.Create(movement); err != nil {
		return nil, err
	}

	// Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "REGISTER_DAMAGE",
		EntityType: "INVENTORY",
		EntityID:   &inventory.ID,
		NewValues: map[string]interface{}{
			"quantity_damaged": input.Quantity,
			"photo":            input.EvidencePhotoURL,
		},
	})

	return &RegisterDamageOutput{
		Message:           "Merma registrada exitosamente",
		AdjustedInventory: &inventory.ID,
		MovementID:        movement.ID,
	}, nil
}

// GetFEFOLotsUseCase implementa HU-06: FEFO (First Expired First Out)
type GetFEFOLotsUseCase struct {
	inventoryRepo domain.InventoryRepository
}

func NewGetFEFOLotsUseCase(inventoryRepo domain.InventoryRepository) *GetFEFOLotsUseCase {
	return &GetFEFOLotsUseCase{
		inventoryRepo: inventoryRepo,
	}
}

type FEFOLot struct {
	InventoryID       uuid.UUID  `json:"inventory_id"`
	LotNumber         string     `json:"lot_number"`
	Quantity          int        `json:"quantity"`
	ExpirationDate    *time.Time `json:"expiration_date,omitempty"`
	DaysUntilExpiry   *int       `json:"days_until_expiry,omitempty"`
	WarehouseLocation string     `json:"warehouse_location"`
	ExpirationAlert   string     `json:"expiration_alert,omitempty"`
}

func (uc *GetFEFOLotsUseCase) Execute(productID uuid.UUID) ([]*FEFOLot, error) {
	// Obtener lotes ordenados por FEFO
	lots, err := uc.inventoryRepo.FindByProductFEFO(productID)
	if err != nil {
		return nil, err
	}

	var fefoLots []*FEFOLot
	for _, lot := range lots {
		fefoLot := &FEFOLot{
			InventoryID:       lot.ID,
			LotNumber:         lot.LotNumber,
			Quantity:          lot.Quantity,
			ExpirationDate:    lot.ExpirationDate,
			WarehouseLocation: lot.WarehouseLocation,
		}

		// Calcular días hasta caducidad
		if lot.ExpirationDate != nil {
			daysUntilExp := lot.DaysUntilExpiration()
			fefoLot.DaysUntilExpiry = daysUntilExp

			if daysUntilExp != nil {
				if *daysUntilExp < 0 {
					fefoLot.ExpirationAlert = "PRODUCTO CADUCADO"
				} else if *daysUntilExp < 30 {
					fefoLot.ExpirationAlert = "ALERTA ROJA: Caducidad menor a 30 días" // HU-06
				} else if *daysUntilExp < 60 {
					fefoLot.ExpirationAlert = "Precaución: Caducidad próxima"
				}
			}
		}

		fefoLots = append(fefoLots, fefoLot)
	}

	return fefoLots, nil
}
