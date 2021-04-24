package loop

import (
	"time"

	"github.com/karlseguin/ccache"
)

var lruCache *ccache.Cache

func init() {
	lruCache = ccache.New(ccache.Configure().MaxSize(100))
}

func GetNext(name string, ips []string) string {
	current := 0
	if item := lruCache.Get(name); item != nil && !item.Expired() && item.Value() != nil {
		current = item.Value().(int)
		if current == 0 {
			current = 1
		} else {
			current = 0
		}
	}
	lruCache.Set(name, current, 10*time.Minute)
	return ips[current]
}

func GetCurrent(name string, ips []string) string {
	current := 0
	if item := lruCache.Get(name); item != nil && !item.Expired() && item.Value() != nil {
		current = item.Value().(int)
	}
	return ips[current]
}
