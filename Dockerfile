FROM ubuntu
MAINTAINER Konrad Reiche <konrad@october.news>

ENV PORT 9000

COPY engine /
COPY local.config.toml config.toml
COPY database.yml /
COPY worker/emailsender/templates worker/emailsender/templates

CMD ["/engine", "--config", "config.toml"]
