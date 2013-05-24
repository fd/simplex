package collate

import (
	"bytes"
	"testing"
)

func TestNil(t *testing.T) {
	var (
		buf Buffer
		col = New("en")
	)

	equals(t, must(col.Key(&buf, nil)), []byte{0x01})
}

func TestFalse(t *testing.T) {
	var (
		buf Buffer
		col = New("en")
	)

	equals(t, must(col.Key(&buf, false)), []byte{0x02})
}

func TestTrue(t *testing.T) {
	var (
		buf Buffer
		col = New("en")
	)

	equals(t, must(col.Key(&buf, true)), []byte{0x03})
}

func TestInt(t *testing.T) {
	var (
		buf Buffer
		col = New("en")
	)

	equals(t, must(col.Key(&buf, -6)), []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFA})

	buf.Reset()
	equals(t, must(col.Key(&buf, 6)), []byte{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06})

	buf.Reset()
	equals(t, must(col.Key(&buf, uint(6))), []byte{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06})
}

func TestString(t *testing.T) {
	var (
		buf Buffer
		col = New("en")
	)

	buf.Reset()
	equals(t, must(col.Key(&buf, "hello")), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

	buf.Reset()
	equals(t, must(col.Key(&buf, String("allo"))), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

	buf.Reset()
	equals(t, must(col.Key(&buf, String("Allo"))), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

	buf.Reset()
	equals(t, must(col.Key(&buf, String("âllo"))), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

	buf.Reset()
	equals(t, must(col.Key(&buf, String("Âllo"))), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

	buf.Reset()
	equals(t, must(col.Key(&buf, String("bike"))), []byte{0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00})

}

func equals(t *testing.T, a, b []byte) {
	if bytes.Compare(a, b) != 0 {
		t.Errorf("not equal: %x != %x", a, b)
	}
}

func must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
