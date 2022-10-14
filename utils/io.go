package utils

import (
	"context"
	"io"
)

func CopyWithContext(ctx context.Context, dst io.WriteCloser, src io.ReadCloser, buf []byte, addTraffic func(int)) (int, error) {
	defer dst.Close()
	defer src.Close()

	count := 0
	for {
		if ctx.Err() == context.Canceled {
			break
		}

		nr, er := src.Read(buf)
		if er != nil {
			if er == io.EOF {
				break
			}

			return 0, er
		}

		if nr > 0 {
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