package envcache

import (
	"fmt"
	"os"
)

var cache = map[string]any{}

// CacheEntries caches all of the provided entries inside of an internal map.
// The entries map must consist of keys that correspond to the environment
// variable name, and the value should be a function converting the value
// to an acceptable type, or nil if no conversion is needed
func CacheNewEntries(entries map[string]func(string) (any, error)) error {
	for key, f := range entries {
		if f == nil {
			val, ok := os.LookupEnv(key)
			if ok {
				cache[key] = val
				continue
			} else {
				return fmt.Errorf("environment key %q not found", key)
			}
		}
		res, err := GetFromEnv(key, f)
		if err != nil {
			return err
		}
		cache[key] = res
	}
	return nil
}

// Get returns the value for a previously cached environment variable.
// NOTE: Developers has to be ensured that the environment variable previously
// have been cached and converted to the appropriate type. Otherwise this function
// will panic
func Get[T any](key string) T {
	val, ok := cache[key]
	if !ok {
		panic(fmt.Sprintf("developer error, key %q not found", key))
	}

	res, ok := val.(T)
	if !ok {
		panic(fmt.Sprintf("developer error, could not cast %q to type %T", val, res))
	}

	return res
}

// GetFromEnv uses os.LookupEnv to fetch the value of the environment variable
// and converts it to an appropriate type using the provided function (cannot be nil).
func GetFromEnv[T any](key string, f func(string) (T, error)) (T, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		res, _ := f(val)
		return res, fmt.Errorf("environment key %q not found", key)
	}
	return f(val)
}
