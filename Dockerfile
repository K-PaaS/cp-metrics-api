FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum main.go ./

COPY action ./action
COPY api ./api
COPY app ./app
COPY common ./common
COPY docs ./docs
COPY model ./model
COPY yaml ./yaml
COPY *.properties ./

RUN ls -al

RUN go mod download

RUN go mod tidy

#COPY --from=itinance/swag /root/swag /usr/local/bin

#RUN swag init

RUN go build -o main .

# Check if properties files exist, if not create metric.properties
RUN if [ -z "$(ls /build/*.properties 2>/dev/null)" ]; then touch /build/metric.properties; fi

WORKDIR /dist

RUN cp /build/main .

RUN cp /build/*.properties .

FROM scratch

COPY --from=builder /dist/main .

COPY --from=builder /dist/*.properties .

ENV PROFILE=prod \
    DATABASE_URL=${DATABASE_URL} \
    DATABASE_NAME=cp \
    DATABASE_TERRAMAN_ID=${DATABASE_TERRAMAN_ID} \
    DATABASE_TERRAMAN_PASSWORD=${DATABASE_TERRAMAN_PASSWORD} \
    VAULT_URL=${VAULT_URL} \
    VAULT_ROLE_ID=${VAULT_ROLE_ID} \
    VAULT_SECRET_ID=${VAULT_SECRET_ID}

ENTRYPOINT ["/main"]
