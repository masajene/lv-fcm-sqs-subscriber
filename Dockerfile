FROM golang:1.21 as build
WORKDIR /app
# Copy dependencies list
COPY go/go.mod go/go.sum ./
# Build with optional lambda.norpc tag
COPY go/ .
# COPY serviceAccountKey.json ./fcm/serviceAccountKey.json
ENV ARCH="arm64"

RUN GOOS=linux GOARCH=${ARCH} CGO_ENABLED=0 go build -tags lambda.norpc -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o main main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]
