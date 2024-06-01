package repo

import (
	"sync"

	"github.com/VitoNaychev/eval-web-service/service"
)

type InMemoryExprErrorRepository struct {
	exprErrors map[string]*service.ExpressionError
	mu         sync.Mutex
}

func NewInMemoryExprErrorRepository() *InMemoryExprErrorRepository {
	return &InMemoryExprErrorRepository{
		exprErrors: make(map[string]*service.ExpressionError),
	}
}

func (repo *InMemoryExprErrorRepository) Increment(exprError *service.ExpressionError) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if existing, exists := repo.exprErrors[exprError.Expression]; exists {
		existing.Frequency++
	} else {
		exprError.Frequency = 1
		repo.exprErrors[exprError.Expression] = exprError
	}

	return nil
}

func (repo *InMemoryExprErrorRepository) GetAll() ([]service.ExpressionError, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var allErrors []service.ExpressionError
	for _, exprError := range repo.exprErrors {
		allErrors = append(allErrors, *exprError)
	}

	return allErrors, nil
}
