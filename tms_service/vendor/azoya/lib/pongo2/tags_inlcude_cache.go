package pongo2

import "errors"

type Cache interface {
	FromFile(path string) (value []byte, err error)
}

type emptyCache struct {
}

func (p *emptyCache) FromFile(string) ([]byte, error) {
	return []byte{}, errors.New("empty cahce instance")
}

var (
	emptyCacheInstance Cache = &emptyCache{}
	IncludeCache       Cache
)

func SetIncludeCache(cache Cache) error {

	if cache == nil {
		return errors.New("SetIncludeCache: set cache instance is nil")
	}

	IncludeCache = cache

	return nil
}

func GetIncludeCache() Cache {

	if IncludeCache != nil {
		return IncludeCache
	}

	return emptyCacheInstance
}
