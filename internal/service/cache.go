package service

import (
	"context"
	"time"

	"github.com/tupyy/osdviewer/internal/entity"
)

type FleetManagerReader interface {
	GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error)
}

// wrapper around FleetManager for cache
type FleetManagerCache struct {
	lastRead time.Time
	cacheTtl time.Duration
	cache    map[Environment][]entity.Cluster
	reader   FleetManagerReader
}

func NewDefaultFleetManagerCache(reader FleetManagerReader) *FleetManagerCache {
	return NewFleetManagerCache(reader, defaultCacheTTL)
}

func NewFleetManagerCache(reader FleetManagerReader, cacheTtl time.Duration) *FleetManagerCache {
	return &FleetManagerCache{
		lastRead: time.Now(),
		cacheTtl: cacheTtl,
		reader:   reader,
		cache:    make(map[Environment][]entity.Cluster),
	}
}

func (c *FleetManagerCache) GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error) {
	if c.lastRead.Add(c.cacheTtl).Before(time.Now()) {
		if clusters, ok := c.cache[env]; ok {
			return clusters, nil
		}
	}

	clusters, err := c.reader.GetClusters(ctx, env)
	if err != nil {
		return nil, err
	}

	// save to cache
	c.cache[env] = clusters
	c.lastRead = time.Now()

	return clusters, nil
}
