package renderer

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

const dataSize = 20
const numWorkers = 4

type WGPool struct {
	data []int
	jobs chan int
	wg   sync.WaitGroup
}

func NewWGPool() *WGPool {
	wp := &WGPool{
		data: make([]int, dataSize),
		jobs: make(chan int, dataSize),
	}
	for range numWorkers {
		go wp.worker()
	}
	return wp
}

func (wp *WGPool) worker() {
	for range wp.jobs {
		expensiveWork()
		wp.wg.Done()
	}
}

func (wp *WGPool) Process() {
	wp.wg.Add(len(wp.data))
	for i := range wp.data {
		wp.jobs <- i
	}
	wp.wg.Wait()
}

type ChanPool struct {
	data     []int
	jobs     chan int
	doneChan chan bool
}

func NewChanPool() *ChanPool {
	cp := &ChanPool{
		data:     make([]int, dataSize),
		jobs:     make(chan int, dataSize),
		doneChan: make(chan bool, dataSize),
	}
	for range numWorkers {
		go cp.worker()
	}
	return cp
}

func (cp *ChanPool) worker() {
	for range cp.jobs {
		expensiveWork()
		cp.doneChan <- true
	}
}

func (cp *ChanPool) Process() {
	n := len(cp.data)

	for i := 0; i < n; i += numWorkers {
		batch := numWorkers
		if i+batch > n {
			batch = n - i
		}

		for j := 0; j < batch; j++ {
			cp.jobs <- i + j
		}

		for j := 0; j < batch; j++ {
			<-cp.doneChan
		}
	}
}

//go:noinline
func expensiveWork() {
	for i := range 100000 {
		_ = i * i
	}
}

func BenchmarkPool(b *testing.B) {
	b.Run("WG", func(b *testing.B) {
		wp := NewWGPool()
		for b.Loop() {
			wp.Process()
		}
	})

	b.Run("CH", func(b *testing.B) {
		cp := NewChanPool()
		for b.Loop() {
			cp.Process()
		}
	})
}

func BenchmarkBitShift(b *testing.B) {
	b.Run("shift", func(b *testing.B) {
		for b.Loop() {
			for i := range uint32(20) {
				t := i << 29 >> 29
				t = t + 1
			}
		}
	})
	b.Run("if check", func(b *testing.B) {
		for b.Loop() {
			c := 1
			for i := range uint32(20) {
				if i%8 == 0 {
					c = 0
				}
				c++
			}
		}
	})
}

type TilesTest struct {
	w, h             int
	fW, fH           int
	offsetW, offsetH int
}

func TestTileGen(t *testing.T) {
	t.Run("Gen", func(t *testing.T) {
		want := []uint8{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		}
		const w, h = 14, 12

		got := [w * h]uint8{}

		tileLength := 4

		tiles := make([]TilesTest, 0)

		wOffSet := 0
		hOffSet := 0
		tW, tH := tileLength, tileLength
		for {
			tt := TilesTest{
				w:  tW,
				h:  tH,
				fW: w, fH: h,
				offsetW: wOffSet,
				offsetH: hOffSet,
			}
			fmt.Printf("%+v\n", tt)
			tiles = append(tiles, tt)

			wOffSet += tileLength
			offOffSetW := wOffSet + tileLength

			if offOffSetW > w {
				// fmt.Println("smaller w")
				if offOffSetW-w < tileLength && offOffSetW-w > 0 {
					wOffSet = w - (tileLength - (offOffSetW - w))
					tW = tileLength - (offOffSetW - w)
				} else {
					tW = tileLength
					wOffSet = 0
					hOffSet += tileLength
				}
			}

			offOffSetH := hOffSet + tileLength

			if offOffSetH > h {
				// fmt.Println("smaller h")
				if offOffSetH-h < tileLength && offOffSetH-h > 0 {
					hOffSet = h - (tileLength - (offOffSetH - h))
					tH = tileLength - (offOffSetH - h)
				} else {
					tH = tileLength
				}
			}

			if hOffSet >= h || wOffSet >= w {
				break
			}
		}

		fmt.Println("tiles = ", len(tiles))

		for i, tt := range tiles {
			for y := range tt.h {
				for x := range tt.w {
					index := (y+tt.offsetH)*tt.fW + (x + tt.offsetW)
					if index >= len(got) {
						fmt.Println("jumped", index)
						continue
					}
					got[index] = uint8(i)
				}
			}
		}

		for y := range h {
			for x := range w {
				fmt.Printf("%3d ", got[y*w+x])
			}
			fmt.Println()
		}

		for y := range h {
			for x := range w {
				if want[y*w+x] == got[y*w+x] {
					// t.Errorf("Want x:%d y:%d = %d, got x:%d y:%d = %d", x, y, want[y*w+x], x, y, got[y*w+x])
				}
			}
		}

	})
}

func BenchmarkPow(b *testing.B) {
	b.Run("pow", func(b *testing.B) {
		for b.Loop() {
			p1 := float32(3)
			exp := float32(50)
			p := math.Pow(float64(p1), float64(exp))
			p++
		}
	})

	b.Run("for", func(b *testing.B) {
		for b.Loop() {
			p := 3
			i := p
			for range 50 - 1 {
				p *= i
			}
			p++
		}
	})

	b.Run("unrolled", func(b *testing.B) {
		for b.Loop() {
			p := 3
			x2 := p * p
			x4 := x2 * x2
			x8 := x4 * x4
			x16 := x8 * x8
			x32 := x16 * x16
			x64 := x32 * x32

			x64++
		}
	})
}
