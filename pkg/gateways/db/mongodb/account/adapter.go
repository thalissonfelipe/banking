package account

import (
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type accountAdpater struct {
	ID        vos.ID     `bson:"id"`
	Name      string     `bson:"name"`
	CPF       vos.CPF    `bson:"cpf"`
	Secret    vos.Secret `bson:"secret"`
	Balance   int        `bson:"balance"`
	CreatedAt time.Time  `bson:"created_at"`
}

func (a accountAdpater) convertToDomainAccount() *entities.Account {
	return &entities.Account{
		ID:        a.ID,
		Name:      a.Name,
		CPF:       a.CPF,
		Secret:    a.Secret,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}
