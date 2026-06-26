# Stage 1: Build frontend
FROM node:20-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.23-alpine AS go-build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
COPY --from=frontend-build /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# Stage 3: Runtime
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go-build /app/server .
COPY --from=go-build /app/config.example.yaml .
EXPOSE 8080
ENTRYPOINT ["./server"]
