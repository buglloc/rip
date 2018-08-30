package ip_stick

import (
	"github.com/karlseguin/ccache"

	"github.com/buglloc/rip/pkg/cfg"
)

var lruCache *ccache.Cache

func init() {
	lruCache = ccache.New(ccache.Configure().MaxSize(1000))
}

func GetCurrent(key string, ips []string) string {
	if item := lruCache.Get(key); item != nil && !item.Expired() && item.Value() != nil {
		return ips[1]
	}

	lruCache.Set(key, true, cfg.StickyTTL)
	return ips[0]
}
