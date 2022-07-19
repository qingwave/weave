package ratelimit

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/time/rate"
)

const defaultCacheSize = 2048

type LimitType string

const (
	ServerLimitType LimitType = "server"
	IPLimitType     LimitType = "ip"
)

type LimitConfig struct {
	LimitType LimitType `yaml:"limitType"`
	Burst     int       `yaml:"burst"`
	QPS       int       `yaml:"qps"`
	CacheSize int       `yaml:"cacheSize"`
}

func (c *LimitConfig) Validate() error {
	if c.QPS == 0 || c.Burst == 0 {
		return fmt.Errorf("LimitConfig Burst and QPS cannot be empty")
	}
	if c.QPS > c.Burst {
		return fmt.Errorf("LimitConfig QPS(%d) must less than Burst(%d)", c.QPS, c.Burst)
	}
	if c.CacheSize == 0 {
		c.CacheSize = defaultCacheSize
	}
	return nil
}

type RateLimiter struct {
	limitType          LimitType
	keyFunc            func(*gin.Context) interface{}
	cache              *lru.Cache
	rateLimiterFactory func() *rate.Limiter
}

func NewRateLimiter(conf *LimitConfig) (*RateLimiter, error) {
	if conf == nil {
		return nil, errors.New("invaild config")
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	var keyFunc func(*gin.Context) interface{}
	switch conf.LimitType {
	case ServerLimitType:
		keyFunc = func(c *gin.Context) interface{} {
			return ""
		}
	case IPLimitType:
		keyFunc = func(c *gin.Context) interface{} {
			return c.ClientIP()
		}
	default:
		return nil, fmt.Errorf("unknow limit type %s", conf.LimitType)
	}

	c, err := lru.New(conf.CacheSize)
	if err != nil {
		return nil, err
	}

	rateLimiterFactory := func() *rate.Limiter {
		return rate.NewLimiter(rate.Limit(conf.QPS), conf.Burst)
	}

	return &RateLimiter{
		limitType:          conf.LimitType,
		keyFunc:            keyFunc,
		cache:              c,
		rateLimiterFactory: rateLimiterFactory,
	}, nil
}

func (rl *RateLimiter) Accept(c *gin.Context) error {
	key := rl.keyFunc(c)
	limiter := rl.get(key)

	if !limiter.Allow() {
		return fmt.Errorf("limit reached on %s for key %v", rl.limitType, key)
	}

	return nil
}

func (rl *RateLimiter) get(key interface{}) *rate.Limiter {
	value, found := rl.cache.Get(key)
	if !found {
		new := rl.rateLimiterFactory()
		rl.cache.Add(key, new)
		return new
	}
	return value.(*rate.Limiter)
}
