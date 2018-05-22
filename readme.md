Plugin based crawler

## Main componenets: 
* seed - not implemented
* queue - supply with what to fetch next
* fetcher - fetching data from remote endpoint using middlewares. could be basic auth, jwt based auth, body limiter, proxy or circuit breaker. Based on http.RoundTripper
* processor - do some stuff on data 
    * contentFilter - filter content based on MetaData and crawlEntry
    * extract - given any type of document try to parse and extract links
    * linkFilter
    * output
        * store fetched content and links in DB
        * put links to queue
* config - everything is controlled by config.yaml, implement any part of crawler and put into config.yaml as an entry crawler will pickup that plugin (similar logic to "database/sql")
* (crawler itself)

### Queue
* implement `crawler/queue.Queue` interface
* create factory for your queue
* register using `crawler.queue.RegisterQueueFactory` 
* import in root folder of executable as you do with database drivers 
```
import _ "github.com/mynick/myqueue"
```
* put your settings in config.yaml

#### Currently available queues: local

### Fetcher
* Steps are same as <a href="#queue">Queue</a>
* implement `crawler/fetcher.Middleware` and register using `crawler/fetcherRegisterMiddlewareFactory`

#### Currently available fetcher: circuitbreaker, limit (body limiter)

### Filter
* Again steps are similar to <a href="#queue">Queue</a>, but filter has more types:

  * Link filter - filtered links will not be added to queue and saved in storage
  * Content filter - filter based on content
  * Strategy - since we have multiple filters, strategy decides based on multiple filter results, what to do next.

#### Currently available link filters: depth, domain, social/facebook, storage unique element

#### Currently available content filters: contenttype, porn

#### Currently available strategies: cumulative, negative_high

### Link extractors
* Again steps are similar to <a href="#queue">Queue</a>
* implement `crawler/link.Extractor` interface and register using `crawler/link.RegisterExtractorFactory`

#### Currently available link extractors: html - extract links from raw html

### Storage
* Again steps are similar to <a href="#queue">Queue</a>
* implement `crawler/storage.LinkStore` and `crawler/storage.PageStore` interface depending on task
* Note: sitemap could be here, since sitemap has been used to output data

#### Currently available storages: memory

## Reliability
* Smallest component in crawling process is `task.Task`, which has States
* When queue receives new URL to process it creates a task and marks state of this task as `TaskStateInitial`
* when crawler receives task and starts processing changes state to `TaskStateProcessing`, after that to `TaskStateProcessed`

* Each transition of state should be handled in order to get reliable crawling process. For this purpose implement `task.StateTransitioner` interface, in order to define your own reliability policy

#### Currently available StateTransitioner: unreliable - does nothing on state transitions. possible reliable would save each task and its state

### overall crawling process
```
Run() {
    for i = 0; i < 10; i++ {
        go func() {
            taskProducer := queue.Tasks()
            for task := range <-taskProducer {
                TaskStateTransitioner.Transition(task, Processing)
                endpointResponse := fetcher(task.url)
                crawlEntry = endpointResponse
                decision := contentFilters(crawlEntry)
                links := linkExtractor(crawlEntry)
                filteredLinks := linkFilters(links)
                queue.Put(filteredLinks)
                storage.Store(crawlEntry, filteredLinks)
                TaskStateTransitioner.Transition(task, Processed)
            }
        }
    }
}
```
