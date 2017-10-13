package utils

import (
	"errors"
	"reflect"
	"runtime"
	"sort"
	"time"
)

// 本地时间
var loc = time.Local

// Job 定义要执行的任务结构体
type Job struct {
	// 两次任务之间的周期
	interval uint64

	// 要执行的任务函数
	jobFunc string

	// 时间单位
	uint string
	
	// 定义开始运行的时间
	atTime string

	// 最后一次运行的时间
	lastRun time.Time

	// 下次运行的时间
	nextRun time.Time	
}

// NewJob 新建一个任务
func NewJob(intervel uint64) *job {
	return &job{
		intervel,
		"",
		"",
		"",
		""
	}
}

// 是否运行任务
func (j *job) shouldRun() bool {
	return time.Now().After(j.nextRun)
}

// 在给定的时间点上运行任务
func (j *Job) At(t string) *Job {
	hour := int((t[0]-'0')*10 + (t[1] - '0'))
	min := int((t[3]-'0')*10 + (t[4] - '0'))
	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		panic("time format error.")
	}

	mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), int(hour), int(min), 0, 0, loc)

	if j.uint == "days" {
		if time.Now().After(mock) {
			j.lastRun = mock
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, hour, min, 0, 0, loc)
		}
	} else if j.uint == "weeks" {
		if time.Now().After(mock) {
			i := mock.Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), hour, min, 0, 0, loc)
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, hour, min, 0, 0, loc)
		}
	}

	return j
}

// 新建一个任务
func NewScheduler() *Scheduler {
	
}

// 计算下次运行时间
func (j *Job) scheduleNextRun() {
	if j.lastRun == time.Unix(0, 0) {
		if j.unit == "weeks" {
			i := time.Now().Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), 0, 0, 0, 0, loc)

		} else {
			j.lastRun = time.Now()
		}
	}

	if j.period != 0 {
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	} else {
		switch j.unit {
		case "minutes":
			j.period = time.Duration(j.interval * 60)
			break
		case "hours":
			j.period = time.Duration(j.interval * 60 * 60)
			break
		case "days":
			j.period = time.Duration(j.interval * 60 * 60 * 24)
			break
		case "weeks":
			j.period = time.Duration(j.interval * 60 * 60 * 24 * 7)
			break
		case "seconds":
			j.period = time.Duration(j.interval)
		}

		j.nextRun = j.lastRun.Add(j.period * time.Second)
	}
}

// 精确到秒
func (j *Job) Second() *Job {
	if j.interval != 1 {
		panic("")
	}

	job := j.Second()

	return job
}

// 
func (j *Job) Seconds() *Job {
	j.uint = "seconds"
	
	return j
}

// 精确到分钟
func (j *Job) Minute() *Job {
	if j.interval != 1 {
		panic("")
	}

	job = j.Minutes()

	return job
}

//
func (j *Job) Minutes() *Job {
	j.unit = "minutes"

	return j
}

// 精确到小时
func (j *Job) Hour() *Job {
	if j.interval != 1 {
		panic("")
	}

	job = j.Hours()

	return job
}

//
func (j *Job) Hours() *Job {
	j.unit = "hours"

	return j
}

// 
func (j *Job) Day() *Job {
	if j.interval != 1 {
		panic("")
	}

	job = j.Days()

	return job
}

//
func (j *Job) Days() *Job {
	j.unit = "days"

	return j
}

//
func (j *Job) Week() *Job {
	if j.interval != 1 {
		panic("")
	}

	job = j.Weeks()

	return job
}

//
func (j *Job) Weeks() *Job {
	j.unit = "weeks"
	return j
}

// 建任务实例
var defaultScheduler = NewScheduler()
var jobs = defaultScheduler.jobs

func Every(interval uint64) *Job {
	return defaultScheduler.Every(interval)
}

func RunPending() {
	defaultScheduler.RunPending()
}

func RunAll() {
	defaultScheduler.RunAll()
}

func RunAllwithDelay(d int) {
	defaultScheduler.RunAllwithDelay(d)
}

func Start() chan bool {
	return defaultScheduler.Start()
}

func Clear() {
	defaultScheduler.Clear()
}

func Remove(j interface{}) {
	defaultScheduler.Remove(j)
}

func NextRun() (job *Job, time time.Time) {
	return defaultScheduler.NextRun()
}
