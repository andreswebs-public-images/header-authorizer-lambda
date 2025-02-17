FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023 AS runtime
COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]
