package flashtext

import (
	"context"
	"math"
	"sync/atomic"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025 2025/12/8 下午5:04
* @Package:
 */

const (
	defaultAlpha  = 0.2
	defaultBuffer = 512
)

type densityUpdate struct {
	matches int
	runes   int
}

type stats struct {
	ctx          context.Context
	cancel       context.CancelFunc
	matchDensity uint64
	alpha        float64 // smoothing factor for EWMA (Exponential Weighted Moving Average)
	// EWMA 公式: newDensity = alpha*currentDensity + (1-alpha)*oldDensity
	// alpha 控制本次观测的权重，值越大更新越快，值越小更新越平滑
	//高频变动词库  短文本  0.3~0.5     matchDensity 需要快速跟随变化
	//稳定词库     长文本  0.1~0.2     平滑更新，避免单次异常影响 capEstimate
	//混合场景            0.2         默认通用值，平衡响应速度和稳定性
	updateChan chan densityUpdate
}

func newStats(ctx context.Context, alpha float64, buffer int) *stats {
	ctx, cancel := context.WithCancel(ctx)
	s := &stats{
		ctx:          ctx,
		cancel:       cancel,
		alpha:        alpha,
		matchDensity: float64ToBits(0.01),
		updateChan:   make(chan densityUpdate, buffer),
	}
	go s.dynamicCalculate()
	return s
}

func (s *stats) dynamicCalculate() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case upd := <-s.updateChan:
			if upd.runes == 0 {
				continue
			}
			currentDensity := float64(upd.matches) / float64(upd.runes)
			for {
				oldBits := atomic.LoadUint64(&s.matchDensity)
				oldDensity := float64FromBits(oldBits)
				newDensity := s.alpha*currentDensity + (1.0-s.alpha)*oldDensity
				newBits := float64ToBits(newDensity)
				if atomic.CompareAndSwapUint64(&s.matchDensity, oldBits, newBits) {
					break
				}
			}
		}
	}
}

func (s *stats) add(matches int, runes int) {
	select {
	case s.updateChan <- densityUpdate{matches, runes}:
	default:
	}
}

func (s *stats) getDensity() float64 {
	return float64FromBits(atomic.LoadUint64(&s.matchDensity))
}
func (s *stats) close() {
	s.cancel()
}

func float64ToBits(f float64) uint64   { return math.Float64bits(f) }
func float64FromBits(b uint64) float64 { return math.Float64frombits(b) }
