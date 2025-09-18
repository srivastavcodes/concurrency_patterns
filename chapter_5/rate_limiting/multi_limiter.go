package main

import (
	"context"
	"slices"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

type MultiLimiter struct {
	limiters []RateLimiter
}

func NewMultiLimiter(limiters ...RateLimiter) *MultiLimiter {
	byLimit := func(a, b RateLimiter) int {
		switch {
		case a.Limit() > b.Limit():
			return 1
		case a.Limit() < b.Limit():
			return -1
		default:
			return 0
		}
	}
	slices.SortFunc(limiters, byLimit)
	return &MultiLimiter{limiters: limiters}
}

func (ml *MultiLimiter) Wait(ctx context.Context) error {
	for _, lim := range ml.limiters {
		if err := lim.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ml *MultiLimiter) Limit() rate.Limit {
	return ml.limiters[0].Limit()
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
