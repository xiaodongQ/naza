package main

import (
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/taskpool"
	"sync"
	"time"
)

var (
	taskNum       = 1000 * 1000
	initWorkerNum = 1 //1000 * 20 //1000 * 10
)

func originGo() {
	nazalog.Debug("> BenchmarkOriginGo")
	var wg sync.WaitGroup
	for j := 0; j < 1; j++ {
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			go func() {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			}()
		}
		wg.Wait()
	}
	nazalog.Debug("< BenchmarkOriginGo")
}

func taskPool() {
	var poolNum = 1

	nazalog.Debug("> BenchmarkTaskPool")
	var wg sync.WaitGroup
	var ps []taskpool.Pool
	for i := 0; i < poolNum; i++ {
		p, _ := taskpool.NewPool(func(option *taskpool.Option) {
			option.InitWorkerNum = initWorkerNum
		})
		ps = append(ps, p)
	}

	for j := 0; j < 1; j++ {
		//b.StartTimer()
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			ps[i % poolNum].Go(func() {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			})
		}
		wg.Wait()
		//b.StopTimer()
		//idle, busy := p.Status()
		//nazalog.Debugf("done, worker num. idle=%d, busy=%d", idle, busy) // 此时还有个别busy也是正常的，因为只是业务方的任务代码执行完了，可能还没回收到idle队列中
		//p.KillIdleWorkers()
		//idle, busy = p.Status()
		//nazalog.Debugf("killed, worker num. idle=%d, busy=%d", idle, busy)
	}
	nazalog.Debug("< BenchmarkTaskPool")
	idle, busy := ps[0].Status()
	nazalog.Debugf("killed, worker num. idle=%d, busy=%d", idle, busy)
}

func main() {
	taskPool()
	//originGo()
	nazalog.Debug("waiting exit.")
	time.Sleep(1000 * time.Second)
	//nazalog.Debug("bye.")
}
