FROM golang:alpine as builder

WORKDIR /src/app
COPY . .

RUN apk --no-cache add git

RUN go get -d -v
RUN go build \
    && ls -l

FROM alpine

WORKDIR /app

COPY --from=builder /src/app/app /app/trackjs-exporter

EXPOSE 9197

CMD [ "/app/trackjs-exporter" ]
