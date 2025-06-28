package pocket_reader

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	pocketapi "github.com/motemen/go-pocket/api"
	pocketauth "github.com/motemen/go-pocket/auth"

	"github.com/hugoshaka/pocket-to-pinboard/pkg/common"
)

func Read(ctx context.Context, consumerKey string) (common.Articles, error) {
	client, err := getClient(ctx, consumerKey)
	if err != nil {
		return nil, fmt.Errorf("could not get client: %w", err)
	}
	articles, err := getAll(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}
	fmt.Printf("Fetched %d articles", len(articles))
	return articles, nil
}

const pageSize = 30

func getAll(ctx context.Context, client *pocketapi.Client) (common.Articles, error) {
	var result common.Articles
	i := 0
	for {
		// The http client used by the pocket lib doesn't make requests with a context :(
		// So we want to catch if the context is cancelled between two API calls.
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		opts := &pocketapi.RetrieveOption{
			Sort:       pocketapi.SortNewest,
			DetailType: pocketapi.DetailTypeSimple,
			Count:      pageSize,
			Offset:     i * pageSize,
		}
		log.Printf("Retrieving articles starting from offset %d", i*pageSize)

		response, err := client.Retrieve(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch results: %w", err)
		}
		for key, item := range response.List {
			// Skip deleted items
			if item.Status == pocketapi.ItemStatusDeleted {
				continue
			}
			// Warn about malformed items
			if item.URL() == "" {
				log.Printf("Item URL is empty: %s - %#v", key, item)
				continue
			}
			title := item.GivenTitle
			if title == "" || title == "Share via" {
				title = item.ResolvedTitle
			}
			tags := make([]string, 0, len(item.Tags))
			for tag, _ := range item.Tags {
				tags = append(tags, tag)
			}
			result = append(result, common.Article{
				Title:       title,
				Description: item.Excerpt,
				URL:         item.URL(),
				Tags:        tags,
				Date:        time.Time(item.TimeAdded),
				Favourite:   item.Favorite == 1,
			})
		}
		i++

		// The lib doesn't expose the "total" field of the response, so I just get until
		// I reach the end.
		if len(response.List) < pageSize {
			break
		}
	}
	return result, nil
}

func getClient(ctx context.Context, consumerKey string) (*pocketapi.Client, error) {
	callBack := make(chan error, 1)
	callBackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			http.Error(w, "Not Found", 404)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err := fmt.Fprintln(w, "You are logged in, you can close this page.")
		callBack <- err
	}))
	defer callBackServer.Close()
	defer close(callBack)

	requestToken, err := pocketauth.ObtainRequestToken(consumerKey, callBackServer.URL)
	if err != nil {
		return nil, err
	}

	url := pocketauth.GenerateAuthorizationURL(requestToken, callBackServer.URL)
	fmt.Printf("Open this URL to log in: %s\n", url)

	select {
	case err := <-callBack:
		if err != nil {
			return nil, fmt.Errorf("failed to fetch log in: %w", err)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	log.Println("Callback received")

	creds, err := pocketauth.ObtainAccessToken(consumerKey, requestToken)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain access token: %w", err)
	}

	log.Println("Access token obtained")

	client := pocketapi.NewClient(consumerKey, creds.AccessToken)
	log.Println("Client built")
	return client, nil
}
