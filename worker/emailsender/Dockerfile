FROM ubuntu
MAINTAINER Konrad Reiche <konrad@october.news>

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY emailsender /
COPY config.toml /

CMD ["/emailsender", "--config", "config.toml"]
