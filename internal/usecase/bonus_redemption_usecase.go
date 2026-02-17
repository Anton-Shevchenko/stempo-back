package usecase

import (
	"errors"
	"time"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"github.com/stempo/backend/pkg/utils"
)

type BonusRedemptionUsecase interface {
	GenerateRedemptionCode(cardID, userID uint) (*entity.BonusRedemption, error)
	RedeemBonus(code string, scannedByUserID uint) (*entity.BonusRedemption, error)
	GetRedemptionsByCard(cardID, userID uint) ([]entity.BonusRedemption, error)
}

type bonusRedemptionUsecase struct {
	redemptionRepo repository.BonusRedemptionRepository
	cardRepo       repository.LoyaltyCardRepository
	businessRepo   repository.BusinessRepository
	employeeUsecase EmployeeUsecase
}

func NewBonusRedemptionUsecase(
	redemptionRepo repository.BonusRedemptionRepository,
	cardRepo repository.LoyaltyCardRepository,
	businessRepo repository.BusinessRepository,
	employeeUsecase EmployeeUsecase,
) BonusRedemptionUsecase {
	return &bonusRedemptionUsecase{
		redemptionRepo: redemptionRepo,
		cardRepo:       cardRepo,
		businessRepo:   businessRepo,
		employeeUsecase: employeeUsecase,
	}
}

func (u *bonusRedemptionUsecase) GenerateRedemptionCode(cardID, userID uint) (*entity.BonusRedemption, error) {
	// Get the loyalty card
	card, err := u.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	// Verify that the user owns the card
	if card.UserID != userID {
		return nil, errors.New("unauthorized: you can only generate codes for your own cards")
	}

	// Check if user has enough stamps
	if card.Stamps < card.StampsRequired {
		return nil, errors.New("insufficient stamps: you need more stamps to use this bonus")
	}

	// Check if there's already an active (unused and not expired) redemption code
	existingRedemptions, err := u.redemptionRepo.FindByCardID(cardID)
	if err == nil {
		now := time.Now()
		for _, redemption := range existingRedemptions {
			if !redemption.Used && redemption.ExpiresAt.After(now) {
				return nil, errors.New("you already have an active redemption code")
			}
		}
	}

	// Generate unique code
	code, err := utils.GenerateUniqueCode()
	if err != nil {
		return nil, errors.New("failed to generate code")
	}

	// Create redemption with 1 minute expiration
	redemption := &entity.BonusRedemption{
		CardID:    cardID,
		Code:      code,
		ExpiresAt: time.Now().Add(1 * time.Minute),
		Used:      false,
	}

	if err := u.redemptionRepo.Create(redemption); err != nil {
		return nil, err
	}

	// Reload with relations
	return u.redemptionRepo.FindByCode(code)
}

func (u *bonusRedemptionUsecase) RedeemBonus(code string, scannedByUserID uint) (*entity.BonusRedemption, error) {
	// Find redemption by code
	redemption, err := u.redemptionRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("invalid redemption code")
	}

	// Check if already used
	if redemption.Used {
		return nil, errors.New("this code has already been used")
	}

	// Check if expired
	if time.Now().After(redemption.ExpiresAt) {
		return nil, errors.New("this code has expired")
	}

	// Get the card
	card, err := u.cardRepo.FindByID(redemption.CardID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	// Verify that scannedByUserID is the business owner or an employee
	business, err := u.businessRepo.FindByID(card.BusinessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	// Check if user is business owner
	if business.OwnerID != scannedByUserID {
		// Check if user is an employee
		isEmployee, err := u.employeeUsecase.IsEmployee(card.BusinessID, scannedByUserID)
		if err != nil || !isEmployee {
			return nil, errors.New("unauthorized: only business owner or employee can redeem bonuses")
		}
	}

	// Verify that card still has enough stamps (double-check)
	if card.Stamps < card.StampsRequired {
		return nil, errors.New("insufficient stamps")
	}

	// Mark as used and update card
	now := time.Now()
	redemption.Used = true
	redemption.UsedAt = &now
	redemption.ScannedBy = &scannedByUserID

	// Reset stamps to 0 (spend the collected stamps)
	card.Stamps = 0

	// Update redemption
	if err := u.redemptionRepo.Update(redemption); err != nil {
		return nil, err
	}

	// Update card
	if err := u.cardRepo.Update(card); err != nil {
		return nil, err
	}

	// Reload redemption with relations
	return u.redemptionRepo.FindByCode(code)
}

func (u *bonusRedemptionUsecase) GetRedemptionsByCard(cardID, userID uint) ([]entity.BonusRedemption, error) {
	// Verify that the user owns the card
	card, err := u.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	if card.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return u.redemptionRepo.FindByCardID(cardID)
}
