// Copyright 2015 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package tree

import (
	"regexp"
	"regexp/syntax"

	"github.com/cockroachdb/cockroachdb-parser/pkg/util/cache"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/syncutil"
)

// RegexpCacheKey allows cache keys to take the form of different types,
// as long as they are comparable and can produce a pattern when needed
// for regexp compilation. The pattern method will not be called until
// after a cache lookup is performed and the result is a miss.
type RegexpCacheKey interface {
	Pattern() (string, error)
}

// A RegexpCache is a cache used to store compiled regular expressions.
// The cache is safe for concurrent use by multiple goroutines. It is also
// safe to use the cache through a nil reference, where it will act like a valid
// cache with no capacity.
type RegexpCache struct {
	mu    syncutil.Mutex
	cache *cache.UnorderedCache
}

// NewRegexpCache creates a new RegexpCache of the given size.
// The underlying cache internally uses a hash map, so lookups
// are cheap.
func NewRegexpCache(size int) *RegexpCache {
	return &RegexpCache{
		cache: cache.NewUnorderedCache(cache.Config{
			Policy: cache.CacheLRU,
			ShouldEvict: func(s int, key, value interface{}) bool {
				return s > size
			},
		}),
	}
}

// GetRegexp consults the cache for the regular expressions stored for
// the given key, compiling the key's pattern if it is not already
// in the cache.
func (rc *RegexpCache) GetRegexp(key RegexpCacheKey) (*regexp.Regexp, error) {
	// syntax.Perl is the default flag for regexp.Compile.
	return rc.getRegexpInternal(key, syntax.Perl)
}

// GetRegexpWithFlags consults the cache for the regular expressions stored for
// the given key, compiling the key's pattern with the given flags if it is not
// already in the cache.
func (rc *RegexpCache) GetRegexpWithFlags(
	key RegexpCacheKey, flags syntax.Flags,
) (*regexp.Regexp, error) {
	return rc.getRegexpInternal(key, flags)
}

func (rc *RegexpCache) getRegexpInternal(
	key RegexpCacheKey, flags syntax.Flags,
) (*regexp.Regexp, error) {
	if rc != nil {
		re := rc.lookup(key)
		if re != nil {
			return re, nil
		}
	}

	pattern, err := key.Pattern()
	if err != nil {
		return nil, err
	}
	var re *regexp.Regexp
	if flags == syntax.Perl {
		// Avoid the redundant 'parse - stringify - parse (within Compile)'
		// sequence in the common case.
		re, err = regexp.Compile(pattern)
	} else {
		var parsed *syntax.Regexp
		parsed, err = syntax.Parse(pattern, flags)
		if err != nil {
			return nil, err
		}
		re, err = regexp.Compile(parsed.String())
	}
	if err != nil {
		return nil, err
	}

	if rc != nil {
		rc.update(key, re)
	}
	return re, nil
}

// lookup checks for the regular expression in the cache in a
// synchronized manner, returning it if it exists.
func (rc *RegexpCache) lookup(key RegexpCacheKey) *regexp.Regexp {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	v, ok := rc.cache.Get(key)
	if !ok {
		return nil
	}
	return v.(*regexp.Regexp)
}

// update invalidates the regular expression for the given pattern.
// If a new regular expression is passed in, it is inserted into the cache.
func (rc *RegexpCache) update(key RegexpCacheKey, re *regexp.Regexp) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cache.Del(key)
	if re != nil {
		rc.cache.Add(key, re)
	}
}

// Len returns the number of compiled regular expressions in the cache.
func (rc *RegexpCache) Len() int {
	if rc == nil {
		return 0
	}
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.cache.Len()
}
