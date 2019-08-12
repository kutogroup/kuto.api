package pkg_test

import (
	"os"
	"testing"
	"time"

	"waha.api/pkg"
)

var cdn = pkg.NewCDN("cdn.kutoapps.com", "1SYOD1M9RIAS97U1X1RZ", "anBa2vd6jZtIgNZ0QNJaYAMKvhE3Lp7gQRQLFt3H", time.Minute)

func TestCDN(t *testing.T) {
	f, err := os.Open("cdn.go")
	if err != nil {
		t.Error("open file error", err)
	}

	stat, _ := f.Stat()
	err = cdn.Put("image", "test.txt", f, stat.Size())
	if err != nil {
		t.Error("cdn upload failed, err=", err)
	}
}
