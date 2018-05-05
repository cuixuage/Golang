package main

import "fmt"

//buffer_size is 3
var channelBuffer = make(chan int, 3)
var finalProducer = make(chan bool)
var finalConsumer = make(chan bool)

func produce() {
	for i := 0; i < 10; i++ {
		channelBuffer <- i
	}

	//finish
	finalProducer <- true
}

func consume() {
	for {
		select {
		// when producer finish
		case <- finalProducer:
			finalConsumer <- true
		//continue to read from channelBuffer
		case message := <- channelBuffer :
			fmt.Print(message)
		}
	}
}

//TODO close the channel

func main() {
	go produce()
	go consume()
	consume := <- finalConsumer
	fmt.Print("all over, consumer=",consume)
	close(channelBuffer)
	close(finalConsumer)
	close(finalProducer)
}

