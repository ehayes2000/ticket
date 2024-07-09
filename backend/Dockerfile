FROM golang:1.22.5  
WORKDIR /service
COPY  . /service/
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y sqlite3 bash

RUN chmod +x scripts/make_db.sh && ./scripts/make_db.sh
RUN go env -w CGO_ENABLED=1
RUN go build .
EXPOSE 1323
CMD ["./backend"]