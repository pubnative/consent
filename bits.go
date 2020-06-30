package consent

type bitWriter struct {
	data []byte
	cur  byte
	used uint
}

func (b *bitWriter) AppendByte(add byte, size uint) {
	if size+b.used > 8 {
		b.cur |= add << (8 - size) >> b.used
		b.data = append(b.data, b.cur)
		size -= 8 - b.used
		b.cur = 0
		b.used = 0
	}

	b.cur |= add << (8 - size) >> b.used
	b.used += size
}

func (b *bitWriter) AppendBools(bools []bool) {
	for i := range bools {
		var v byte
		if bools[i] {
			v = 1
		}
		b.AppendByte(v, 1)
	}
}

func (b *bitWriter) AppendBit(x bool) {
	var v byte
	if x {
		v = 1
	}
	b.AppendByte(v, 1)
}

func (b *bitWriter) AppendBits(bits bitWriter) {
	for i := range bits.data {
		b.AppendByte(bits.data[i], 8)
	}

	b.AppendByte(bits.cur>>(8-bits.used), bits.used)
}

func (b *bitWriter) AppendInt(add int64, size uint) {
	for size > 0 {
		byteSize := size
		if byteSize > 8 {
			byteSize = 8
		}
		size -= byteSize
		b.AppendByte(byte(add>>size), byteSize)
	}
}

func (b *bitWriter) Bytes() []byte {
	if b.used == 0 {
		return b.data
	}

	return append(b.data, b.cur)
}

type bitReader struct {
	cur  byte
	left uint
	data []byte
}

func newBitReader(buf []byte) bitReader {
	return bitReader{
		data: buf,
	}
}

func (b *bitReader) ReadByte(size uint) (byte, bool) {
	if size <= b.left {
		out := b.cur >> (8 - size)
		b.cur = b.cur << size
		b.left -= size
		return out, true
	}

	if len(b.data) == 0 {
		return 0, false
	}

	out := b.cur >> (8 - size)
	size -= b.left

	out |= b.data[0] >> (8 - size)
	b.cur = b.data[0] << size
	b.data = b.data[1:]
	b.left = 8 - size
	return out, true
}

func (b *bitReader) ReadInt(size uint) (int64, bool) {
	var out int64
	for size > 0 {
		bitSize := size
		if bitSize > 8 {
			bitSize = 8
		}
		size -= bitSize
		v, ok := b.ReadByte(bitSize)
		if !ok {
			return 0, false
		}
		out = out<<bitSize | int64(v)
	}
	return out, true
}

func (b *bitReader) ReadBit() (bool, bool) {
	v, ok := b.ReadByte(1)
	return v > 0, ok
}
