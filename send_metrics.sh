#!/bin/bash

randrange() {
    min=$1
    max=$2

    ((range = max - min))

    ((maxrand = 2**30))
    ((limit = maxrand - maxrand % range))

    while true; do
        ((r = RANDOM * 2**15 + RANDOM))
        ((r < limit)) && break
    done

    ((num = min + r % range))
    echo $num
}

for i in {1..100}
do
echo "Running $i times"
VAR=$(randrange 2 50)
echo "Set VAR=$VAR"
sleep 1
curl --request POST \
  --url http://localhost:8081/agentmetrics \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/2023.5.8' \
  --data '{
	"resp_time":'$VAR',
	"req_cnt":1,
	"req_size":20
}'
done
