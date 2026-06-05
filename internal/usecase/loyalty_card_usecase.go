package usecase

import (
	"errors"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type LoyaltyCardUsecase interface {
	Create(card *entity.LoyaltyCard) error
	GetByID(id uint) (*entity.LoyaltyCard, error)
	GetByUserID(userID uint) ([]entity.LoyaltyCard, error)
	GetByBusinessID(businessID uint) ([]entity.LoyaltyCard, error)
	GetByUserIDAndBusinessID(userID, businessID uint) ([]entity.LoyaltyCard, error)
	AddStamp(cardID, scannedByUserID uint) error
	Update(card *entity.LoyaltyCard) error
	Delete(id uint) error
}

type loyaltyCardUsecase struct {
	cardRepo      repository.LoyaltyCardRepository
	businessRepo  repository.BusinessRepository
	employeeUsecase EmployeeUsecase
}

func NewLoyaltyCardUsecase(
	cardRepo repository.LoyaltyCardRepository,
	businessRepo repository.BusinessRepository,
	employeeUsecase EmployeeUsecase,
) LoyaltyCardUsecase {
	return &loyaltyCardUsecase{
		cardRepo:       cardRepo,
		businessRepo:   businessRepo,
		employeeUsecase: employeeUsecase,
	}
}

func (u *loyaltyCardUsecase) Create(card *entity.LoyaltyCard) error {
	existing, _ := u.cardRepo.FindByUserAndBusiness(card.UserID, card.BusinessID)
	if existing != nil {
		return errors.New("loyalty card already exists")
	}

	return u.cardRepo.Create(card)
}

func (u *loyaltyCardUsecase) GetByID(id uint) (*entity.LoyaltyCard, error) {
	return u.cardRepo.FindByID(id)
}

func (u *loyaltyCardUsecase) GetByUserID(userID uint) ([]entity.LoyaltyCard, error) {
	return u.cardRepo.FindByUserID(userID)
}

func (u *loyaltyCardUsecase) GetByBusinessID(businessID uint) ([]entity.LoyaltyCard, error) {
	return u.cardRepo.FindByBusinessID(businessID)
}

func (u *loyaltyCardUsecase) GetByUserIDAndBusinessID(userID, businessID uint) ([]entity.LoyaltyCard, error) {
	return u.cardRepo.FindByUserIDAndBusinessID(userID, businessID)
}

func (u *loyaltyCardUsecase) AddStamp(cardID, scannedByUserID uint) error {
	card, err := u.cardRepo.FindByID(cardID)
	if err != nil {
		return errors.New("card not found")
	}

	// Allow if the user is scanning their own card
	if card.UserID == scannedByUserID {
		card.Stamps++
		return u.cardRepo.Update(card)
	}

	// Check if scannedByUserID is the business owner or an employee
	business, err := u.businessRepo.FindByID(card.BusinessID)
	if err != nil {
		return errors.New("business not found")
	}

	// Check if user is business owner
	if business.OwnerID == scannedByUserID {
		card.Stamps++
		return u.cardRepo.Update(card)
	}

	// Check if user is an employee
	isEmployee, err := u.employeeUsecase.IsEmployee(card.BusinessID, scannedByUserID)
	if err != nil {
		return errors.New("failed to verify employee status")
	}

	if isEmployee {
		card.Stamps++
		return u.cardRepo.Update(card)
	}

	return errors.New("unauthorized: only card owner, business owner, or employee can add stamps")
}

func (u *loyaltyCardUsecase) Update(card *entity.LoyaltyCard) error {
	return u.cardRepo.Update(card)
}

func (u *loyaltyCardUsecase) Delete(id uint) error {
	return u.cardRepo.Delete(id)
}
