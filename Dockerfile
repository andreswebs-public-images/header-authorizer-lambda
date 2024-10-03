FROM golang:1.23 AS build

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"
ARG GOPROXY="https://proxy.golang.org"

ENV GOARCH="${TARGETARCH}"
ENV GOOS="${TARGETOS}"
ENV GOPROXY="${GOPROXY}"

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go


FROM public.ecr.aws/lambda/provided:al2023 AS runtime

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"

ENV GOARCH="${TARGETARCH}"
ENV GOOS="${TARGETOS}"

COPY --from=build "/app/bootstrap" "/bootstrap"
ENTRYPOINT ["/bootstrap"]
