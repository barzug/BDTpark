FROM ubuntu:16.04

MAINTAINER Sergey Barsukov

# Обновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql
#
ENV PGVER 9.6

RUN apt-get install -y wget

RUN echo deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main > /etc/apt/sources.list.d/pgdg.list

RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | \
         apt-key add -

RUN apt-get -y update



RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `postgres` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql -c "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    psql -c "GRANT ALL ON DATABASE postgres TO docker;" &&\
    psql -d postgres -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >>\
 /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "max_wal_size = 1GB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 128MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 64MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

#
# Сборка проекта
#

# Установка golang
RUN wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.8.linux-amd64.tar.gz
RUN apt-get install -y git


# Выставляем переменную окружения для сборки проекта
ENV GOPATH /opt/go
#ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin


#Устанавливаем требуемые пакеты
RUN go get -v github.com/jackc/pgx
RUN go get -v github.com/qiangxue/fasthttp-routing
RUN go get -v github.com/valyala/fasthttp
RUN go get -v github.com/fasthttp-contrib/render


ADD / /

EXPOSE 5000

USER postgres
CMD service postgresql start && go run main.go