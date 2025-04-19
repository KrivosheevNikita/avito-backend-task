package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
}
