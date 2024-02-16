package util

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

type profiling struct {
	file     *os.File
	execTime time.Time
}

func NewProfiling(filename string) *profiling {
	p := profiling{
		execTime: time.Now(),
	}
	p.start(filename)
	return &p
}

func (m *profiling) start(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	m.file = f
}

func (m *profiling) Close() {
	if m.file == nil {
		return
	}
	defer m.file.Close()
	pprof.WriteHeapProfile(m.file)
	fmt.Println("exec time: ", time.Since(m.execTime))
}
