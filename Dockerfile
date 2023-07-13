FROM ubuntu:22.04

RUN apt update
RUN apt install -y bluez bluetooth usbutils

ADD ./bin/app /app

CMD ["/app"]