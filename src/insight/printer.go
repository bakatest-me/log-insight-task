package insight

import (
	"fmt"
	"log-insign-task/domain"
	"sort"
)

// utility

type Printer interface {
	map[string]int | map[int]int | []domain.KV[int, int] | []domain.KV[string, int]
}

func Print[T Printer](list T) {
	switch l := any(list).(type) {
	case map[string]int:
		for k, v := range l {
			fmt.Printf("%s: %d\n", k, v)
		}
	case map[int]int:
		for k, v := range l {
			fmt.Printf("%d: %d\n", k, v)
		}
	case []domain.KV[string, int]:
		for _, v := range l {
			fmt.Printf("%v: %v\n", v.Key, v.Value)
		}
	case []domain.KV[int, int]:
		for _, v := range l {
			fmt.Printf("%v: %v\n", v.Key, v.Value)
		}
	}
}

func sortByKeyInt(m map[int]int) (ss []domain.KV[int, int]) {
	for k, v := range m {
		ss = append(ss, domain.KV[int, int]{Key: k, Value: v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key < ss[j].Key
	})

	return ss
}

func sortTopRank(m map[string]int) (ss []domain.KV[string, int]) {
	for k, v := range m {
		ss = append(ss, domain.KV[string, int]{Key: k, Value: v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

func mergeMapInt(m1, m2 map[int]int) map[int]int {
	for k, v := range m2 {
		m1[k] += v
	}
	return m1
}

func mergeMap(m1, m2 map[string]int) map[string]int {
	for k, v := range m2 {
		m1[k] += v
	}
	return m1
}
