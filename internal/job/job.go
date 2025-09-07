package job

import (
	"context"
	"time"
)

// Job represents a settlement request.
type Job struct {
	ID        string
	From, To  time.Time
	CreatedAt time.Time
	Cancel    context.CancelFunc
}
