package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	inf = "[INFO ] "
	war = "[WARN ] "
	err = "[ERROR] "
	fat = "[FATAL] "
	que = "[QUERY] "
	maxFileLength = 25
	maxFunctionLength = 20
	// length of padding needs to be >= max(maxFileLength, maxFunctionLength)
	// enables use of fancy slicing, too lazy to use array...
	padding = "                         "
)

type Logger struct {
	mut sync.Mutex
	out io.Writer
	buf []byte
}

var mut = &sync.Mutex{}
var outLogger = &Logger{out: os.Stdout, mut: *mut}
var errLogger = &Logger{out: os.Stderr, mut: *mut}

func Info(v ...interface{}) {
	outLogger.output(inf, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	outLogger.output(inf, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	outLogger.output(war, fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	outLogger.output(war, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	errLogger.output(err, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	errLogger.output(err, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	errLogger.output(fat, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	errLogger.output(fat, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Query(v ...interface{}) {
	outLogger.output(que, fmt.Sprint(v...))
}

func (l *Logger) output(level string, content string) {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.formatHeader(level)
	l.buf = append(l.buf, content...)
	if len(content) == 0 || content[len(content)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	l.out.Write(l.buf)
}

func (l *Logger) formatHeader(level string) {
	l.buf = l.buf[:0]
	l.addTimeStamp()
	l.buf = append(l.buf, level...)
	l.addCaller()
}

func (l *Logger) addTimeStamp() {
	t := time.Now()
	t = t.UTC()
	year, month, day := t.Date()
	itoa(&l.buf, year, 4)
	l.buf = append(l.buf, '/')
	itoa(&l.buf, int(month), 2)
	l.buf = append(l.buf, '/')
	itoa(&l.buf, day, 2)
	l.buf = append(l.buf, ' ')
	hour, min, sec := t.Clock()
	itoa(&l.buf, hour, 2)
	l.buf = append(l.buf, ':')
	itoa(&l.buf, min, 2)
	l.buf = append(l.buf, ':')
	itoa(&l.buf, sec, 2)
	l.buf = append(l.buf, ' ')
}

// taken from native log package
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) addCaller() {
	var pc uintptr
	var file string
	var function string
	var line int
	var ok bool
	pc, file, line, ok = runtime.Caller(4)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		file = "???"
		function = "???"
		line = 0
	} else {
		function = details.Name()
		function = getRemainingAfter('.', function)
		file = getRemainingAfter('/', file)
	}
	l.buf = append(l.buf, file...)
	l.buf = append(l.buf, ':')
	itoa(&l.buf, line, 3)
	l.buf = append(l.buf, ':')
	if len(file)+4 < maxFileLength {
		l.buf = append(l.buf, padding[len(file)+3:maxFileLength-1]...)
	}
	l.buf = append(l.buf, function...)
	l.buf = append(l.buf, ':')
	if len(function)+1 < maxFunctionLength {
		l.buf = append(l.buf, padding[len(function):maxFunctionLength-1]...)
	}
}

func getRemainingAfter(c byte, s string) string {
	short := s
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == c {
			short = s[i+1:]
			break
		}
	}
	return short
}
