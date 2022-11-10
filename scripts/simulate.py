import argparse
import base64
import datetime
import math
import io
import time
from kafka import KafkaProducer
from faker import Faker
import random

# this list of topics should be inclued in docker-compose file kafka service section.
TOPICS = [
    'kv-topic-01',
    'kv-topic-02',
    'kv-topic-03'
]

# this topic should be inclued in docker-compose file kafka service section.
KAFKA_OUT_TOPIC = 'karane-intopic'

def partition_logs(log_lines, package_amount):
    logs_per_package = math.ceil(len(log_lines) / package_amount)
    return [log_lines[i:i + logs_per_package] for i in range(0, len(log_lines), logs_per_package)]


def generate_simple_messages(number):
    fake = Faker()
    num_topics = len(TOPICS) - 1
    lines = []
    for _ in range(number):
        topic = TOPICS[random.randint(0,num_topics)]
        name = fake.name()
        line = '%s\t%s' % (topic, name)
        lines.append(line)

    return lines



if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("events", help="Number of events the script will send", type=int)
    parser.add_argument("interval", help="Interval between log shipping", type=float)
    parser.add_argument("--bootstrap_servers", help="The Kafka servers to produce logs", type=str, default="kafka:9092")
    parser.add_argument("--package_amount", help="The amount of packages that will be sent", type=int)
    
    args = parser.parse_args()

    print("Starting...")

    producer = KafkaProducer(bootstrap_servers=args.bootstrap_servers.split(",")
                             )

    print("Producer created")

    # create all the log lines that will be sent
    log_lines = generate_simple_messages(args.events)
    print("LogLines:\n" + str(log_lines))

    # create a package for each log if package_amount is not configured
    package_amount = args.package_amount if args.package_amount else args.events
    partitioned_log_lines = partition_logs(log_lines, package_amount)
    print("PartitionedLines:\n" + str(partitioned_log_lines))

    for log_lines in partitioned_log_lines:
        # send data
        producer.flush()
        message = "\n".join(log_lines)
        print("KAFKA_OUT_TOPIC: " + KAFKA_OUT_TOPIC)
        print("message: " + message)

        # future = producer.send(KAFKA_OUT_TOPIC, msgpack.packb(message))
        future = producer.send(KAFKA_OUT_TOPIC, message.encode('utf-8'))

        # get metadata
        record_metadata = future.get()

        print("Message sent to partition {}:".format(record_metadata.partition))
        print(message)

        time.sleep(args.interval)
        # block until all async messages are sent

    print("Finished")
    producer.flush()

