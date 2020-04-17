package idworker

import (
	"fmt"
	"sync"
	"testing"
)

func TestIdWorker(t *testing.T) {
	var scene sync.Map
	arr := make([]*IdWorker, 0)
	for i := 0; i < 3; i++ {
		worker, _ := NewIdWorker(int64(i))
		arr = append(arr, worker)
	}
	var wg sync.WaitGroup
	count := 3
	for i := 0; i < count; i++ {
		wg.Add(1)
		worker := arr[i]
		go func() {
			defer wg.Add(-1)
			for j := 0; j < 10000; j++ {

				id, err := worker.NextId()
				if err != nil {
					t.Errorf("ID NextId is err! %s \n", err.Error())
					return
				}
				if _, ok := scene.Load(id); ok {
					t.Error("ID is not unique!\n")
					return
				}
				scene.Store(id, 1)
				fmt.Println(id)
			}
		}()
	}
	wg.Wait()
	// 成功生成 idworker ID
	fmt.Println("All", count*10000, "idworker ID Get successed!")
}
