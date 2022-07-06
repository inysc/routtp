package routtp

import (
	"math"
	"net/http"
	"path"
	"unsafe"
)

// AbortIndex represents a typical value used in abort functions.
const AbortIndex int8 = math.MaxInt8 >> 1

type HandlerFunc = http.HandlerFunc
type HandlersChain []http.HandlerFunc

func assert(guard bool, text string) {
	if guard {
		panic(text)
	}
}

func JoinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func longestPrefix(a, b string) (idx int) {
	max := len(a)
	if max > len(b) {
		max = len(b)
	}
	idx = max
	for i := 0; i < max; i++ {
		if a[i] != b[i] {
			idx = i
			break
		}
	}

	return idx
}

func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func InvalidPath(path string) {
	assert(len(path) == 0, "must be at least a byte")
	assert(path[0] != '/', "must begin with '/'")

	for i := 1; i < len(path); i++ {

	}
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
