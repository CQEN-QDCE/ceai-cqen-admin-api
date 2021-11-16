package services

import (
	"sync"
)

// Parallelize parallelizes function calls
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func RemoveStringElementFromArray(array []string, element string) []string {
	index := -1

	for i, e := range array {
		if e == element {
			index = i
			break
		}
	}

	if index != -1 {
		ret := make([]string, 0)
		ret = append(ret, array[:index]...)
		return append(ret, array[index+1:]...)
	} else {
		return array
	}
}
