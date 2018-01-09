package store

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func initBenchmarkStore(size int) (*store, error) {
	if size < 1 {
		return nil, errors.New("expected size > 0")
	}
	s, err := NewStore(Params{CleaningPeriod: 60 * time.Second}, testClock{}, &testDumper{})
	if err != nil {
		return nil, err
	}
	for i := 0; i < size; i++ {
		switch i % 3 {
		case 0:
			s.Set("item"+strconv.Itoa(i), "value"+strconv.Itoa(i), time.Hour)
		case 1:
			var list []string
			for j := 0; j < 10; j++ {
				list = append(list, fmt.Sprintf("list item %d", j))
			}
			s.ListSet("item"+strconv.Itoa(i), list, time.Hour)
		case 2:
			dict := map[string]string{}
			for j := 0; j < 10; j++ {
				dict["dkey"+strconv.Itoa(j)] = "dict item " + strconv.Itoa(j)
			}
			s.DictSet("item"+strconv.Itoa(i), dict, time.Hour)
		}
	}
	return s.(*store), nil
}

func BenchmarkStore(b *testing.B) {
	for _, size := range []int{
		1,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				initBenchmarkStore(size)
			}
		})
	}
}

func BenchmarkStore_Get(b *testing.B) {
	for _, size := range []int{
		1,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = s.Get("item0")
				if err != nil {
					b.Fatal("unexpected error: " + err.Error())
					b.Fail()
				}
			}
		})
	}
}

func BenchmarkStore_Set(b *testing.B) {
	for _, size := range []int{
		1,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Set("key", "value", time.Hour)
			}
		})
	}
}

func BenchmarkStore_ListGet(b *testing.B) {
	for _, size := range []int{
		3,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := s.ListGet("item1", 0)
				if err != nil {
					b.Fatal("unexpected error: " + err.Error())
					b.Fail()
				}
			}
		})
	}
}

func BenchmarkStore_ListSet(b *testing.B) {
	for _, size := range []int{
		2,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.ListSet("key", []string{"a", "b", "c"}, time.Hour)
			}
		})
	}
}

func BenchmarkStore_DictGet(b *testing.B) {
	for _, size := range []int{
		3,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := s.DictGet("item2", "dkey0")
				if err != nil {
					b.Fatal("unexpected error: " + err.Error())
					b.Fail()
				}
			}
		})
	}
}

func BenchmarkStore_DictSet(b *testing.B) {
	for _, size := range []int{
		3,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.DictSet("key", map[string]string{"a": "a", "b": "b", "c": "c"}, time.Hour)
			}
		})
	}
}

func BenchmarkStore_Remove(b *testing.B) {
	for _, size := range []int{
		1,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Remove("item0")
			}
		})
	}
}

func BenchmarkStore_Keys(b *testing.B) {
	var keys []string
	for _, size := range []int{
		1,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("with store size %d", size), func(b *testing.B) {
			s, err := initBenchmarkStore(size)
			if err != nil {
				b.Fatal("failed to init store: " + err.Error())
				b.FailNow()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				keys = s.Keys()
			}
		})
	}
}
