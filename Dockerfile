FROM golang:1.24 AS build
WORKDIR /build
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download
COPY src/ ./
RUN CGO_ENABLED=0 go build -tags lambda.norpc -o app .

FROM public.ecr.aws/lambda/provided:al2023 AS runtime
COPY --from=build /build/app /app
ENTRYPOINT [ "/app" ]
