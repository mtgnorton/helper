package helper

import (
	"math"
	"sync"
)

type RollingAverage struct {
	mu      sync.Mutex
	count   int
	average float64 // 平均值
	ma      float64 // 滑动平均值
	cma     float64 // 修正后的滑动平均值
	Beta    float64
}

func NewRollingAverage(initialValue, Beta float64) *RollingAverage {
	return &RollingAverage{
		count:   0,
		average: initialValue,
		ma:      initialValue,
		cma:     initialValue,
		Beta:    Beta,
	}
}
func (ra *RollingAverage) Add(v float64) {
	ra.mu.Lock()
	ra.ma = ra.ma*ra.Beta + v*(1-ra.Beta)
	ra.average = (ra.average*float64(ra.count) + v) / float64(ra.count+1)
	ra.count++
	ra.cma = ra.ma / (1 - math.Pow(ra.Beta, float64(ra.count)))
	ra.mu.Unlock()
}

func (ra *RollingAverage) Get() (int, float64, float64, float64) {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	return ra.count, ra.average, ra.ma, ra.cma
}
