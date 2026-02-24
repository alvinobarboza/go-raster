package scene

import (
	"sync"
	"testing"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func valueTester(v1, v2, v3 transforms.Vec3) transforms.Vec3 {
	return v1.Divide(2).Add(v2.Divide(2)).Add(v3.Divide(2))
}

func referenceTester(vs []transforms.Vec3, i1, i2, i3 int) transforms.Vec3 {
	return vs[i1].Divide(2).Add(vs[i2].Divide(2)).Add(vs[i3].Divide(2))
}

func BenchmarkPassByValue(b *testing.B) {
	b.Run("value", func(b *testing.B) {
		v := transforms.NewVec3(2, 3, 4)
		for b.Loop() {
			v = valueTester(v, v, v)
		}
	})

	b.Run("reference", func(b *testing.B) {
		vs := make([]transforms.Vec3, 0)

		for range 200 {
			vs = append(vs, transforms.NewVec3(2, 3, 4))
		}

		for b.Loop() {
			referenceTester(vs, 0, 3, 4)
		}
	})
}

func doTask(threads int) func(task func(), c bool) {
	tasks := make(chan func())

	for range threads {
		go func() {
			for t := range tasks {
				t()
			}
		}()
	}

	return func(task func(), c bool) {
		if c {
			close(tasks)
			return
		}
		tasks <- task
	}
}

func BenchmarkConcurrent(b *testing.B) {
	const size = 1000000
	const numGoroutines = 8
	const chunkSize = size / numGoroutines
	b.Run("NoMutex", func(b *testing.B) {
		var arr [size]int

		for b.Loop() {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for g := range numGoroutines {
				go func(id int) {
					defer wg.Done()
					start := id * chunkSize
					end := start + chunkSize

					for i := start; i < end; i++ {
						arr[i] = i
					}
				}(g)
			}

			wg.Wait()
		}
	})

	b.Run("WithMutex", func(b *testing.B) {
		var arr [size]int
		var mu sync.Mutex

		for b.Loop() {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for g := range numGoroutines {
				go func(id int) {
					defer wg.Done()
					start := id * chunkSize
					end := start + chunkSize

					for i := start; i < end; i++ {
						mu.Lock()
						arr[i] = i
						mu.Unlock()
					}
				}(g)
			}

			wg.Wait()
		}
	})

	b.Run("WithChanneledMutex", func(b *testing.B) {
		var arr [size]int
		var mu sync.Mutex

		poolTask := doTask(numGoroutines)

		for b.Loop() {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for g := range numGoroutines {
				poolTask(func() {
					defer wg.Done()
					start := g * chunkSize
					end := start + chunkSize

					for i := start; i < end; i++ {
						mu.Lock()
						arr[i] = i
						mu.Unlock()
					}
				}, false)
			}

			wg.Wait()
		}

		poolTask(func() {}, true)
	})

	b.Run("NoMutexChanneled", func(b *testing.B) {
		var arr [size]int

		poolTask := doTask(numGoroutines)

		for b.Loop() {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for g := range numGoroutines {
				poolTask(func() {
					defer wg.Done()
					start := g * chunkSize
					end := start + chunkSize

					for i := start; i < end; i++ {
						arr[i] = i
					}
				}, false)
			}

			wg.Wait()
		}

		poolTask(func() {}, true)
	})
}
