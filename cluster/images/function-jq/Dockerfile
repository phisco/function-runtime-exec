FROM phisco/function-runtime-exec:v0.0.2 as runtime
FROM alpine:3.18
RUN apk add --no-cache jq
COPY --from=runtime /usr/local/bin/function-runtime-exec /usr/local/bin/function-runtime-exec
CMD ["function-runtime-exec", "start", "-d", "--", "/bin/sh", "-c", "jq '.desired.composite.resource.labels |= {\"labelizer.xfn.crossplane.io/processed\": \"true\"} + .'"]

