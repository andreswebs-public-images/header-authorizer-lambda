# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /build
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download
COPY src/ ./

ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 go build -tags lambda.norpc -o app .

FROM public.ecr.aws/lambda/provided:al2023 AS runtime
COPY --from=build /build/app /app
ENTRYPOINT [ "/app" ]
