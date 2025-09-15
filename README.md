##Setup sequence

1. docker-compose up -d
2. docker exec -it kafka bash (once kafka docker is up)
3. create a kafka topic

```kafka-topics.sh --create --topic <topic-name> --partitions <number> --replication-factor <number> --bootstrap-server localhost:9092```

```--topic  - topic name

    -- how many partitions? - dictates how many consumers it can accomodate

    -- replication-factor - how many machines/instances
```


4. monitor in realtime kafka 

```kafka-console-consumer.sh --topic orders --from-beginning --bootstrap-server localhost:9092```


5. run go . (in main dir of the endpoint)


TODO
consumer endpoints that writes to db?