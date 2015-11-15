#a mini dns server

run:

    docker run --rm -dit -p 5555:5353 gesellix/mini-dns

cross compile the binary:

    docker build -t dns-cross -f Dockerfile.cross .
    docker create --name mini-dns dns-cross -
    docker cp mini-dns:/usr/src/app/mini-dns-linux-amd64 .
    docker rm mini-dns
