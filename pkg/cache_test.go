package pkg_test

import (
	"testing"
	"time"

	"waha.api/pkg"
)

var cache = pkg.NewCache(":6379", 5, 300*time.Second)

func TestCache(t *testing.T) {
	err := cache.Set("test1", "123")
	if err != nil {
		t.Error("redis set t1 error=", err)
		return
	}
	v1, err := cache.Get("test1")
	if err != nil {
		t.Error("redis get v2 error=", err)
		return
	}
	if v1 != "123" {
		t.Error("redis failed")
		return
	}

	if !cache.Exists("test1") {
		t.Error("redis exist failed")
		return
	}

	err = cache.Del("test1")
	if err != nil || cache.Exists("test1") {
		t.Error("redis del failed, err=", err)
		return
	}

}
