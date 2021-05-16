package usecase

import "github.com/thalissonfelipe/banking/pkg/domain/transfer"

type Transfer struct {
	repository transfer.Repository
}

func NewTransfer(repo transfer.Repository) *Transfer {
	return &Transfer{repository: repo}
}
