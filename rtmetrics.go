// Package rtmetrics provide a easy way to send go's runtime data to
// the monitoring service of your choice (graphite, librato, whatever you want).
package rtmetrics

import (
	"runtime"
	"strings"
	"time"
)

// Defaul interval for
var Interval = 1 * time.Second

// Run function is to called with your own implementation of Collector. prefix should not end with a "."
func Run(c Collector, prefix string) {
	m := &runtime.MemStats{}
	ticker := time.Tick(Interval)
	prefix = strings.TrimRight(prefix, ".")
	for _ = range ticker {
		runtime.ReadMemStats(m)
		c.Measure(prefix+".gomaxprocs", uint64(runtime.GOMAXPROCS(0)))
		c.Measure(prefix+".numcpu", uint64(runtime.NumCPU()))
		c.Measure(prefix+".numcgocalls", uint64(runtime.NumCgoCall()))
		c.Measure(prefix+".numgoroutine", uint64(runtime.NumGoroutine()))
		c.Measure(prefix+".memory.alloc", m.Alloc)
		c.Measure(prefix+".memory.totalalloc", m.TotalAlloc)
		c.Measure(prefix+".memory.sys", m.Sys)
		c.Measure(prefix+".memory.lookups", m.Lookups)
		c.Measure(prefix+".memory.mallocs", m.Mallocs)
		c.Measure(prefix+".memory.frees", m.Frees)
		c.Measure(prefix+".memory.heap.alloc", m.HeapAlloc)
		c.Measure(prefix+".memory.heap.sys", m.HeapSys)
		c.Measure(prefix+".memory.heap.idle", m.HeapIdle)
		c.Measure(prefix+".memory.heap.inuse", m.HeapInuse)
		c.Measure(prefix+".memory.heap.released", m.HeapReleased)
		c.Measure(prefix+".memory.heap.objects", m.HeapObjects)
		c.Measure(prefix+".gc.pausetotal", m.PauseTotalNs)
		c.Measure(prefix+".gc.pause", m.PauseNs[(m.NumGC+255)%256])
		c.Measure(prefix+".gc.num", uint64(m.NumGC))
		c.Flush()
	}
}

// Collector is the interface one of your structure should implement in order for Run to be able to send data
type Collector interface {
	Measure(name string, value uint64)
	// Flush is called at the end of the collection cycle.
	Flush()
}
