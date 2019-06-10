package shouts

import (
	"github.com/pocockn/models/api/shouts"
)

// Repository represents the shout's repository context.
type Repository interface {
	Fetch(id uint) (shouts shouts.Shout, err error)
	FetchAll() (shouts []shouts.Shout, err error)
	Store(shout *shouts.Shout) error
}
