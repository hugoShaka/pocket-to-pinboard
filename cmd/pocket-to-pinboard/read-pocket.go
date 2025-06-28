package main

import (
	"context"
	"fmt"
	pocket_reader "github.com/hugoshaka/pocket-to-pinboard/pkg/pocket-reader"
	"github.com/spf13/cobra"
)

func NewReadPocketCommand(ctx context.Context) (*cobra.Command, error) {
	var pocketConsumerKey string
	readPocketCmd := &cobra.Command{
		Use:   "read-pocket",
		Short: "Lists pocket's content",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("read-pocket called with credentials:", pocketConsumerKey)
			articles, err := pocket_reader.Read(ctx, pocketConsumerKey)
			for _, article := range articles {
				fmt.Printf("%s - %s - %s\n", article.Date.String(), article.Title, article.URL)
				if len(article.Tags) > 0 {
					fmt.Printf("Tags: %#v\n", article.Tags)
				}
			}
			return err
		},
	}
	readPocketCmd.Flags().StringVar(&pocketConsumerKey, "consumer-key", "", "Pocket's consumer key")
	if err := readPocketCmd.MarkFlagRequired("consumer-key"); err != nil {
		return nil, err
	}
	return readPocketCmd, nil
}
