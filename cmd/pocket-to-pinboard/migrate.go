package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	pinboard_writer "github.com/hugoshaka/pocket-to-pinboard/pkg/pinboard-writer"
	pocket_reader "github.com/hugoshaka/pocket-to-pinboard/pkg/pocket-reader"
)

func NewMigrateCommand(ctx context.Context) (*cobra.Command, error) {
	var pocketConsumerKey string
	var pinboardAPIKey string
	readPocketCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Lists all pocket content, them migrate it to pinboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔎 Reading from Pocket")
			articles, err := pocket_reader.Read(ctx, pocketConsumerKey)
			if err != nil {
				return err
			}

			fmt.Println("✍️ Importing into Pinboard")
			err = pinboard_writer.Write(ctx, articles, pinboardAPIKey)
			if err != nil {
				return err
			}
			fmt.Println("🎉 Done")
			return nil
		},
	}
	const pocketConsumerKeyFlag = "pocket-consumer-key"
	readPocketCmd.Flags().StringVar(&pocketConsumerKey, pocketConsumerKeyFlag, "", "Pocket's consumer key")
	if err := readPocketCmd.MarkFlagRequired(pocketConsumerKeyFlag); err != nil {
		return nil, err
	}
	const pinboardAPIKeyFlag = "pinboard-api-key"
	readPocketCmd.Flags().StringVar(&pinboardAPIKey, pinboardAPIKeyFlag, "", "Pinboard's API key")
	if err := readPocketCmd.MarkFlagRequired(pinboardAPIKeyFlag); err != nil {
		return nil, err
	}
	return readPocketCmd, nil
}
