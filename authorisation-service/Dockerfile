ARG DOCKER_REPO=nexus.tools.devopenocean.studio

############################
# STEP 1 build executable binary
############################
FROM ${DOCKER_REPO}/golang-build:builder-v5 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o ./bin/svc

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

# Port on which the grpcsrv will be exposed.
EXPOSE 8132 8080 8090 8888 9100

# Run the svc binary.
CMD ["./svc"]
