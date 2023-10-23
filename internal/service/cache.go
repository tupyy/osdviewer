package service

import (
	"context"
	"time"

	"github.com/tupyy/osdviewer/internal/entity"
)

type FleetManagerReader interface {
	GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error)
}

type cacheData struct {
	Data     []entity.Cluster
	ttl      time.Duration
	lastRead time.Time
}

func (c *cacheData) NextHit() time.Time {
	return c.lastRead.Add(c.ttl)
}

func newCacheData(data []entity.Cluster, ttl time.Duration) *cacheData {
	return &cacheData{
		Data:     data,
		ttl:      ttl,
		lastRead: time.Now(),
	}
}

// wrapper around FleetManager for cache
type FleetManagerCache struct {
	cacheTtl time.Duration
	cache    map[Environment]*cacheData
	reader   FleetManagerReader
}

func NewDefaultFleetManagerCache(reader FleetManagerReader) *FleetManagerCache {
	return NewFleetManagerCache(reader, defaultCacheTTL)
}

func NewFleetManagerCache(reader FleetManagerReader, cacheTtl time.Duration) *FleetManagerCache {
	return &FleetManagerCache{
		cacheTtl: cacheTtl,
		reader:   reader,
		cache:    make(map[Environment]*cacheData),
	}
}

func (c *FleetManagerCache) GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error) {
	if data, ok := c.cache[env]; ok {
		if data.NextHit().After(time.Now()) {
			return data.Data, nil
		}
	}

	clusters, err := c.reader.GetClusters(ctx, env)
	if err != nil {
		return nil, err
	}

	// save to cache
	c.cache[env] = newCacheData(clusters, c.cacheTtl)

	return clusters, nil
}
