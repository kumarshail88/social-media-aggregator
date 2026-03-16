# social-media-aggregator

This project comprises 2 services, `ingestor` and `feed`. 


## Ingestor
The function of `ingestor` is to read the posts from Mastodon api/stream
and publish to pub/sub redis (NOT IMPLEMENTED). 

### How to run:
The service runs in 2 modes: stream and api. 
First set the mastodon token in the env.

```
export MASTODON_TOKEN=<token>
```

And then simply run
```shell
cd ingestor
go run . --mode=api
```

The stream mode is also partially functional but the Mastodon streaming api has restricted the events to only publish DELETE events.

## Feed (NOT IMPLEMENTED)
The function of `feed` is to provide an api for the user to get the timeline. This service is not implemented due to lack of time. 
But the idea for this service is cache the user's timeline into Redis and fan out the timeline via an SSE connection. 