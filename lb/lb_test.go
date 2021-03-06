package lb

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"sync"
	"testing"

	"github.com/rs/xid"
)

func TestNewLoadBalancer(t *testing.T) {
	lb := NewLoadBalancer(10, func(id int64) (interface{}, error) {
		return "get " + strconv.Itoa(int(id)), nil
	})
	mu := sync.Mutex{}
	var result = map[interface{}]int{}

	wg := &sync.WaitGroup{}
	num := 200000
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			instance, _ := lb.Get()
			result[instance]++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(result)
	for _, n := range result {
		if n != num/10 {
			t.Errorf("error expect %d got %d", num/10, n)
		}
	}
}

func TestLoadBalancer_Remove(t *testing.T) {
	lb := NewLoadBalancer(10, func(id int64) (interface{}, error) {
		return "get " + strconv.Itoa(int(id)), nil
	})
	for i := 0; i < 10; i++ {
		lb.Get()
		if len(lb.lists) != i+1 {
			t.Errorf("lb.lists error want %d got %d", i+1, len(lb.lists))
		}
	}

	lb.Remove(0)
	lb.Remove(10)
	lb.Remove(1)
	lb.Remove(2)

	for i := 0; i < 1000; i++ {
		lb.Get()
	}

	for _, list := range lb.lists {
		fmt.Println(list.id)
	}
	fmt.Println(lb.current)
}

func TestLoadBalancer_RemoveWithCallbackError(t *testing.T) {
	lb := NewLoadBalancer(10, func(id int64) (interface{}, error) {
		return "get " + strconv.Itoa(int(id)), nil
	})

	for i := 0; i < 10; i++ {
		lb.Get()
	}

	if lb.current != 9 {
		t.Errorf("want 9 got %d", lb.current)
	}
	lb.Remove(9)

	if lb.current != 8 {
		t.Errorf("want 10 got %d", lb.current)
	}
	for i := 0; i < 10; i++ {
		lb.Get()
	}

	if len(lb.lists) != 10 {
		t.Errorf("want 10 got %d", len(lb.lists))
	}
}

func BenchmarkNewLoadBalancer(b *testing.B) {
	lb := NewLoadBalancer(1000, func(id int64) (interface{}, error) {
		return "get " + strconv.Itoa(int(id)), nil
	})
	for i := 0; i < b.N; i++ {
		lb.Get()
	}
}

func TestAAA(t *testing.T) {
	s := make([][]byte, 10)
	r := make([]uint32, 10)
	for i := 0; i < 10; i++ {
		s[i] = xid.New().Bytes()
		r[i] = crc32.ChecksumIEEE(s[i]) % 10
	}
	//i := crc32.ChecksumIEEE(a)
	fmt.Println(r, s)
}
