package job

import (
	"context"
	"time"
)

type Job struct {
	ID        string
	From, To  time.Time
	CreatedAt time.Time
	Cancel    context.CancelFunc
}
