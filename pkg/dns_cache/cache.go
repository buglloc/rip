package dns_cache

import (
	"net"
	"strconv"
	"time"

	"github.com/karlseguin/ccache"
)

type Cache struct {
	lruCache *ccache.Cache
}

func NewCache() *Cache {
	return &Cache{
		lruCache: ccache.New(ccache.Configure()),
	}
}

func (c *Cache) Get(reqType uint16, domain string) *net.IP {
	if item := c.lruCache.Get(makeKey(reqType, domain)); item != nil && !item.Expired() && item.Value() != nil {
		return item.Value().(*net.IP)
	}
	return nil
}

func (c *Cache) Set(reqType uint16, domain string, ttl time.Duration, ip *net.IP) error {
	key := makeKey(reqType, domain)
	if ip == nil {
		c.lruCache.Delete(key)
		return nil
	}

	c.lruCache.Set(key, ip, ttl)
	return nil
}

func makeKey(reqType uint16, domain string) string {
	return strconv.Itoa(int(reqType)) + domain
}
