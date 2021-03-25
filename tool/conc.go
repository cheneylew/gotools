package tool

import (
	"math"
	"fmt"
)

type ChanData struct {
	Success bool
	Message string
	Items []interface{}
}

func Conc(jobs []interface{}, process func(a interface{})) {
	ch := make(chan bool)
	for _, job := range jobs {
		go func(c chan bool, job interface{}) {
			process(job)
			c <- true
		}(ch, job)
	}

	finished := 0
	exit := false
	for {

		select {
		case ok := <-ch:
			if ok {
				finished++
			}
			if finished == len(jobs) {
				exit = true
				break
			}
		}
		if exit {
			break
		}
	}

	close(ch)
}

func ConcQueue(jobs []interface{}, process func(jobsSlice []interface{}) []interface{}, maxConc uint8, sliceSize int) []interface{} {
	perPageNum := sliceSize
	pageTotalCount := int(math.Ceil(float64(len(jobs)) / float64(perPageNum)))
	cd := make(chan ChanData)
	maxConcurrencyNum := maxConc
	var curConNum uint8
	curConNum = 0
	curIndex := 0
	finished := 0
	addedCount := 0
	var totalResults []interface{}
	for {
		if curIndex < pageTotalCount && curConNum < maxConcurrencyNum {
			var endIndex int
			if curIndex*perPageNum+perPageNum >= len(jobs) {
				endIndex = len(jobs)
			} else {
				endIndex = curIndex*perPageNum + perPageNum
			}
			targetJobItems := jobs[curIndex*perPageNum : endIndex]

			go func(cd chan ChanData, pageId int, targetJobItems []interface{}) {
				results := process(targetJobItems)
				cd <- ChanData{
					Success: true,
					Message: fmt.Sprintf("%d已完成", pageId),
					Items:   results,
				}
			}(cd, curIndex, targetJobItems)
			curConNum++
			curIndex++
			addedCount++
		}
		//并发满了，等待完成后进入下一个;或者并发未满，但并发池子里还有未完成的任务，等待所有完成才退出。
		if curConNum == maxConcurrencyNum || (addedCount == pageTotalCount && finished != pageTotalCount) {
			next := false
			exit := false
			select {
			case data := <-cd:
				if data.Success && data.Items != nil {
					totalResults = append(totalResults, data.Items...)
				}
				finished++
				curConNum--
				next = true
				if finished == pageTotalCount {
					exit = true
					break
				}
			}
			if exit {
				break
			}
			if next {
				continue
			}
		}
	}

	close(cd)

	return totalResults
}
