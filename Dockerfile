FROM scratch

MAINTAINER Robert Jacob <xperimental@solidproject.de>

COPY knwifimap /knwifimap
COPY web /web

VOLUME /data

EXPOSE 8080
CMD [ "/knwifimap", "-a", ":8080", "-f", "/data/database.sqlite" ]
