ARG BUILD_BASE=golang:1.17-alpine
ARG ISO_BASE=alpine

FROM ${BUILD_BASE} as build
WORKDIR /go/src/github.com/github.com/oka-is/alice
ENV BUILD_DEPS make
COPY . .
RUN apk add --update --no-cache $BUILD_DEPS
RUN make build\:linux

FROM ${ISO_BASE} as iso
WORKDIR /app
ENV PATH="/app:${PATH}"
ENV ISO_DEPS bash busybox-extras curl
COPY --from=build /go/src/github.com/github.com/oka-is/alice/build/linux ./alice
COPY docker/launch-entrypoint.sh ./
RUN apk add --update --no-cache $ISO_DEPS
RUN addgroup --gid 2019 user && \
    adduser --disabled-password --uid 2019 --ingroup user --gecos user user
RUN chown -R user:user ./
USER user
SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["alice"]
