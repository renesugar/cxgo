package libs

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dennwc/cxgo/types"
)

const (
	stdintH = "stdint.h"
)

func init() {
	RegisterLibrary(stdintH, func(c *Env) *Library {
		return &Library{
			Header: incStdInt(c.Env),
			Types:  typesStdInt(c.Env),
		}
	})
}

func fixedIntTypeDefs(buf *strings.Builder, part string) {
	for _, unsigned := range []bool{false, true} {
		name := "int"
		if unsigned {
			name = "uint"
		}
		for _, sz := range intSizes {
			_, _ = fmt.Fprintf(buf, "#define %s%s%d_t %s\n",
				name, part, sz,
				buildinFixedIntName(sz, unsigned),
			)
		}
		buf.WriteByte('\n')
	}
}

func fixedIntTypes() string {
	var buf strings.Builder
	fixedIntTypeDefs(&buf, "")
	fixedIntTypeDefs(&buf, "_least")
	fixedIntTypeDefs(&buf, "_fast")
	return buf.String()
}

func maxIntTypeDefs(buf *strings.Builder, part string, sz int) {
	for _, unsigned := range []bool{false, true} {
		name := "int"
		if unsigned {
			name = "uint"
		}
		_, _ = fmt.Fprintf(buf, "typedef %s %s%s_t;\n",
			buildinFixedIntName(sz, unsigned),
			name, part,
		)
	}
	buf.WriteByte('\n')
}

func maxIntTypes(e *types.Env) string {
	var buf strings.Builder
	maxIntTypeDefs(&buf, "ptr", e.PtrSize()*8)
	maxIntTypeDefs(&buf, "max", intSizes[len(intSizes)-1])
	return buf.String()
}

func intMinMaxDef(buf *strings.Builder, name string, min, max int64) {
	_, _ = fmt.Fprintf(buf, "#define %s_MIN %#x\n", name, min)
	_, _ = fmt.Fprintf(buf, "#define %s_MAX %#x\n", name, max)
}

func uintMaxDef(buf *strings.Builder, name string, max uint64) {
	_, _ = fmt.Fprintf(buf, "#define %s_MAX %#x\n", name, max)
}

func intLimitsDef(buf *strings.Builder, part string, min, max int64, umax uint64) {
	intMinMaxDef(buf, "INT"+part, min, max)
	uintMaxDef(buf, "UINT"+part, umax)
}

func intLimitsDefs(buf *strings.Builder, part string) {
	for i, sz := range intSizes {
		intLimitsDef(buf, part+strconv.Itoa(sz), minInts[i], maxInts[i], maxUints[i])
	}
	buf.WriteByte('\n')
}

func intLimits() string {
	var buf strings.Builder
	intLimitsDefs(&buf, "")
	intLimitsDefs(&buf, "_LEAST")
	intLimitsDefs(&buf, "_FAST")
	return buf.String()
}

func intSizeInd(sz int) int {
	switch sz {
	case 0, 1:
		return 0
	case 2:
		return 1
	case 4:
		return 2
	case 8:
		return 3
	default:
		panic(sz)
	}
}

func maxIntLimits(e *types.Env) string {
	var buf strings.Builder
	i := intSizeInd(e.PtrSize())
	intLimitsDef(&buf, "PTR", minInts[i], maxInts[i], maxUints[i])
	buf.WriteByte('\n')

	i = len(intSizes) - 1
	intLimitsDef(&buf, "MAX", minInts[i], maxInts[i], maxUints[i])
	buf.WriteByte('\n')
	return buf.String()
}

func otherLimits(e *types.Env) string {
	var buf strings.Builder

	i := intSizeInd(e.PtrSize())
	intMinMaxDef(&buf, "PTRDIFF", minInts[i], maxInts[i])
	buf.WriteByte('\n')

	// TODO: SIG_ATOMIC

	i = intSizeInd(e.PtrSize())
	uintMaxDef(&buf, "SIZE", maxUints[i])
	buf.WriteByte('\n')

	i = intSizeInd(e.C().WCharSize())
	if e.C().WCharSigned() {
		intMinMaxDef(&buf, "WCHAR", minInts[i], maxInts[i])
	} else {
		intMinMaxDef(&buf, "WCHAR", 0, int64(maxUints[i]))
	}
	buf.WriteByte('\n')

	i = intSizeInd(e.C().WIntSize())
	if e.C().WCharSigned() {
		intMinMaxDef(&buf, "WINT", minInts[i], maxInts[i])
	} else {
		intMinMaxDef(&buf, "WINT", 0, int64(maxUints[i]))
	}
	buf.WriteByte('\n')
	return buf.String()
}

var (
	intSizes = []int{8, 16, 32, 64}
	minInts  = []int64{math.MinInt8, math.MinInt16, math.MinInt32, math.MinInt64}
	maxInts  = []int64{math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
	maxUints = []uint64{math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}
)

func incStdInt(e *types.Env) string {
	var buf strings.Builder
	buf.WriteString("#include <" + BuiltinH + ">\n")
	buf.WriteString(fixedIntTypes())
	buf.WriteString(maxIntTypes(e))
	buf.WriteString(intLimits())
	buf.WriteString(maxIntLimits(e))
	buf.WriteString(otherLimits(e))
	return buf.String()
}

func typesStdInt(e *types.Env) map[string]types.Type {
	m := make(map[string]types.Type, 2*len(intSizes))
	for _, unsigned := range []bool{false, true} {
		for _, sz := range intSizes {
			name := buildinFixedIntName(sz, unsigned)
			if unsigned {
				m[name] = types.UintT(sz / 8)
			} else {
				m[name] = types.IntT(sz / 8)
			}
		}
	}
	m["intptr_t"] = e.IntPtrT()
	m["uintptr_t"] = e.UintPtrT()
	max := intSizes[len(intSizes)-1] / 8
	m["intmax_t"] = types.IntT(max)
	m["uintmax_t"] = types.UintT(max)
	return m
}
