package libs

// https://pubs.opengroup.org/onlinepubs/009695399/basedefs/float.h.html

const (
	floatH = "float.h"
)

func init() {
	RegisterLibrary(floatH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#define FLT_RADIX 2
#define DECIMAL_DIG 10
#define FLT_DIG 6
#define DBL_DIG 10
#define LDBL_DIG 10
#define FLT_MIN_10_EXP -37
#define DBL_MIN_10_EXP -37
#define LDBL_MIN_10_EXP -37
#define FLT_MAX_10_EXP +37
#define DBL_MAX_10_EXP +37
#define LDBL_MAX_10_EXP +37
#define FLT_MAX 1E+37
#define DBL_MAX 1E+37
#define LDBL_MAX 1E+37
#define FLT_EPSILON 1E-5
#define DBL_EPSILON 1E-9
#define LDBL_EPSILON 1E-9
#define FLT_MIN 1E-37
#define DBL_MIN 1E-37
#define LDBL_MIN 1E-37
`,
		}
	})
}
