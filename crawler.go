package crawler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/iahmedov/crawler/model"
	"github.com/pkg/errors"

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
		linkStrategy    filterstrategy.StrategyBuilder
		contentStrategy filterstrategy.StrategyBuilder
	}

	linkExtractors []link.Extractor
	linkStores     []storage.LinkStore
	pageStores     []storage.PageStore
}

func NewCrawler(config Config) (Crawler, error) {
	c := &crawler{
		workers: config.Workers,
		config:  config,
	}

	var err error
	c.linkStores, err = storage.LoadLinkStores(c.config.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "could not load link storage")
	}

	c.pageStores, err = storage.LoadPageStores(c.config.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "could not load page storage")
	}

	c.stateTransitioner, err = task.LoadStateTransitioner(c.config.StateTransition)
	if err != nil {
		return nil, errors.Wrap(err, "could not load state transitioner")
	}
	c.queue, err = queue.LoadQueue(c.config.Queue)
	if err != nil {
		return nil, errors.Wrap(err, "could not load queue")
	}
	c.queue.SetStateTransitioner(c.stateTransitioner) // TODO: looks ugly

	c.filters.link, err = filterlink.LoadFilters(c.config.Filters.Link)
	if err != nil {
		return nil, errors.Wrap(err, "could not load link filters")
	}

	c.filters.content, err = filtercontent.LoadFilters(c.config.Filters.Content)
	if err != nil {
		return nil, errors.Wrap(err, "could not load content filters")
	}

	c.filters.linkStrategy, err = filterstrategy.LoadStrategyBuilder(c.config.Filters.LinkStrategy)
	if err != nil {
		return nil, errors.Wrap(err, "could not load link filter strategy")
	}

	c.filters.contentStrategy, err = filterstrategy.LoadStrategyBuilder(c.config.Filters.ContentStrategy)
	if err != nil {
		return nil, errors.Wrap(err, "could not load content filter strategy")
	}

	c.linkExtractors, err = link.LoadExtractors(c.config.Extractors)
	if err != nil {
		return nil, errors.Wrap(err, "could not load link extractors")
	}

	return c, nil
}

func (c *crawler) Run(ctx context.Context) error {
	c.executor = ExecutorWithContext(ctx)

	c.executor.Go(func(ctx context.Context) error {
		tasks := c.queue.Tasks(ctx)
		client, err := c.makeClient(c.config.Fetchers)
		if err != nil {
			return err
		}

		for t := range tasks {
			c.stateTransitioner.Transition(t, task.TaskStateProcessing, "")
			entry, err := c.fetchWith(client, t)
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
	contentStrategy := c.filters.contentStrategy.Build()
	for _, f := range c.filters.content {
		contentStrategy.AddFilterResult(f(*entry))
	}
	if contentStrategy.Decision().IsNegative() {
		return c.storeCrawlEntry(entry)
	}

	uniqueLinks := map[string]url.URL{}
	for _, linkExtractor := range c.linkExtractors {
		urls := linkExtractor.Extract(*entry)
		for _, u := range urls {
			uniqueLinks[u.String()] = u
		}
	}

	childDepth := t.Depth() + 1
	for _, v := range uniqueLinks {
		linkStrategy := c.filters.linkStrategy.Build()
		for _, f := range c.filters.link {
			linkStrategy.AddFilterResult(f(v, childDepth))
		}
		decision := linkStrategy.Decision()
		if decision.IsNegative() {
			continue
		}

		entry.Links = append(entry.Links, v)
	}

	return c.storeCrawlEntry(entry)
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

func (c *crawler) fetchWith(client *http.Client, t task.Task) (*model.CrawlEntry, error) {
	path := t.URL()
	response, err := client.Get(path.String())
	entry := model.NewCrawlEntry(path, t.Depth())
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

func (c *crawler) makeClient(fetcherConfigs []fetcher.Config) (*http.Client, error) {
	middlewares, err := fetcher.LoadMiddlewares(c.config.Fetchers)
	if err != nil {
		return nil, err
	}
	httpRoundTripper := http.DefaultTransport.RoundTrip
	middleware := fetcher.Chain(fetcher.NoopMiddleware, middlewares...)

	httpClient := *http.DefaultClient
	httpClient.Transport = fetcher.MiddlewareTransport(middleware(httpRoundTripper))
	return &httpClient, nil
}
