package gollect

import "strings"

// A list of Golang's builtin package paths except internal prefixed packages.
// Packages differ depending on the environment, but we ignore it
// because they are not so important in competition programming.
//
// Writing directly for better performance
// generated by cmd/list-builtin-packages
var builtinPackages = map[string]interface{}{
	"archive/tar":          struct{}{},
	"archive/zip":          struct{}{},
	"bufio":                struct{}{},
	"bytes":                struct{}{},
	"compress/bzip2":       struct{}{},
	"compress/flate":       struct{}{},
	"compress/gzip":        struct{}{},
	"compress/lzw":         struct{}{},
	"compress/zlib":        struct{}{},
	"container/heap":       struct{}{},
	"container/list":       struct{}{},
	"container/ring":       struct{}{},
	"context":              struct{}{},
	"crypto":               struct{}{},
	"crypto/aes":           struct{}{},
	"crypto/cipher":        struct{}{},
	"crypto/des":           struct{}{},
	"crypto/dsa":           struct{}{},
	"crypto/ecdh":          struct{}{},
	"crypto/ecdsa":         struct{}{},
	"crypto/ed25519":       struct{}{},
	"crypto/elliptic":      struct{}{},
	"crypto/hmac":          struct{}{},
	"crypto/md5":           struct{}{},
	"crypto/rand":          struct{}{},
	"crypto/rc4":           struct{}{},
	"crypto/rsa":           struct{}{},
	"crypto/sha1":          struct{}{},
	"crypto/sha256":        struct{}{},
	"crypto/sha512":        struct{}{},
	"crypto/subtle":        struct{}{},
	"crypto/tls":           struct{}{},
	"crypto/x509":          struct{}{},
	"crypto/x509/pkix":     struct{}{},
	"database/sql":         struct{}{},
	"database/sql/driver":  struct{}{},
	"debug/buildinfo":      struct{}{},
	"debug/dwarf":          struct{}{},
	"debug/elf":            struct{}{},
	"debug/gosym":          struct{}{},
	"debug/macho":          struct{}{},
	"debug/pe":             struct{}{},
	"debug/plan9obj":       struct{}{},
	"embed":                struct{}{},
	"encoding":             struct{}{},
	"encoding/ascii85":     struct{}{},
	"encoding/asn1":        struct{}{},
	"encoding/base32":      struct{}{},
	"encoding/base64":      struct{}{},
	"encoding/binary":      struct{}{},
	"encoding/csv":         struct{}{},
	"encoding/gob":         struct{}{},
	"encoding/hex":         struct{}{},
	"encoding/json":        struct{}{},
	"encoding/pem":         struct{}{},
	"encoding/xml":         struct{}{},
	"errors":               struct{}{},
	"expvar":               struct{}{},
	"flag":                 struct{}{},
	"fmt":                  struct{}{},
	"go/ast":               struct{}{},
	"go/build":             struct{}{},
	"go/build/constraint":  struct{}{},
	"go/constant":          struct{}{},
	"go/doc":               struct{}{},
	"go/doc/comment":       struct{}{},
	"go/format":            struct{}{},
	"go/importer":          struct{}{},
	"go/parser":            struct{}{},
	"go/printer":           struct{}{},
	"go/scanner":           struct{}{},
	"go/token":             struct{}{},
	"go/types":             struct{}{},
	"hash":                 struct{}{},
	"hash/adler32":         struct{}{},
	"hash/crc32":           struct{}{},
	"hash/crc64":           struct{}{},
	"hash/fnv":             struct{}{},
	"hash/maphash":         struct{}{},
	"html":                 struct{}{},
	"html/template":        struct{}{},
	"image":                struct{}{},
	"image/color":          struct{}{},
	"image/color/palette":  struct{}{},
	"image/draw":           struct{}{},
	"image/gif":            struct{}{},
	"image/jpeg":           struct{}{},
	"image/png":            struct{}{},
	"index/suffixarray":    struct{}{},
	"io":                   struct{}{},
	"io/fs":                struct{}{},
	"io/ioutil":            struct{}{},
	"log":                  struct{}{},
	"log/syslog":           struct{}{},
	"math":                 struct{}{},
	"math/big":             struct{}{},
	"math/bits":            struct{}{},
	"math/cmplx":           struct{}{},
	"math/rand":            struct{}{},
	"mime":                 struct{}{},
	"mime/multipart":       struct{}{},
	"mime/quotedprintable": struct{}{},
	"net":                  struct{}{},
	"net/http":             struct{}{},
	"net/http/cgi":         struct{}{},
	"net/http/cookiejar":   struct{}{},
	"net/http/fcgi":        struct{}{},
	"net/http/httptest":    struct{}{},
	"net/http/httptrace":   struct{}{},
	"net/http/httputil":    struct{}{},
	"net/http/internal":    struct{}{},
	"net/http/pprof":       struct{}{},
	"net/mail":             struct{}{},
	"net/netip":            struct{}{},
	"net/rpc":              struct{}{},
	"net/rpc/jsonrpc":      struct{}{},
	"net/smtp":             struct{}{},
	"net/textproto":        struct{}{},
	"net/url":              struct{}{},
	"os":                   struct{}{},
	"os/exec":              struct{}{},
	"os/signal":            struct{}{},
	"os/user":              struct{}{},
	"path":                 struct{}{},
	"path/filepath":        struct{}{},
	"plugin":               struct{}{},
	"reflect":              struct{}{},
	"regexp":               struct{}{},
	"regexp/syntax":        struct{}{},
	"runtime":              struct{}{},
	"runtime/cgo":          struct{}{},
	"runtime/coverage":     struct{}{},
	"runtime/debug":        struct{}{},
	"runtime/metrics":      struct{}{},
	"runtime/pprof":        struct{}{},
	"runtime/race":         struct{}{},
	"runtime/trace":        struct{}{},
	"sort":                 struct{}{},
	"strconv":              struct{}{},
	"strings":              struct{}{},
	"sync":                 struct{}{},
	"sync/atomic":          struct{}{},
	"syscall":              struct{}{},
	"testing":              struct{}{},
	"testing/fstest":       struct{}{},
	"testing/iotest":       struct{}{},
	"testing/quick":        struct{}{},
	"text/scanner":         struct{}{},
	"text/tabwriter":       struct{}{},
	"text/template":        struct{}{},
	"text/template/parse":  struct{}{},
	"time":                 struct{}{},
	"time/tzdata":          struct{}{},
	"unicode":              struct{}{},
	"unicode/utf16":        struct{}{},
	"unicode/utf8":         struct{}{},
	"unsafe":               struct{}{},
}

// package path prefixes treated as same as buitin packages.
var thirdPartyPackagePathPrefixes []string

// !! this function populates global variable !!
func setThirdPartyPackagePathPrefixes(s []string) {
	thirdPartyPackagePathPrefixes = s
}

func isBuiltinPackage(path string) bool {
	_, ok := builtinPackages[path]
	for i := 0; !ok && i < len(thirdPartyPackagePathPrefixes); i++ {
		ok = strings.HasPrefix(path, thirdPartyPackagePathPrefixes[i])
	}
	return ok
}
