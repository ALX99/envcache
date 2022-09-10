package envcache

import (
	"reflect"
	"strconv"
	"testing"
)

func TestCacheEntries(t *testing.T) {
	makeAndAddToMap := func(keys string, funcs func(string) (any, error)) map[string]func(string) (any, error) {
		m := map[string]func(string) (any, error){}
		m[keys] = funcs
		return m
	}

	t.Setenv("__1", "a")
	t.Setenv("__2", "432")
	type args struct {
		entries map[string]func(string) (any, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test caching existing environment variable", args{makeAndAddToMap("__1", nil)}, false},
		{"test caching existing environment variable with function", args{makeAndAddToMap("__2", func(s string) (any, error) { return strconv.Atoi(s) })}, false},
		{"test caching non-existent environment variable", args{makeAndAddToMap("__3", nil)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CacheNewEntries(tt.args.entries); (err != nil) != tt.wantErr {
				t.Errorf("CacheNewEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	m := map[string]func(string) (any, error){}
	t.Setenv("__1", "32")
	t.Setenv("__2", "t")
	t.Setenv("__3", "false")
	t.Setenv("__4", "abcde")

	m["__1"] = func(s string) (any, error) { return strconv.Atoi(s) }
	m["__2"] = func(s string) (any, error) { return strconv.ParseBool(s) }
	m["__3"] = func(s string) (any, error) { return strconv.ParseBool(s) }
	m["__4"] = nil

	if err := CacheNewEntries(m); err != nil {
		t.Errorf("CacheNewEntries() error = %v", err)
	}

	if got := Get[int]("__1"); !reflect.DeepEqual(got, 32) {
		t.Errorf("Get() = %v, want %v", got, 32)
	}
	if got := Get[bool]("__2"); !reflect.DeepEqual(got, true) {
		t.Errorf("Get() = %v, want %v", got, true)
	}
	if got := Get[bool]("__3"); !reflect.DeepEqual(got, false) {
		t.Errorf("Get() = %v, want %v", got, false)
	}
	if got := Get[string]("__4"); !reflect.DeepEqual(got, "abcde") {
		t.Errorf("Get() = %v, want %v", got, "abcde")
	}
}

func TestGetFromEnv(t *testing.T) {
	t.Setenv("__1", "32")
	t.Setenv("__2", "t")
	t.Setenv("__3", "false")
	t.Setenv("__4", "_")

	tests := []struct {
		key     string
		f       func(string) (any, error)
		wantErr bool
		want    any
	}{
		{"__1", func(s string) (any, error) { return strconv.Atoi(s) }, false, 32},
		{"__2", func(s string) (any, error) { return strconv.ParseBool(s) }, false, true},
		{"__3", func(s string) (any, error) { return strconv.ParseBool(s) }, false, false},
		{"__4", func(s string) (any, error) { return strconv.ParseBool(s) }, true, false},
		{"__5", func(s string) (any, error) { return s, nil }, true, ""},
	}

	for _, tt := range tests {
		got, err := GetFromEnv(tt.key, tt.f)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetFromEnv() error = %v, wantErr %v", err, tt.wantErr)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetFromEnv() = %v, want %v", got, tt.want)
		}
	}
}
