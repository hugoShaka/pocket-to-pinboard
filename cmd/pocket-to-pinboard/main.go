package main

import (
	"context"
	"log"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "pocket-to-pinboard",
		Short: "A program to fetch content from pocket and upload to pinboard.",
	}
	ctx := context.Background()
	pocketReadCmd, err := NewReadPocketCommand(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cmd.AddCommand(pocketReadCmd)
	migrateCmd, err := NewMigrateCommand(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cmd.AddCommand(migrateCmd)
	if err := fang.Execute(ctx, cmd); err != nil {
		os.Exit(1)
	}
}
