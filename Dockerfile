
FROM --platform=linux/amd64 golang:1.17-alpine 

ADD . /home

RUN apk add --no-cache bash vim procps
        
WORKDIR /home

# CMD ["go","run","cmd/simple_forwarder/main.go"]
# CMD ["go","run","cmd/buffered_forwarder/main.go"]
CMD ["go","run","cmd/fanout/main.go"]