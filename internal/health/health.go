package health

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Snapshot is a runtime health snapshot of the current process.
type Snapshot struct {
	Status     string      `json:"status"`
	Goroutines int         `json:"goroutines"`
	Memory     MemoryInfo  `json:"memory"`
	Runtime    RuntimeInfo `json:"runtime"`
	Timestamp  string      `json:"timestamp"`
}

// MemoryInfo contains memory statistics in MB.
type MemoryInfo struct {
	AllocMB      float64 `json:"allocMB"`
	TotalAllocMB float64 `json:"totalAllocMB"`
	SysMB        float64 `json:"sysMB"`
	NumGC        uint32  `json:"numGC"`
}

// RuntimeInfo contains Go runtime metadata.
type RuntimeInfo struct {
	Version string `json:"version"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	CPUs    int    `json:"cpus"`
}

// Collect returns a health snapshot for the current process.
func Collect() Snapshot {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return Snapshot{
		Status:     "healthy",
		Goroutines: runtime.NumGoroutine(),
		Memory: MemoryInfo{
			AllocMB:      float64(mem.Alloc) / 1024 / 1024,
			TotalAllocMB: float64(mem.TotalAlloc) / 1024 / 1024,
			SysMB:        float64(mem.Sys) / 1024 / 1024,
			NumGC:        mem.NumGC,
		},
		Runtime: RuntimeInfo{
			Version: runtime.Version(),
			OS:      runtime.GOOS,
			Arch:    runtime.GOARCH,
			CPUs:    runtime.NumCPU(),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// FormatText formats a snapshot into a human-readable text block.
func FormatText(s Snapshot) string {
	var b strings.Builder
	b.WriteString("nagobot Health\n")
	b.WriteString("==============\n\n")
	b.WriteString(fmt.Sprintf("Status: %s\n\n", s.Status))
	b.WriteString("Memory:\n")
	b.WriteString(fmt.Sprintf("  Allocated: %.2f MB\n", s.Memory.AllocMB))
	b.WriteString(fmt.Sprintf("  Total Allocated: %.2f MB\n", s.Memory.TotalAllocMB))
	b.WriteString(fmt.Sprintf("  System: %.2f MB\n", s.Memory.SysMB))
	b.WriteString(fmt.Sprintf("  GC Cycles: %d\n\n", s.Memory.NumGC))
	b.WriteString("Runtime:\n")
	b.WriteString(fmt.Sprintf("  Go Version: %s\n", s.Runtime.Version))
	b.WriteString(fmt.Sprintf("  OS/Arch: %s/%s\n", s.Runtime.OS, s.Runtime.Arch))
	b.WriteString(fmt.Sprintf("  CPUs: %d\n", s.Runtime.CPUs))
	b.WriteString(fmt.Sprintf("  Goroutines: %d\n", s.Goroutines))
	return b.String()
}
