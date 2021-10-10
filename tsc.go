package titime

import (
	"bufio"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

func rdtscp() (uint64, uint64)

func isTscStable() bool {
	bs, err := ioutil.ReadFile("/sys/devices/system/clocksource/clocksource0/available_clocksource")
	if err != nil {
		return false
	}
	return strings.Contains(string(bs), "tsc")
}

func isTscPerCpuStable() bool {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return false
	}
	defer file.Close()

	var line string
	reader := bufio.NewReaderSize(file, 1024)
	for {
		line, err = readLine(reader)
		if err != nil {
			return false
		}
		if strings.HasPrefix(line, "flags") {
			break
		}
	}

	hasConstantTsc := strings.Contains(line, "constant_tsc")
	hasNonStopTsc := strings.Contains(line, "nonstop_tsc")
	hasRdtscp := strings.Contains(line, "rdtscp")

	return hasConstantTsc && hasNonStopTsc && hasRdtscp
}

func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

func cyclesPerSec(anchor time.Time) (uint64, uint64) {
	cps, lastMonotonic, lastTsc := cyclesPerSecInner()
	nanosFromAnchor := lastMonotonic.Sub(anchor).Nanoseconds()
	cyclesFlied := float64(cps) * float64(nanosFromAnchor) / 1_000_000_000.0
	cyclesFromAnchor := lastTsc - uint64(cyclesFlied)

	return cps, cyclesFromAnchor
}

func cyclesPerSecInner() (uint64, time.Time, uint64) {
	var cyclesPerSec float64
	var lastMonotonic time.Time
	var lastTsc uint64
	var oldCycles = 0.0

	for {
		t1, tsc1 := monotonicWithTsc()
		for {
			t2, tsc2 := monotonicWithTsc()
			lastMonotonic = t2
			lastTsc = tsc2
			elapsedNanos := t2.Sub(t1).Nanoseconds()
			if elapsedNanos > 10_000_000 {
				cyclesPerSec = float64(tsc2-tsc1) * 1_000_000_000.0 / float64(elapsedNanos)
				break
			}
		}
		delta := math.Abs(cyclesPerSec - oldCycles)
		if delta/cyclesPerSec < 0.00001 {
			break
		}
		oldCycles = cyclesPerSec
	}

	return uint64(cyclesPerSec), lastMonotonic, lastTsc
}

func monotonicWithTsc() (time.Time, uint64) {
	tsc, _ := rdtscp()
	return time.Now(), tsc
}
