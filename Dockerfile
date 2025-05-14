FROM --platform=amd64 golang:1.24 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./main.go

FROM --platform=amd64 alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=verysecretpassword

VOLUME [ "/data" ]

CMD ["./scheduler"]