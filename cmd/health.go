package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show system health information",
	Long:  `Display current system health including memory usage, goroutines, and uptime.`,
	RunE:  runHealth,
}

var healthJSON bool

func init() {
	healthCmd.Flags().BoolVar(&healthJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(healthCmd)
}

func runHealth(cmd *cobra.Command, args []string) error {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	health := map[string]any{
		"status":     "healthy",
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]any{
			"allocMB":      float64(mem.Alloc) / 1024 / 1024,
			"totalAllocMB": float64(mem.TotalAlloc) / 1024 / 1024,
			"sysMB":        float64(mem.Sys) / 1024 / 1024,
			"numGC":        mem.NumGC,
		},
		"runtime": map[string]any{
			"version":   runtime.Version(),
			"os":        runtime.GOOS,
			"arch":      runtime.GOARCH,
			"cpus":      runtime.NumCPU(),
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if healthJSON {
		data, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	// Human-readable output
	fmt.Println("nagobot Health")
	fmt.Println("==============")
	fmt.Println()
	fmt.Printf("Status: %s\n", health["status"])
	fmt.Println()
	fmt.Println("Memory:")
	memInfo := health["memory"].(map[string]any)
	fmt.Printf("  Allocated: %.2f MB\n", memInfo["allocMB"])
	fmt.Printf("  Total Allocated: %.2f MB\n", memInfo["totalAllocMB"])
	fmt.Printf("  System: %.2f MB\n", memInfo["sysMB"])
	fmt.Printf("  GC Cycles: %d\n", memInfo["numGC"])
	fmt.Println()
	fmt.Println("Runtime:")
	rtInfo := health["runtime"].(map[string]any)
	fmt.Printf("  Go Version: %s\n", rtInfo["version"])
	fmt.Printf("  OS/Arch: %s/%s\n", rtInfo["os"], rtInfo["arch"])
	fmt.Printf("  CPUs: %d\n", rtInfo["cpus"])
	fmt.Printf("  Goroutines: %d\n", health["goroutines"])

	return nil
}
