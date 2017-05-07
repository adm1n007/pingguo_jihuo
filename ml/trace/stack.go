package trace

import (
    "fmt"
    "runtime"
    "bytes"
    "strings"
    "io/ioutil"
)

var (
    dunno           = []byte("???")
    centerDot       = []byte("·")
    dot             = []byte(".")
    slash           = []byte("/")
    mainFuncName    = []byte("main")
    goexitFuncName  = []byte("goexit")
)

func source(lines [][]byte, n int) []byte {
    if n < 0 || n >= len(lines) {
        return dunno
    }
    return bytes.Trim(lines[n], " \t")
}

func function(n string) []byte {
    name := []byte(n)

    if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
        name = name[lastslash+1:]
    }

    if period := bytes.Index(name, dot); period >= 0 {
        name = name[period+1:]
    }

    return bytes.Replace(name, centerDot, dot, -1)
}

func stack(skip, depth int) []byte {
    entries := make([]callStackEntry, depth)
    if getCallStack(skip + 1, entries) == 0 {
        return nil
    }

    return []byte(strings.Join(getTrackBack(entries), "\r\n"))
}

func getCallStack(skip int, entries []callStackEntry) int {
    depth := len(entries) + skip
    for i := skip; i < depth; i++ {
        pc, file, line, ok := runtime.Caller(i)
        if ok == false {
            continue
        }
        entries[i - skip] = callStackEntry{
            pc: pc,
            file: file,
            line: line,
        }
    }

    length := len(entries)
    for i := 0; i != length / 2; i++ {
        entries[i], entries[length - i - 1] = entries[length - i - 1], entries[i]
    }

    return length
}

func getTrackBack(entries []callStackEntry) []string {
    var lastFile string
    var traceBack []string
    var lines [][]byte

    mainFound := false

    for _, pc := range entries {
        fn := runtime.FuncForPC(pc.pc)
        if fn == nil {
            continue
        }

        funcName := function(fn.Name())

        if mainFound == false && bytes.Equal(funcName, mainFuncName) {
            mainFound = true
            traceBack = nil
            continue
        }

        file, line := pc.file, pc.line
        traceBack = append(traceBack, fmt.Sprintf("%s:%d (0x%x)", file, line, pc.pc))

        detail := fmt.Sprintf("\t%s", funcName)

        if Config.ReadSource {
            if file != lastFile {
                lastFile = file
                data, err := ioutil.ReadFile(file)
                lines = nil

                if err == nil {
                    lines = bytes.Split(data, []byte{'\n'})
                }
            }

            if lines != nil {
                detail = fmt.Sprintf("\t%s: %s", funcName, source(lines, line - 1))
            }
        }

        traceBack = append(traceBack, detail)
    }

    traceBack = append(traceBack, "")

    return traceBack
}
