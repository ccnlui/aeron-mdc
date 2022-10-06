package client

import (
	"aeron-mdc/config"
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start a multi destination cast client",
	Run: func(cmd *cobra.Command, args []string) {
		runClient(cmd.Context())
	},
}

func runClient(ctx context.Context) {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	defer a.Close()

	// MDC subscribes to server
	channel := generateNextClientChannel()
	fmt.Println("[info] subscribing to channel:", channel)
	sub := <-a.AddSubscription(channel, int32(config.TimeStream))
	for !sub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	defer sub.Close()
	log.Println("[info] subscription connected to media driver:", sub)

	buf := bytes.NewBuffer(make([]byte, 64))
	onMessage := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		buf.Reset()
		buffer.WriteBytes(buf, offset, length)
		fmt.Println("received:", buf.String())
	}
	assembler := aeron.NewFragmentAssembler(onMessage, 512)
	idleStrategy := idlestrategy.NewDefaultBackoffIdleStrategy()
	for {
		if ctx.Err() != nil {
			return
		}
		workDone := sub.Poll(assembler.OnFragment, 10)
		idleStrategy.Idle(workDone)
	}
}

func generateNextClientChannel() string {
	min := 4096
	max := 65535
	rand.Seed(time.Now().Unix())
	port := rand.Intn(max-min+1) + min
	return fmt.Sprintf(config.ClientChannelFormat, port)
}
