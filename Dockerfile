FROM alpine:edge AS build

RUN apk update
RUN apk upgrade
RUN apk add --update go=1.12.1-r0 gcc=8.3.0-r0 g++=8.3.0-r0 make git

ADD . $GOPATH/src/github.com/VideoCoin/cli
WORKDIR $GOPATH/src/github.com/VideoCoin/cli

# make sure you run make dep before to make this image
RUN make
RUN cp -r build /

FROM alpine:edge

COPY --from=build /build /home/
WORKDIR /home

ENTRYPOINT ["./cli"]