package rpenc

type RPEnc struct {
	Shift1 byte
	Shift2 byte
	Table  [256]byte
}

func (e *RPEnc) Init(key []byte) {
	for i := 0; i < 256; i++ {
		e.Table[i] = byte(i)
	}
	var shift byte
	for i := 0; i < 256; i++ {
		a := key[i%len(key)]
		shift += (a + e.Table[i]) & 255

		b := e.Table[i]
		e.Table[i] = e.Table[shift]
		e.Table[shift] = b
	}
	return
}

func (e *RPEnc) Get(data []byte) []byte {
	res := data
	for i := 0; i < len(data); i++ {
		e.Shift1++
		a := e.Table[e.Shift1]

		e.Shift2 += a
		b := e.Table[e.Shift2]

		e.Table[e.Shift2] = a
		e.Table[e.Shift1] = b

		c := (a + b) & 255
		d := e.Table[c]

		res[i] ^= d
	}
	return res
}
