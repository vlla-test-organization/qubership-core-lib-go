package logging

import (
	"context"
	"time"
)

// A Record is what a Logger asks its handler to write
type Record struct {
	PackageName string
	Time        time.Time
	Lvl         Lvl
	Message     string
	Ctx         context.Context
}
