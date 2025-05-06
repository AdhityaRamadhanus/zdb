package zdb

func Mask64(x uint) uint64 {
	bits := 0
	for ; x > 0; x /= 2 {
		bits += 1
	}

	return uint64(1 << bits)
}
