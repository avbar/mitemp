FROM golang:1.20-alpine as build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app github.com/avbar/mitemp/cmd/app

FROM scratch

COPY --from=build /src/bin/app /app
COPY --from=build /src/config.yml /config.yml

CMD ["/app"]