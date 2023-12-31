# This is debian:bookworm-slim (i.e. Debian 12, testing)
FROM debian:bookworm-slim@sha256:9bd077d2f77c754f4f7f5ee9e6ded9ff1dff92c6dce877754da21b917c122c77

ARG TARGETOS
ARG TARGETARCH

# TODO(negz): Find a better way to get an OCI runtime? Ideally we'd grab a
# static build of crun (or runc) that we could drop into a distroless image. We
# slightly prefer crun for its nascent WASM and KVM capabilities, but they only
# offer static builds for amd64 and arm64 and building our own takes a long
# time.
RUN apt-get update && apt-get install -y ca-certificates crun && rm -rf /var/lib/apt/lists/*

COPY bin/${TARGETOS}\_${TARGETARCH}/function-runtime-exec /usr/local/bin/

# We run xfn as root in order to grant it CAP_SETUID and CAP_SETGID, which are
# required in order to create a user namespace with more than one available UID
# and GID. xfn invokes all of the logic that actually fetches, caches, and runs
# a container as an unprivileged user (relative to the root/initial user
# namespace - the user is privileged inside the user namespace xfn creates).
# 
# It's possible to run xfn without any root privileges at all - uncomment the
# following line to do so. Note that in this mode xfn will only be able to
# create containers with a single UID and GID (0), so Containerized Functions
# that don't run as root may not work.
# USER 65532

ENTRYPOINT ["function-runtime-exec"]
