FROM golang:1.16.3 as build-env

WORKDIR /opendata
COPY sentinel/go.mod sentinel/go.sum /opendata/
RUN go mod download && go mod verify

COPY ./sentinel /opendata/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM ipfs/go-ipfs:v0.9.1 as ipfs

# Now comes the actual target image, which aims to be as small as possible.
FROM busybox:1.31.1-glibc
LABEL maintainer="Manu <manu@set.health>"

# Get the ipfs binary, entrypoint script, and TLS CAs from the build container.
ENV SRC_DIR /go-ipfs
COPY --from=ipfs /usr/local/bin/ipfs /usr/local/bin/ipfs
COPY --from=ipfs /usr/local/bin/start_ipfs /usr/local/bin/start_ipfs
COPY --from=ipfs /sbin/su-exec /sbin/su-exec
COPY --from=ipfs /sbin/tini /sbin/tini
COPY --from=ipfs /usr/local/bin/fusermount /usr/local/bin/fusermount
COPY --from=ipfs /etc/ssl/certs /etc/ssl/certs
COPY --from=build-env /opendata/app /usr/local/bin/opendata_sentinel
COPY ./sentinel/init /usr/local/bin/init_server

# Add suid bit on fusermount so it will run properly
RUN chmod 4755 /usr/local/bin/fusermount

# Fix permissions on start_ipfs (ignore the build machine's permissions)
RUN chmod 0755 /usr/local/bin/start_ipfs
RUN chmod 0755 /usr/local/bin/init_server

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=ipfs /lib/libdl.so.2 /lib/

# Copy over SSL libraries.
COPY --from=ipfs /usr/lib/libssl.so* /usr/lib/
COPY --from=ipfs /usr/lib/libcrypto.so* /usr/lib/

# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm UDP; should be exposed to the public
EXPOSE 4001/udp
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
EXPOSE 8080


# Create the fs-repo directory and switch to a non-privileged user.
ENV IPFS_PATH /data/ipfs
ENV IPFS_PROFILE server,badgerds
ENV PINSET_URL https://api.set.health/pinset/openview-health/bu-5643105772503040

RUN mkdir -p $IPFS_PATH \
  && adduser -D -h $IPFS_PATH -u 1000 -G users ipfs \
  && chown ipfs:users $IPFS_PATH

# Create mount points for `ipfs mount` command
RUN mkdir /ipfs /ipns \
  && chown ipfs:users /ipfs /ipns

# Expose the fs-repo as a volume.
# start_ipfs initializes an fs-repo if none is mounted.
# Important this happens after the USER directive so permissions are correct.
VOLUME $IPFS_PATH

# The default logging level
ENV IPFS_LOGGING ""
ENV OPENDATA_NODE "anonymous"

# This just makes sure that:
# 1. There's an fs-repo, and initializes one if there isn't.
# 2. The API and Gateway are accessible from outside the container.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/init_server"]
