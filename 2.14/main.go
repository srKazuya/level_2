package main

import (
	"fmt"
	"time"
)


func or(chans ...<-chan interface{}) <- chan interface{}{
	lenght := len(chans)
	if lenght == 0{
		return nil
	}
	if lenght == 1{
		return chans[0]
	}
	doneCh := make(chan interface{})
	go func() {
		select {
		case <-chans[0]:
		case <-chans[1]:
		case <-or(chans[2:]...):
		}
		defer close(doneCh)
	}()
	return doneCh
}



func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))

}
