ARG DOCKER_REPO=nexus.tools.devopenocean.studio

############################
# STEP 1 build executable binary
############################

FROM ${DOCKER_REPO}/golang-build:builder-v3 as builder

WORKDIR /app
COPY . ./


RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ./bin/svc

############################
# STEP 2 build a small image
############################

FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /app/bin/svc /svc
COPY --from=builder /app/api/service.proto /
COPY --from=builder /app/sha /
COPY --from=builder /app/version /

# Port on which the service will be exposed.
EXPOSE 8080 8888 9100

# Run the svc binary.
CMD ["./svc"]
