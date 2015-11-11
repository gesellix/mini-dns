FROM alpine:edge
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/mini-dns
COPY . $APPPATH

RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc \
    && cd $APPPATH && go get -d && go build -o /bin/mini-dns \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mini-dns" ]
CMD [ "-port=5353", "-printf=true" , "-debug=true" ]
