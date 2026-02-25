package renderer

import (
	"fmt"
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

type TilesTest struct {
	char             uint8
	w, h             int
	fW, fH           int
	offsetW, offsetH int
}

func TestTileGen(t *testing.T) {
	t.Run("Gen", func(t *testing.T) {
		want := []uint8{
			0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2,
			0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2,
			3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5,
			3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5,
			6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8,
			6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8,
		}
		const w, h = 12, 6

		got := [w * h]uint8{}

		til := 9
		tw, th := 4, 2

		tiles := make([]TilesTest, 0)

		y := 0
		twAcc := 0
		for i := range til {

			tiles = append(tiles, TilesTest{
				char: uint8(i),
				w:    tw, h: th,
				fW: w, fH: h,
				offsetW: twAcc,
				offsetH: th * y,
			})

			twAcc += tw

			if twAcc >= w {
				twAcc = 0
				y++
			}
		}

		for _, tt := range tiles {
			for y := range tt.h {
				for x := range tt.w {
					index := (y+tt.offsetH)*tt.fW + (x + tt.offsetW)
					if index >= len(got) {
						continue
					}
					got[index] = tt.char
				}
			}
		}

		for y := range h {
			for x := range w {
				fmt.Printf("%v ", got[y*w+x])
			}
			fmt.Println()
		}

		for y := range h {
			for x := range w {
				if want[y*w+x] != got[y*w+x] {
					t.Errorf("Want x:%d y:%d = %d, got x:%d y:%d = %d", x, y, want[y*w+x], x, y, got[y*w+x])
				}
			}
		}

	})
}
