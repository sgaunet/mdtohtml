FROM golang:1.16.7-alpine AS builder
LABEL stage=builder

RUN apk add --no-cache git upx
ENV GOPATH /go
COPY src/ /go/src/
COPY go.mod /go/src/
COPY go.sum /go/src/
WORKDIR /go/src/

RUN echo $GOPATH
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build . 
RUN upx mdtohtml



FROM alpine:3.12.3 AS final
LABEL maintainer="Sylvain Gaunet <sgaunet@gmail.com>"

RUN addgroup -S mdtohtml_group -g 1000 && adduser -S mdtohtml -G mdtohtml_group --uid 1000

WORKDIR /usr/bin/
COPY --from=builder /go/src/mdtohtml .

USER mdtohtml

CMD [ "/usr/bin/mdtohtml"] 