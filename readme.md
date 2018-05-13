main componenets: 
* queue - supply with what to fetch next
* fetcher - fetch data from remote endpoint
* processor - do some stuff on data 
    * parse - given any type of document try to parse and set result to MetaData
    * contentFilter - filter content based on MetaData and crawlEntry
    * generate links
    * linkFilter
    * output
        * store fetched content and links in DB
        * put links to queue
* (crawler itself)

// queue
queueStrategy := FIFO{}
seed: Seed(Link("https://abc.com")), // Seed(Database("connstring")), Seed(Sitemap("https://abc.com/sitemap.xml"))
queue := Queue(queueStrategy, seed) // store links, states, 

// Fetcher
httpClient := makeHttpClientWith(Auth, CircuitBreaker, WithBodyLimit(100K, put wrapper around Body.ReadCloser)) use http.Transport
fetcher := httpClient.GET("..........")

// processor
//  parsers (JSON, PDF, HTML) (ParseRequest) -> ParseResponse
//      crawlEntry.MetaData = crawlEntry.Put("parserName", ParseResponse)
//  link cleaners (not necessary at this moment)
//  generateLinks(crawlEntry.MetaData) -> []Links
//  filteredLinks = linkFilters(links)
//  uniqueLinks = unique(filteredLinks)
//  outputs (Storage, Queue) (CrawlEntry)

// linkFilters
linkFilters: []Filter{
    RobotsFilter,
    FacebookFilter,
    DepthFilter,
    SitemapFilter("https://abc.com/sitemap.xml"),
    HasCrawledDBChecker(linkStorage),
}
strategy: FirstPositive, FirstNegative, Cumulative,...
filter := chainFilter(strategy, filters)

// storage
linkStorage := cached{DBLinkStorage{}} // place where we store document links
documentStorage := cached{DBDocumentStorage{}} // place where we store documents

// reliability
using tasks with states
StateTransitioner: taskStorage := DBTaskStorage{}
StateTransitioner: taskStorage := None{} // non reliable

config := crawler.Config{
    queue: queue.Cached{queue.DatabaseQueue()},
    dataFilter: []Filter{
        ContainsFilter("abc")
    },
    action: chainAction(basicAuth, endpointRequester)
    linkStorage: linkStorage,
    documentStorage: documentStorage,
}
crawler := NewCrawler(config)
crawler.Run()

// initialization
NewCrawler(config) {
    crawler := &Crawler{}
    crawler.queue := config.queue
    return crawler
}

// crawling process
Run() {
    for i = 0; i < 10; i++ {
        go func() {
            taskProducer := queue.Tasks()
            for task := range taskProducer {
                endpointResponse := endpoint(url)
                crawlEntry = endpointResponse
                TaskStateTransitioner.Transition(task, Processing)
                Process(task, crawlEntry)
                TaskStateTransitioner.Transition(task, Processed)
            }
        }
    }
}
