FROM alpine:3.6

COPY zazkia /

# Web UI
EXPOSE 9191

# useable by routes
EXPOSE 9200-9300

CMD ./zazkia -v -f /data/zazkia-routes.json