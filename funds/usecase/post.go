package usecase

import "context"

type IncreaseCmd func(ctx context.Context, by int) error

func (u UseCase) Increase() IncreaseCmd {
	return func(ctx context.Context, by int) error {
		return nil
	}
}

type DecreaseCmd func(ctx context.Context, by int) error

func (u UseCase) Decrease() DecreaseCmd {
	return func(ctx context.Context, by int) error {
		return nil
	}
}
