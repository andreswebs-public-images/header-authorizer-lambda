FROM golang:1.22 as build
ARG GOARCH=arm64
WORKDIR /app
COPY go.mod go.sum ./
COPY main.go .
RUN GOOS=linux CGO_ENABLED=0 go build -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
ARG GOARCH=arm64
COPY --from=build /app/main /main
ENTRYPOINT ["/main"]
