// 一个端口扫描的小东西, 控制在 2000 并发以内

package app

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	// 待处理作业通道
	jobs = make(chan int, 100)
	// 用于检查是否完成
	wg sync.WaitGroup
)

// 不断从 jobs 取出数据, 能取出就处理, 然后放入 results
// jobs 执行完成后, jobs 通道关闭
// 此时 goroutine 取出失败也会退出, wg -1
func worker(ip string, results chan<- int) {
	defer wg.Done()

	for v := range jobs {
		addr := fmt.Sprintf("%v:%v", ip, v)
		_, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
		if err == nil {
			results <- v
		}
	}
}

func showResult(results <-chan int) {
	r := make([]int, 0, 200)

	for v := range results {
		r = append(r, v)
	}

	fmt.Println(r)
}

func checkFlag(counts, poolCap int) {
	// 检查端口范围
	if counts < 0 {
		log.Fatalf("range error, End > Start")
	}

	// 检查线程池大小
	if poolCap > 2000 {
		log.Fatalf("poolCap is too large")
	}
}

func main() {

	// 获取参数
	var (
		host    string
		start   int
		end     int
		poolCap int
	)

	flag.StringVar(&host, "host", "localhost", "要扫描的主机IP")
	flag.IntVar(&start, "start", 1, "端口起始位置")
	flag.IntVar(&end, "end", 65535, "端口终止位置")
	flag.IntVar(&poolCap, "poolCap", 200, "线程池容量")

	flag.Parse()

	// 计算个数
	counts := end - start + 1

	// 检查 flag 参数
	checkFlag(counts, poolCap)

	// 用于存放结果的 chan
	results := make(chan int, counts)

	log.Println("--- Start ---")

	// poolCap 个 goroutine 做事 --> workpool
	// 因为 wg 的限制, 上限不能超过 13000+
	wg.Add(poolCap)
	for id := 0; id < poolCap; id++ {
		go worker(host, results)
	}

	// 将需要扫描的端口范围放入 jobs chan
	for index := start; index < end; index++ {
		jobs <- index
	}

	// 都放入 jobs 后, 关闭 jobs
	close(jobs)

	// 等待 goroutine 都取出失败后, 表示 results 完整, 关闭 results
	wg.Wait()

	// 关闭通道, 打印结果
	close(results)
	showResult(results)

	log.Println("--- End ---")
}
