package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	GracefulShutdown()
	fmt.Println("Done (GracefulShutdown func)")
	wMutex()
	fmt.Println("Done (wMutex func)")
	generateRandomNums()
	fmt.Println("Done (generateRandomNums func)")
}

func GracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	timer := time.After(10 * time.Second)

	for {
		select {
		case <-timer:
			fmt.Println("Timeout")
			return
		case sig := <-sigChan:
			fmt.Println("Stopped by signal", sig)
			return
		}
	}
}
func generateRandomNums() {
	var num int
	var wg sync.WaitGroup
	buffernumChan := make(chan int, 5)
	unbuffernumChan := make(chan int)
	wg.Add(4)

	go func() {
		defer wg.Done()
		startBuff := time.Now()
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			num = rand.Intn(100)
			buffernumChan <- num
		}
		close(buffernumChan)
		buffTime := time.Since(startBuff)
		fmt.Println("Time taken (generateRandomNums buff func):", buffTime)
	}()
	go func() {
		defer wg.Done()
		for num := range buffernumChan {
			fmt.Println(num)
		}
	}()

	go func() {
		defer wg.Done()
		startUnbuff := time.Now()
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			num = rand.Intn(100)
			unbuffernumChan <- num
		}
		close(unbuffernumChan)
		unbuffTime := time.Since(startUnbuff)
		fmt.Println("Time taken (generateRandomNums unbuff func):", unbuffTime)
	}()
	go func() {
		defer wg.Done()
		for num := range unbuffernumChan {
			fmt.Println(num)
		}
	}()

	wg.Wait()
}

func wMutex() {
	start := time.Now()
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1000)
	mu.Lock()
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Second)
			counter++
		}()
	}
	mu.Unlock()

	wg.Wait()
	fmt.Println("Counter:", counter)
	fmt.Println("Time taken (wMutex func):", time.Since(start))
}
