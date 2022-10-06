package server

import (
	"aeron-mdc/config"
	"aeron-mdc/internal/util"
	"context"
	"log"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start a multi destination cast server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd.Context())
	},
}

func runServer(ctx context.Context) {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	defer a.Close()

	// MDC publishes on server
	pub := <-a.AddPublication(config.ServerChannel, int32(config.TimeStream))
	// for !pub.IsConnected() {
	// 	time.Sleep(time.Millisecond)
	// }
	defer pub.Close()
	log.Println("[info] publication connected to media driver:", pub)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	buf := atomic.MakeBuffer(make([]byte, 64))
	for t := range ticker.C {
		if ctx.Err() != nil {
			return
		}

		bytes := []byte(t.Format(time.RFC3339Nano))
		length := int32(len(bytes))
		buf.PutBytesArray(0, &bytes, 0, length)

		var res int64
		for {
			if res = pub.Offer(buf, 0, length, nil); res > 0 {
				log.Println("[info] sent:", string(bytes))
				break
			}
			if !util.RetryPublicationResult(res) {
				log.Println("[info] dropped:", util.PublicationErrorString(res), string(bytes))
				break
			}
		}
	}
}
