FROM ubuntu:18.04

RUN apt-get update && apt-get install -y gnupg git

# Клонируем проект
USER root
RUN git clone https://github.com/lisa-bella97/tech-db-forum.git
WORKDIR tech-db-forum

RUN apt-get install postgresql-11

# Подключаемся к PostgreSQL и создаем БД
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    psql forum -a -f database/init.sql &&\
    /etc/init.d/postgresql stop