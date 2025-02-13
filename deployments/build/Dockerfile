FROM golang:1.22 AS build

RUN apt-get update && apt-get install -y \
     gcc \
     libc-dev \
     git \
     && rm -rf /var/lib/apt/lists/*

ENV BIN_FILE=/opt/calendar/calendar-app
ENV CODE_DIR=/go/src/app/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

FROM debian:latest

RUN apt-get update && apt-get install -y \
    curl \
    postgresql-client \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /opt/calendar/calendar-app /opt/calendar/calendar-app
COPY --from=build /go/bin/grpcurl /usr/local/bin/grpcurl
COPY --from=build /go/bin/goose /usr/local/bin/goose

ENV CONFIG_FILE=/app/configs/calendar_config.yaml
COPY configs/calendar_config.yaml ${CONFIG_FILE}

COPY ./migrations /migrations

COPY run_migrations.sh /app/run_migrations.sh
RUN chmod +x /app/run_migrations.sh

ENTRYPOINT ["/bin/sh", "-c", "/app/run_migrations.sh && /opt/calendar/calendar-app -config /app/configs/calendar_config.yaml"]

CMD ["sh", "-c", "until nc -z db 5432; do echo 'Waiting for database...'; sleep 2; done; /opt/calendar/calendar-app -config /app/configs/calendar_config.yaml"]
