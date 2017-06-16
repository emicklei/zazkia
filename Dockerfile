FROM alpine:3.6

COPY zazkia /

# Web UI
EXPOSE 9191

# useable by routes if port mapping is possible
EXPOSE 9200-9300

# mysql
EXPOSE 3306

# postgres
EXPOSE 5432

# tomcat, jboss
EXPOSE 8080

CMD ./zazkia -v -f /data/zazkia-routes.json