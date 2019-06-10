package shouts

import (
	"context"

	"github.com/pocockn/models/api/shouts"
)

// Usecase represents the shout's repository context.
type Usecase interface {
	Fetch(ctx context.Context, id string) (shouts shouts.Shout, err error)
	Store(ctx context.Context, shout *shouts.Shout) error
}
