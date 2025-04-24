FROM postgres:15

COPY up.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["postgres"]
