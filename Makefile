broker="localhost:29092"

send_messages:
	@ docker-compose exec kafka /opt/kafka/bin/kafka-console-producer.sh --topic ${topic} --bootstrap-server ${broker}

listen_messages:
	docker-compose exec kafka /opt/kafka/bin/kafka-console-consumer.sh --topic ${topic} --from-beginning --bootstrap-server ${broker}

create_topic: 
	@ docker-compose exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic ${topic} --bootstrap-server ${broker}

list_topics:
	@ docker-compose exec kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server=${broker} --list

target=./...
test:
	@ go test ${target}

.PHONY: build_and_push_image
build_and_push_image: ## Build the forwarder image and push to Digital Ocean registry
ifeq ($(env), stage)
	@ $([ -z "$version" ] && $(eval version="stage"))
# else ifeq ($(env), production)
# 	@ $([ -z "$version" ] && $(eval version="latest"))
else
	@ echo "Invalid environment."; exit 1;
endif
	@ docker build --build-arg env=${env} -t registry.digitalocean.com/azion/karane-forwarder:${version} -f Dockerfile .
	@ docker push registry.digitalocean.com/azion/karane-forwarder:${version}

.PHONY: build_scripts
build_scripts: ## Build the image used to run all the scripts
	@ docker build -t scripts ./scripts

events = 5
period = 1
amount = 1
send_logs:
	# Add lines of log in kafka topic
	@ docker run \
		--privileged \
		-v ${PWD}/scripts:/app \
        -w /app \
		--network="karane-network" \
		-it \
		scripts python3 simulate.py ${events} ${period} --package_amount ${amount} \
		--bootstrap_servers=${broker}
