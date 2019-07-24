FROM ubuntu
MAINTAINER Konrad Reiche <konrad@october.news>

COPY activityrecorder /
COPY config.toml /

CMD ["/activityrecorder", "record", "--config", "config.toml"]
