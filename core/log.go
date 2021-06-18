package core

import (
	"bytes"
	"fmt"
	"log"

	"github.com/mattn/go-colorable"
)

type (
	inner func(interface{}) string
)

var Log *log.Logger

func NewLogger() {
	Log = log.New(colorable.NewColorableStdout(), ColorDim("dev-cli: "), 0)
}

// Color styles
const (
	cb  = "30"
	cr  = "31"
	cg  = "32"
	cw  = "37"
	cgr = "90"

	bb = "40"
	br = "41"
	bg = "42"

	rr = "0"
	rb = "1"
	rd = "2"
	ri = "3"
	rs = "9"
)

var (
	ColorRed   = outer(cr)
	ColorGreen = outer(cg)
	ColorWhite = outer(cw)
	ColorGray  = outer(cgr)

	BgBlack = outer(bb)
	BgRed   = outer(br)
	BgGreen = outer(bg)

	ColorReset     = outer(rr)
	ColorBold      = outer(rb)
	ColorDim       = outer(rd)
	ColorItalic    = outer(ri)
	ColorStrikeout = outer(rs)
)

func outer(n string) inner {
	return func(msg interface{}) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[")
		b.WriteString(n)
		b.WriteString("m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), msg)
	}
}

func LogInfo(v interface{}) {
	Log.Println(ColorGray(v))
}

func LogError(v interface{}) {
	Log.Println(ColorRed(v))
}

func LogSuccess(v interface{}) {
	Log.Println(ColorGreen(v))
}
