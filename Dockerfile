FROM golang:1.17-alpine AS build

WORKDIR /app
COPY . .
RUN go build -o /app/main main.go
RUN chmod +x /app/main

FROM scratch 
WORKDIR /app
COPY --from=build /app/main .
ENTRYPOINT [ "/app/main" ]