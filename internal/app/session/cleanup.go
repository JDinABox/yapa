package session

import (
	"context"
	"time"

	"github.com/JDinABox/yapa/internal/ilog"
)

type Cleanup struct {
	quit chan struct{}
	done chan struct{}
}

func (s *Store) Cleanup(interval time.Duration) *Cleanup {
	c := &Cleanup{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		ctx := context.Background()
		for {
			select {
			case <-ticker.C:
				if err := s.db.CleanupSessions(ctx, time.Now().Unix()); err != nil {
					ilog.Error(err)
				}
			case <-c.quit:
				c.done <- struct{}{}
				return
			}
		}
	}()
	return c
}

func (c *Cleanup) Stop() {
	c.quit <- struct{}{}
	<-c.done
}
