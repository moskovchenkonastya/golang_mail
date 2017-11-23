package main

import (
	"fmt"
	"sync"
	"sort"
	"strings"
)

func ExecutePipeline(arrFunc...job) {
 
	wg := &sync.WaitGroup{}
	wg.Add(len(arrFunc))
	var prev chan interface{}
	// рассмотрим каждую функцию отдельно
	for _, oneFunc := range arrFunc {
		
		next := make(chan interface{})
		go func(workFunc job, in, out chan interface{}) chan interface{} {
			defer wg.Done()
			defer close(out)
			workFunc(in, out)	
			return out
		}(oneFunc, prev, next)
		prev = next
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {

	quotaCh := make(chan struct{}, 1)
	wg1 := &sync.WaitGroup{}

	for data := range in {
		wg1.Add(1)
		go func(data string) {
			var  str1, str3 string
			wg2 := &sync.WaitGroup{}
			wg2.Add(1)
			go func (){
				defer wg2.Done()
				quotaCh <- struct{}{}
				str1 = DataSignerMd5(data)
				<-quotaCh 
				str1 = DataSignerCrc32(str1)
			}()
			str3 = DataSignerCrc32(data) + "~"
			wg2.Wait()
			out <- (str3 + str1) 
			defer wg1.Done()
		}(fmt.Sprintf("%v", data))
	}	
	wg1.Wait()	
}
//
func MultiHash(in, out chan interface{}) {
	
	wg1 := &sync.WaitGroup{}
	for data := range in {
		arr1 := make([]string, 6)
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			wg := &sync.WaitGroup{}
			wg.Add(6)
			for i := 0; i < 6; i++{
				go func(th int, data string) {
					arr1[th] = DataSignerCrc32(fmt.Sprintf("%v",th) + data)
					defer wg.Done()
				}(i, fmt.Sprintf("%v", data))
			}
			wg.Wait()
			out <- strings.Join(arr1, "")
		}()	
	}
	wg1.Wait()
}

func CombineResults(in, out chan interface{}) {
	var arr [] string
	for data := range in {
		arr = append(arr, fmt.Sprintf("%v", data))
	}
	sort.Strings(arr)
	out<- strings.Join(arr, "_") 
}