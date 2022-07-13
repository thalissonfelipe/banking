package services

import "context"

type Auth interface {
	Autheticate(ctx context.Context, cpfStr, secretStr string) (token string, err error)
}
