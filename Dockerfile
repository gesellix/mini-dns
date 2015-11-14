FROM alpine:edge
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/mini-dns
COPY . $APPPATH

RUN echo '@community http://dl-4.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories
RUN apk upgrade --update --available && apk add go@community git mercurial libc-dev gcc libgcc
RUN cd $APPPATH && go get -d && go build -o /bin/mini-dns
#RUN cd $APPPATH && go get -d && go build -o /bin/mini-dns \
#    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mini-dns" ]
CMD [ "-port=5353", "-printf=true" , "-debug=true" ]
