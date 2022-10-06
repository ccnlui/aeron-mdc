package cmd

import (
	"aeron-mdc/cmd/client"
	"aeron-mdc/cmd/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println("cancelling...")
		cancel()
	}()

	cmd := newRootCmd()
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("[fatal] error executing command: %v", err)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "aeron-mdc",
		Short: "Aeron multi destination cast example",
	}
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(client.Cmd)
	return rootCmd
}
