package c

/*
int sshdmain(const char *b, const char *e);
*/
import "C"

import "unsafe"

func Match(data []byte) int {
	b := unsafe.Pointer(&data[0])
	e := unsafe.Add(b, len(data))
	return int(C.sshdmain((*C.char)(b), (*C.char)(e)))
}
