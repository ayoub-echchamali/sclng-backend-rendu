package api

import "time"

type Cache struct {
	LastUpdated time.Time
	Content interface{}
}

func (s *ApiServer) ReadCache(path string) (Cache, bool) {
    cacheEntry, exists := s.Cache[path]
    
    if !exists {
        return Cache{}, false // Cache miss
    }
    
    return cacheEntry, true // Cache hit
}

func (s *ApiServer) WriteCache(path string, content interface{}) {
    cacheEntry := Cache{
        LastUpdated: time.Now(),
        Content:     content,
    }
    
    s.Cache[path] = cacheEntry
}