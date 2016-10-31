FROM scratch

MAINTAINER Robert Jacob <xperimental@solidproject.de>

COPY knwifimap /knwifimap
COPY web /web
COPY database.sqlite /database.sqlite

EXPOSE 8080
CMD [ "/knwifimap", "-a", ":8080", "-f", "/database.sqlite" ]
