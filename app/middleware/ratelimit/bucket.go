package ratelimit

type RateLimiter struct {
	BucketSize int
	LeakRatePreMinute int
}
