package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	echo "github.com/labstack/echo/v4"
)

type Stats struct {
	Uptime       time.Time      `json:"uptime"`
	RequestCount uint64         `json:"request_count"`
	Statuses     map[string]int `json:"statuses"`
	mutex        sync.RWMutex
}

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
	}
}

func (s *Stats) Process() echo.MiddlewareFunc {
	return s.process
}

// Process is the middleware function.
func (s *Stats) process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			s.RequestCount++
			status := strconv.Itoa(ctx.Response().Status)
			s.Statuses[status]++
		}()

		if err := next(ctx); err != nil {
			return err
		}

		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Stats) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}
