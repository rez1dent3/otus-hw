# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_FILE /bin/calendar_scheduler
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

ENV BIN_FILE /bin/calendar_scheduler
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE} -config /etc/config.yaml
