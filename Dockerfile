FROM ubuntu:19.04
ENV DEBIAN_FRONTEND noninteractive

# Installing packets
RUN apt-get update && apt-get upgrade -y && apt-get install -y gnupg git postgresql-11 postgresql-contrib

# Cloning project
USER root
RUN git clone https://github.com/lisa-bella97/tech-db-forum.git
WORKDIR tech-db-forum

# Starting PostgreSQL and creating a database
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    psql forum -a -f database/init.sql &&\
    /etc/init.d/postgresql stop