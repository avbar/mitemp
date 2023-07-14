FROM scratch

ADD ./bin/app /app
ADD ./config.yml /config.yml

CMD ["/app"]