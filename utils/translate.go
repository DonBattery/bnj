package utils

var objSizes = map[string]struct {
	w int
	h int
}{
	"bunny": {
		w: 16,
		h: 16,
	},
}

func SizeFromObjType(objType string) (width, height int) {
	if val, ok := objSizes[objType]; ok {
		return val.w, val.h
	}
	return 0, 0
}
