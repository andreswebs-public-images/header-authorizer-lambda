FROM golang:1.22 AS build

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"
ARG GOPROXY="https://proxy.golang.org"
ARG APP_NAME="main"

ENV GOARCH="${TARGETARCH}"
ENV GOOS="${TARGETOS}"
ENV GOPROXY="${GOPROXY}"

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -o "${APP_NAME}" .

FROM public.ecr.aws/lambda/provided:al2023

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"
ARG APP_NAME="main"

ENV GOARCH="${TARGETARCH}"
ENV GOOS="${TARGETOS}"

COPY --from=build "/app/${APP_NAME}" "/${APP_NAME}"
ENTRYPOINT ["/${APP_NAME}}"]
