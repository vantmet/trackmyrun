FROM ubuntu:latest
WORKDIR /opt/tmr
COPY html ./html
COPY tmr .
RUN apt-get update && apt-get -y install ca-certificates
EXPOSE 5000
ENTRYPOINT ["./tmr"]