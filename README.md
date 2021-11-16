# url-shortener
> Written in Hexagonal Architecture

This is a small API that provides basic REST endpoints to shorten a URL, get information about the URL, update the URL, and get statistics on most accessed URLs.

The technology behind it: 
* [Golang 1.17](https://golang.org/doc/go1.17)
* [Gin Web Framework](https://gin-gonic.com/docs/)
* [Google Firestore](https://cloud.google.com/firestore)
* [Google PubSub](https://cloud.google.com/pubsub)
* [Redis](https://redis.io/)
* [NanoID](https://github.com/ai/nanoid)
* [Swagger](https://swagger.io/)

## Architecture

![Architecture](doc/arquitecture.png?raw=true "Architecture")

1. Requests can come from many types of devices.
2. Cloud Run is a fully managed serverless platform which already has an integrated load balance. [See the doc here](https://cloud.google.com/run/docs/about-concurrency)
3. Cloud Run will automatically scale to the number of container instances needed to handle all incoming requests. [See the doc here](https://cloud.google.com/run/docs/about-instance-autoscaling)
4. All request to get a URL will be checked in the cache first and any new URL that is generated will cached with a configurable TTL.
5. All request to get a URL that is not found in the cache will be queried on the NoSQL and any new URL that is generated will be stored on the NoSQL.
6. For eache redirect request, a message will be sent to a pubsub topic to record that we had one more access.
7. Here, we have an Apache Beam pipeline running in Dataflow, to group all messages by Id in a fixed time window, that will update the NoSQL database. [This is a separate project, see here](https://github.com/erickhgm/url-shortener-counter)

## Installing / Getting started

### **Using `docker-compose`**

In the terminal run the following command:
```console
docker-compose up
``` 

### **Access the documentation**
After starting the app you can access the documentation and test using the `Try it on` option.

http://localhost:8090/doc/index.html


## Running tests

In the terminal run the following command:

```console
go test -coverprofile=coverage.out -v ./...
```

Show the coverage report on terminal:
```console
go tool cover -func=coverage.out
```

Open the coverage report in HTML format:
```console
go tool cover -html=coverage.out
```

## Public access via Google Cloud Run
https://url-shortener-ztiqwvbfiq-rj.a.run.app/doc/

## Performance testing using JMeter

20 users making 500 parallel requests during 1min 19sec, total 10,000 requests:

![response-time-graph](doc/response-time-graph.png?raw=true "response-time-graph")

CPU utilization:

![cpu-utilization](doc/cpu-utilization.png?raw=true "cpu-utilization")

Memory utilization:

![memory-utilization](doc/memory-utilization.png?raw=true "memory-utilization")

