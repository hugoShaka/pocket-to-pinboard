package pinboard_writer

import (
	"context"
	"errors"
	"fmt"

	"github.com/imwally/pinboard"

	"github.com/hugoshaka/pocket-to-pinboard/pkg/common"
)

func Write(ctx context.Context, articles common.Articles, apiKey string) error {
	if apiKey == "" {
		return errors.New("api key is required")
	}
	if len(articles) == 0 {
		return errors.New("articles is required")
	}
	pinboard.SetToken(apiKey)
	for i, article := range articles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		opts := &pinboard.PostsAddOptions{
			URL:         article.URL,
			Description: article.Title,
			Extended:    []byte(article.Description),
			Tags:        article.Tags,
			Dt:          article.Date,
			Replace:     true,
			Shared:      true,
			Toread:      false,
		}
		fmt.Printf("Adding article %d/%d: %q\n", i, len(articles), article.Title)
		err := pinboard.PostsAdd(opts)
		if err != nil {
			return fmt.Errorf("adding post %q: %w", article.Title, err)
		}
	}
	return nil

}
