package repository

import "github.com/stempo/backend/internal/domain/entity"

type BonusProgramRepository interface {
	Create(program *entity.BonusProgram) error
	FindByID(id uint) (*entity.BonusProgram, error)
	FindByBusinessID(businessID uint) ([]entity.BonusProgram, error)
	FindActiveByBusinessID(businessID uint) ([]entity.BonusProgram, error)
	FindAll(page, pageSize int) ([]entity.BonusProgram, int64, error)
	FindByStatus(status entity.BonusProgramStatus, page, pageSize int) ([]entity.BonusProgram, int64, error)
	Update(program *entity.BonusProgram) error
	Delete(id uint) error
}
