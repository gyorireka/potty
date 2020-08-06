FROM golang:alpine

#RUN apk add git

COPY . /go/src/github.com/gyorireka/potty

WORKDIR /go/src/github.com/gyorireka/potty

RUN go build ./cmd/potty-server

FROM alpine

COPY --from=0 /go/src/github.com/gyorireka/potty/potty-server /potty-server
#COPY keys /keys

ENV PORT 8123
ENV HOST 0.0.0.0

CMD [ "/potty-server" ]