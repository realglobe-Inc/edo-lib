package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。

// キャッシュ用。
type cachingDatedServiceExplorer struct {
	DatedServiceExplorer
	cache util.Cache
}

func NewCachingDatedServiceExplorer(backend DatedServiceExplorer) DatedServiceExplorer {
	return &cachingDatedServiceExplorer{DatedServiceExplorer: backend,
		cache: util.NewCache(func(a1 interface{}, a2 interface{}) bool {
			return a1.(*Stamp).ExpiDate.Before(a2.(*Stamp).ExpiDate)
		}),
	}
}

func (reg *cachingDatedServiceExplorer) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	now := time.Now()
	reg.cache.CleanLesser(&Stamp{ExpiDate: now})

	// 残ってるキャッシュは有効。

	val, prio := reg.cache.Get(servUri)
	if prio == nil {
		// キャッシュしてない。
		servUuid, newCaStmp, err = reg.DatedServiceExplorer.StampedServiceUuid(servUri, nil)
		if err != nil {
			return "", nil, erro.Wrap(err)
		} else if newCaStmp == nil {
			// 無い。
			return "", nil, nil
		} else {
			// あった。
			reg.cache.Put(servUri, servUuid, newCaStmp)
			if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
				// 要求元のキャッシュと同じだった。
				return "", newCaStmp, nil
			} else {
				return servUuid, newCaStmp, nil
			}
		}
	}

	// キャッシュしてた。

	stmp := prio.(*Stamp)
	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return "", stmp, nil
	}
	return val.(string), stmp, nil
}
