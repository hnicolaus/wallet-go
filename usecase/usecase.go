// This file contains the usecase implementation layer.
package usecase

import (
	"github.com/WalletService/repository"
)

type Usecase struct {
	Repository repository.RepositoryInterface
}

func NewUsecase(repo repository.RepositoryInterface) *Usecase {
	return &Usecase{
		Repository: repo,
	}
}
