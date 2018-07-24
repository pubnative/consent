package consent

import (
	"strconv"
	"testing"
)

func TestBitsAppendByte(t *testing.T) {
	var b bitWriter
	if len(b.Bytes()) > 0 {
		t.Fail()
	}

	b.AppendByte(1, 1)
	equal(t, "10000000", toBitsString(b.Bytes()[0]))
	// ^ big endian, adding stuff from left to right

	b.AppendByte(0, 1)
	equal(t, "10000000", toBitsString(b.Bytes()[0]))
	b.AppendByte(9, 2)
	equal(t, "10010000", toBitsString(b.Bytes()[0]))
	b.AppendByte(127, 3)
	equal(t, "10011110", toBitsString(b.Bytes()[0]))
	b.AppendByte(127, 3)
	equal(t, "10011111", toBitsString(b.Bytes()[0]))
	equal(t, "11000000", toBitsString(b.Bytes()[1]))
	b.AppendByte(255, 8)
	equal(t, "10011111", toBitsString(b.Bytes()[0]))
	equal(t, "11111111", toBitsString(b.Bytes()[1]))
	equal(t, "11000000", toBitsString(b.Bytes()[2]))
}

func TestBitsAppendUInt(t *testing.T) {
	var b bitWriter
	b.AppendInt(1, 1)
	equal(t, "10000000", toBitsString(b.Bytes()[0]))
	b.AppendInt(256, 9)
	equal(t, "11000000", toBitsString(b.Bytes()[0]))
	equal(t, "0", toBitsString(b.Bytes()[1]))
	b.AppendInt(1, 1)
	equal(t, "11000000", toBitsString(b.Bytes()[0]))
	equal(t, "100000", toBitsString(b.Bytes()[1]))
}

func TestBitsReadByte(t *testing.T) {
	var b bitReader
	b.data = []byte{
		fromBitsString("11100000"),
		fromBitsString("11100000"),
		fromBitsString("11100001"),
		fromBitsString("11100000"),
	}
	v, ok := b.ReadByte(1)
	if !ok {
		t.Fail()
	}
	equal(t, byte(1), v)

	v, ok = b.ReadByte(2)
	if !ok {
		t.Fail()
	}
	equal(t, byte(3), v)

	v, ok = b.ReadByte(1)
	if !ok {
		t.Fail()
	}
	equal(t, byte(0), v)

	v, ok = b.ReadByte(8)
	if !ok {
		t.Fail()
	}
	equal(t, "1110", toBitsString(v))

	v, ok = b.ReadByte(4)
	if !ok {
		t.Fail()
	}
	equal(t, "0", toBitsString(v))

	v, ok = b.ReadByte(7)
	if !ok {
		t.Fail()
	}
	equal(t, "1110000", toBitsString(v))

	v, ok = b.ReadByte(8)
	if !ok {
		t.Fail()
	}
	equal(t, "11110000", toBitsString(v))

	v, ok = b.ReadByte(1)
	if !ok {
		t.Fail()
	}
	equal(t, "0", toBitsString(v))

	if uint(0) != b.left {
		t.Fail()
	}
}

func TestBitsReadInt(t *testing.T) {
	var b bitReader
	b.data = []byte{
		fromBitsString("11100000"),
		fromBitsString("11100000"),
		fromBitsString("11100000"),
	}
	v, ok := b.ReadInt(1)
	if !ok {
		t.Fail()
	}
	if v != 1 {
		t.Fail()
	}

	v, ok = b.ReadInt(16)
	if !ok {
		t.Fail()
	}
	equal(t, "1100000111000001", strconv.FormatInt(v, 2))

}

func toBitsString(b byte) string {
	return strconv.FormatUint(uint64(b), 2)
}

func fromBitsString(s string) byte {
	i, _ := strconv.ParseUint(s, 2, 8)
	return byte(i)
}
