package utils

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/kutogroup/kuto.api/pkg"
)

//300K，缓存
var pool = pkg.NewBytePool(300, 1024)

func CopyWithContext(ctx context.Context, src, dst net.Conn, addTraffic func(int)) (int, error) {
	defer dst.Close()

	buf := pool.Get()
	defer pool.Put(buf)

	count := 0
	for {
		if ctx.Err() == context.Canceled {
			break
		}

		src.SetDeadline(time.Now().Add(20 * time.Second))
		nr, er := src.Read(buf)
		if er != nil {
			if er == io.EOF {
				break
			}

			return 0, er
		}

		if nr > 0 {
			dst.SetDeadline(time.Now().Add(20 * time.Second))
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				return 0, ew
			}

			if nr != nw {
				return 0, io.ErrShortWrite
			}

			if addTraffic != nil {
				addTraffic(nr)
			}

			count = count + nw
		}
	}

	return count, nil
}