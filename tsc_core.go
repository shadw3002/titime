package titime

import "time"

type TscLevel = int

const (
	TscStable = iota
	TscPerCpuStable
	TscUnstable
)

type TscCore struct {
	level            TscLevel
	cyclesPerSecond  uint64
	cyclesFromAnchor uint64
}

func newTscCore() (core TscCore) {
	anchor := time.Now()

	if isTscStable() {
		core.level = TscStable
		core.cyclesFromAnchor, core.cyclesFromAnchor = cyclesPerSec(anchor)
	} else if isTscPerCpuStable() {
		core.level = TscPerCpuStable
		cpus, err := availableCpus()
		if len(cpus) == 0 || err != nil {
			core.level = TscUnstable
			return
		}

		for _, cpu := range cpus {
			err := pinCPU(cpu)
			if err != nil {
				// TODO
			}

			err = unpinCPU()
			if err != nil {
				// TODO
			}
		}

	} else {
		core.level = TscUnstable
	}
}
