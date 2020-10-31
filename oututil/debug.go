package oututil

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	l         = log.New(os.Stderr, "", 0)
	debugging = os.Getenv("DEBUG") == "TRUE"
	mut       = &sync.Mutex{}
)

func printheader() {
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return
	}

	_, name := path.Split(fun.Name())
	_, callerFile, line, _ := runtime.Caller(2)
	x, _ := filepath.EvalSymlinks(fmt.Sprintf("/proc/%d/cwd", os.Getppid()))

	rel, _ := filepath.Rel(x, callerFile)
	l.Printf("%s:%d > %s\n", rel, line, name)
}

// Printf prints a message with a header
func Printf(format string, v ...interface{}) {
	if debugging {
		mut.Lock()
		printheader()
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, format, v...)
		Indent(&buf, 1)
		fmt.Print(buf.String())
		mut.Unlock()
	}
}

// Println prints a message with a header and appends a \n
func Println(v ...interface{}) {
	if debugging {
		mut.Lock()
		printheader()
		buf := bytes.Buffer{}
		fmt.Fprintln(&buf, v...)
		Indent(&buf, 1)
		fmt.Print(buf.String())
		mut.Unlock()
	}
}

// PrefixLine prefixes all lines in a buffer
func PrefixLine(buf *bytes.Buffer, prefix string, level int) {
	tmp := &bytes.Buffer{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		fmt.Fprintln(tmp, strings.Repeat(prefix, level), scanner.Text())
	}
	buf.Reset()
	buf.Write(tmp.Bytes())
}

// Indent prefixes all lines with \t
func Indent(buf *bytes.Buffer, level int) { PrefixLine(buf, "\t", level) }
