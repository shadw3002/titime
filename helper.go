package titime

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func availableCpus() ([]uint, error) {
	res := []uint{}
	bs, err := ioutil.ReadFile("/sys/devices/system/clocksource/clocksource0/available_clocksource")
	if err != nil {
		return []uint{}, err
	}
	l := strings.Trim(string(bs), " ")
	list := strings.Split(l, ",")
	for _, set := range list {
		if strings.Contains(set, "-") {
			ss := strings.SplitN(set, "-", 2)
			if len(ss) != 2 {
				return []uint{}, fmt.Errorf("TODO")
			}
			from, err := strconv.Atoi(ss[0])
			if err != nil {
				return []uint{}, err
			}
			to, err := strconv.Atoi(ss[1])
			if err != nil {
				return []uint{}, err
			}
			for i := from; i <= to; i++ {
				res = append(res, uint(i))
			}
		} else {
			i, err := strconv.Atoi(set)
			if err != nil {
				return []uint{}, err
			}
			res = append(res, uint(i))
		}
	}

	return res, nil
}
