# Kafka - Go Fanout


Kafka - Go Fanout is a project made for streaming [Kafka](https://kafka.apache.org) messages from a topic and distribute them among different Output Kafka Topics according to a field in the message.

---
## How do I start?

### Using Docker

```
# Start container directly
docker-compose up
```
This will start the container using the available project config file. 

### Running on kubernetes

Install and deploy on the local machine a benthos application using the project config file:
```
# Install on stage kubernetes
helm install forwarder .helm/forwarder --namespace <namespace>
```

### Send messages
```
# Build scripts container
make build_scripts

# Send lots of logs
make send_logs events=20000 period=1 amount=5 broker=<kafka broker address>:<kafka broker port>
```

### Create Topics you need
Inside `scripts/simulate.py`, edit the list TOPICS to the topics you would like to distribute messages to. Then, create all them one by one.
```
# Create topic
make create_topic broker=<kafka broker address>:<kafka broker port> topic=topic-name
```

### Listen messages
```
# Remember you should have the scripts container built

# Send lots of logs
make send_logs events=20000 period=1 amount=5 broker=<kafka broker address>:<kafka broker port>
```
