package main

import (
	"fmt"
	"time"
)

type TestStatus struct {
	acceptCount int
	recvCount int
	parseCount int
	startTime time.Time
	duration int
}

var Testinfo TestStatus

func RecvStatus(statusCh chan string) {
	var statusType string

	Testinfo.duration = 0
	Testinfo.acceptCount = 0
	Testinfo.recvCount = 0
	Testinfo.parseCount = 0
	for {
		statusType = <-statusCh

		if statusType == "A" {
			Testinfo.acceptCount++
		} else if statusType == "R" {
			Testinfo.recvCount++
		} else if statusType == "P" {
			Testinfo.parseCount++
		} else if statusType == "T" {
			Testinfo.startTime = time.Now()
		}
	}
}

func PrintStatus() {
	Testinfo.duration++

	avgA := Testinfo.acceptCount / Testinfo.duration
	avgR := Testinfo.recvCount / Testinfo.duration
	avgP := Testinfo.parseCount / Testinfo.duration

	r2p := Testinfo.recvCount - Testinfo.parseCount

	log := fmt.Sprintf(`A:[%d, %d], R:[%d, %d], P:[%d, %d], R-P:%d`,
		Testinfo.acceptCount, avgA,
		Testinfo.recvCount, avgR,
		Testinfo.parseCount, avgP,
		r2p)
	MainLogger.Info(log)
}
