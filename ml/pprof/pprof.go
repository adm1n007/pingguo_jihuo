package pprof

import (
    . "ml/trace"
    "ml/os2"
    "runtime/pprof"
    // "runtime"
    "os"
)

func Start() {
    fn := os2.Executable()

    f, err := os.Create(fn + ".cpu.prof")
    RaiseIf(err)

    pprof.StartCPUProfile(f)
    // runtime.SetCPUProfileRate(1000000)
}

func Stop() {
    pprof.StopCPUProfile()
    f, err := os.Create(os2.Executable() + ".mem.prof")
    RaiseIf(err)
    pprof.WriteHeapProfile(f)
}
