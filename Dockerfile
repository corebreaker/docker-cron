FROM docker/compose:1.25.0-rc1
LABEL maintainer="Corebreaker"
LABEL description="Have a cron scheduler for docker which run command in a docker container"
LABEL version="0.1.0"

RUN adduser -D docker-cron && echo 'docker:x:134:docker-cron' >/etc/group

COPY docker-cron entry-point.sh /home/docker-cron/
RUN mkdir /projects && chown docker-cron:docker-cron /projects /home/docker-cron/*

USER docker-cron:docker-cron
WORKDIR /home/docker-cron/
VOLUME /projects

ENTRYPOINT ["/home/docker-cron/entry-point.sh"]
CMD ["start"]
