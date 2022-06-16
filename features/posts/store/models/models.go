package models

import (
	"github.com/k0marov/socnet/core/core_values"
	"time"
)

type PostToCreate struct {
	Author    core_values.UserId
	Text      string
	CreatedAt time.Time
}
