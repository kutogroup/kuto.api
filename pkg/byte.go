package pkg

//KutoBytePool Byte池
type KutoBytePool struct {
	c chan []byte
	w int
}

// NewKutoBytePool creates a new KutoBytePool bounded to the given maxSize, with new
// byte arrays sized based on width.
func NewKutoBytePool(maxSize int, width int) (bp *KutoBytePool) {
	return &KutoBytePool{
		c: make(chan []byte, maxSize),
		w: width,
	}
}

// Get gets a []byte from the KutoBytePool, or creates a new one if none are
// available in the pool.
func (bp *KutoBytePool) Get() (b []byte) {
	select {
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = make([]byte, bp.w)
	}
	return
}

// Put returns the given Buffer to the KutoBytePool.
func (bp *KutoBytePool) Put(b []byte) {
	if cap(b) < bp.w {
		// someone tried to put back a too small buffer, discard it
		return
	}

	select {
	case bp.c <- b[:bp.w]:
		// buffer went back into pool
	default:
		// buffer didn't go back into pool, just discard
	}
}

// Width returns the width of the byte arrays in this pool.
func (bp *KutoBytePool) Width() (n int) {
	return bp.w
}
