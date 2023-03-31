### redis 客户端

redis经常用来处理复杂业务而不是单纯的字符串存储，比如使用Redis集合、分布式锁、队列等，因此集成lettuce +
Redisson，按照不同业务而使用不同组件是比较优秀的方案。

Jedis is a straight-forward Redis client that is not thread-safe when applications want to share a single Jedis instance
across multiple threads. The approach to use Jedis in a multi-threaded environment is to use connection pooling. Each
concurrent thread using Jedis gets its own Jedis instance for the duration of Jedis interaction. Connection pooling
comes at the cost of a physical connection per Jedis instance which increases the number of Redis connections.

Lettuce is built on netty and connection instances (StatefulRedisConnection) can be shared across multiple threads. So a
multi-threaded application can use a single connection
regardless the number of concurrent threads that interact with Lettuce.

### 普通的单机请求

```shell

ujucom@j:~$ wrk -t12 -c400 -d30s http://localhost:8080/api/test
Running 30s test @ http://localhost:8080/api/test
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   182.19ms  311.69ms   2.00s    86.90%
    Req/Sec   481.61    265.61     2.94k    76.72%
  168014 requests in 30.08s, 18.77MB read
  Socket errors: connect 0, read 0, write 0, timeout 451
Requests/sec:   5585.06
Transfer/sec:    639.06KB

```

### elastic  
#### 查询方式  
* Query creation  
  Generally the query creation mechanism for Elasticsearch works as described in Query Methods
```java
Criteria miller = new Criteria("lastName").is("Miller")  
  .subCriteria(                                          
    new Criteria().or("firstName").is("John")            
      .or("firstName").is("Jack")                        
  );
Query query = new CriteriaQuery(criteria);
```
  
* StringQuery  
  This class takes an Elasticsearch query as JSON String.
```java
Query query = new StringQuery("{ \"match\": { \"firstname\": { \"query\": \"Jack\" } } } ");
SearchHits<Person> searchHits = operations.search(query, Person.class);
```

* NativeQuery  
  NativeQuery is the class to use when you have a complex query, or a query that cannot be expressed by using the Criteria API,
```java
Query query = NativeQuery.builder()
	.withAggregation("lastNames", Aggregation.of(a -> a
		.terms(ta -> ta.field("last-name").size(10))))
	.withQuery(q -> q
		.match(m -> m
			.field("firstName")
			.query(firstName)
		)
	)
	.withPageable(pageable)
	.build();

SearchHits<Person> searchHits = operations.search(query, Person.class);
```


