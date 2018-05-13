package crawler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/iahmedov/crawler/model"

	"github.com/iahmedov/crawler/fetcher"
	filtercontent "github.com/iahmedov/crawler/filter/content"
	filterlink "github.com/iahmedov/crawler/filter/link"
	filterstrategy "github.com/iahmedov/crawler/filter/strategy"
	"github.com/iahmedov/crawler/link"
	"github.com/iahmedov/crawler/queue"
	"github.com/iahmedov/crawler/storage"
	"github.com/iahmedov/crawler/task"
)

type Crawler interface {
	Run(ctx context.Context) error
}

type crawler struct {
	workers  int
	executor Executor

	config Config

	queue             queue.Queue
	stateTransitioner task.StateTransitioner
	filters           struct {
		link            []filterlink.Filter
		content         []filtercontent.Filter
		linkStrategy    filterstrategy.Strategy
		contentStrategy filterstrategy.Strategy
	}

	linkExtractors []link.Extractor
	linkStores     []storage.LinkStore
	pageStores     []storage.PageStore
}

func NewCrawler(config Config) Crawler {
	c := &crawler{
		workers: config.Workers,
		config:  config,
	}

	c.stateTransitioner = task.LoadStateTransitioner(c.config.StateTransition)
	c.queue = queue.LoadQueue(c.config.Queue)
	c.queue.SetStateTransitioner(c.stateTransitioner) // TODO: ixtiyor, looks

	c.filters.link = filterlink.LoadFilters(c.config.Filters.Link)
	c.filters.content = filtercontent.LoadFilters(c.config.Filters.Content)
	c.filters.linkStrategy = filterstrategy.LoadStrategy(c.config.Filters.LinkStrategy)
	c.filters.contentStrategy = filterstrategy.LoadStrategy(c.config.Filters.ContentStrategy)

	c.linkExtractors = link.LoadExtractors(c.config.Extractors)

	c.linkStores = storage.LoadLinkStores(c.config.Storage)
	c.pageStores = storage.LoadPageStores(c.config.Storage)
	return c
}

func (c *crawler) Run(ctx context.Context) error {
	c.executor = ExecutorWithContext(ctx)

	c.executor.Go(func(ctx context.Context) error {
		tasks := c.queue.Tasks(ctx)
		client := c.makeClient(c.config.Fetchers)
		for t := range tasks {
			c.stateTransitioner.Transition(t, task.TaskStateProcessing, "")
			entry, err := c.fetchWith(client, t.URL())
			if err != nil {
				c.stateTransitioner.Transition(t, task.TaskStateFailed, err.Error())
				c.storeCrawlEntry(entry)
				continue
			}

			_ = c.processEntry(entry, t)
			c.stateTransitioner.Transition(t, task.TaskStateProcessed, "")
		}
		return nil
	}, c.workers)

	return c.executor.Wait()
}

func (c *crawler) Stop() {
	c.executor.Cancel()
}

func (c *crawler) processEntry(entry *model.CrawlEntry, t task.Task) error {
	// for _, f := range c.filters.link
	return nil
}

func (c *crawler) storeCrawlEntry(entry *model.CrawlEntry) error {
	for _, store := range c.pageStores {
		// TODO: no rollback when multiple stores
		if err := store.AppendCrawlEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (c *crawler) fetchWith(client *http.Client, path url.URL) (*model.CrawlEntry, error) {
	response, err := client.Get(path.String())
	entry := model.NewCrawlEntry(path)
	if err != nil {
		entry.Error = err
		return entry, err
	}

	entry.StatusCode = response.StatusCode
	entry.Header = response.Header
	entry.Error = nil
	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		entry.Error = err
		return entry, err
	}

	entry.Document = model.Document(buff)
	return entry, nil
}

func (c *crawler) makeClient(fetcherConfigs []fetcher.Config) *http.Client {
	middlewares := fetcher.LoadMiddlewares(c.config.Fetchers)
	httpRoundTripper := http.DefaultTransport.RoundTrip
	middleware := fetcher.Chain(fetcher.NoopMiddleware, middlewares...)

	httpClient := *http.DefaultClient
	httpClient.Transport = fetcher.MiddlewareTransport(middleware(httpRoundTripper))
	return &httpClient
}
