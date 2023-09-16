### prometheus metrics
### Setup Test 

```

1. Run prometheus in docker.
docker-compose up

2. Start this server - 
simulates agent controller, that will receive metrics on /agentmetrics endpoint
that will be scraped by prometheus.
go run main.go

Note - This endpoint would be called by agent to send metrics every x sec to agent-controller. 

3. Send metrics via API to above service endpoint started in #2.

curl --request POST \
  --url http://localhost:8081/agentmetrics \
  --header 'Content-Type: application/json' \
  --data '{
	"resp_time":50,
	"req_cnt":1,
	"req_size":20
}'

4. Check metrics ok in prometheus.

```

# How to view metrics on Prometheus
```
1. check target is ok.
http://localhost:9090/targets

2. check metrics populated by correct label for metric=resp_time_seconds
http://localhost:9090/graph
Enter below query-
resp_time_seconds
resp_time_seconds{tenantid='tenant1'}
resp_time_seconds{agentid='agent1'}

Note - depending on the scrapr interval set in prometheus.yml, it may take some time for the metrics to show on the prometheus UI
```
