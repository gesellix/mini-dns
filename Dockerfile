FROM alpine:edge
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/mini-dns

ENV ADD_PACKAGES git mercurial libc-dev gcc libgcc go@community
ENV DEL_PACKAGES git mercurial libc-dev gcc libgcc go

COPY . $APPPATH

# see https://golang.org/pkg/net/?m=all#hdr-Name_Resolution
#ENV GODEBUG netdns=1
RUN echo '@community http://dl-4.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories \
    && apk upgrade --update --available \
    && apk add $ADD_PACKAGES \
    && cd $APPPATH && go get -d && go build -o /bin/mini-dns \
    && apk del --purge $DEL_PACKAGES \
    && rm -rf /var/cache/apk/* && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mini-dns" ]
CMD [ "-port=5353" ]
