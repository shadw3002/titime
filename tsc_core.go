package titime

import "time"

type TscLevel = int

const (
	TscStable = iota
	TscPerCpuStable
	TscUnstable
)

type TscCore struct {
	level             TscLevel
	cyclesPerSecond   uint64
	cyclesFromAnchors []uint64
	nanosPerCycles    float64
	anchor            time.Time
}

func NewTscCore() (core *TscCore) {
	core = &TscCore{}
	core.anchor = time.Now()

	if isTscStable() {
		core.cyclesFromAnchors = make([]uint64, 1)
		core.level = TscStable
		core.cyclesPerSecond, core.cyclesFromAnchors[0] = cyclesPerSec(core.anchor)
		core.nanosPerCycles = 1_000_000_000.0 / float64(core.cyclesPerSecond)
	} else if isTscPerCpuStable() {
		core.level = TscPerCpuStable
		cpus, err := availableCpus()
		if len(cpus) == 0 || err != nil {
			core.level = TscUnstable
			return
		}

		max := uint(0)
		for _, cpu := range cpus {
			if cpu > max {
				max = cpu
			}
		}

		res := make([]struct {
			cpu              uint
			cyclesPerSec     uint64
			cyclesFromAnchor uint64
		}, len(cpus))

		for i, cpu := range cpus {
			err := pinCPU(cpu)
			if err != nil {
				panic("")
			}

			res[i].cpu = cpu
			res[i].cyclesPerSec, res[i].cyclesFromAnchor = cyclesPerSec(core.anchor)

			err = unpinCPU()
			if err != nil {
				panic("")
			}
		}

		maxCps := uint64(0)
		minCps := ^uint64(0)
		sumCps := uint64(0)
		core.cyclesFromAnchors = make([]uint64, max+1)
		for _, r := range res {
			cpu, cps, cfa := r.cpu, r.cyclesPerSec, r.cyclesFromAnchor
			if cps > maxCps {
				maxCps = cps
			}
			if cps < minCps {
				minCps = cps
			}
			sumCps += cps
			core.cyclesFromAnchors[cpu] = cfa
		}

		if float64(maxCps-minCps)/float64(minCps) > 0.0005 {
			core.level = TscUnstable
			return
		}

		core.cyclesPerSecond = sumCps / uint64(len(cpus))
		core.nanosPerCycles = 1_000_000_000.0 / float64(core.cyclesPerSecond)

	} else {
		core.level = TscUnstable
	}

	return
}

func (t *TscCore) Now() time.Time {

	switch t.level {
	case TscStable:
		tsc, _ := rdtscp()
		nanos := float64(tsc-t.cyclesFromAnchors[0]) * t.nanosPerCycles
		return t.anchor.Add(time.Nanosecond * time.Duration(nanos))
	case TscPerCpuStable:
		tsc, cpu := rdtscp()
		nanos := float64(tsc-t.cyclesFromAnchors[cpu]) * t.nanosPerCycles
		return t.anchor.Add(time.Nanosecond * time.Duration(nanos))
	default:
		panic("")
	}
	panic("")
}
