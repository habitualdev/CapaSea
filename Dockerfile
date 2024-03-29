FROM ubuntu:21.04

LABEL maintaner = "https://github.com/habitualdev"

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install wget ca-certificates -y \
    && useradd -m capa

WORKDIR /home/capa

RUN wget "https://github.com/habitualdev/CapaSea/releases/download/v0.5.0/CapaSea"

COPY ./init.sh /home/capa/init.sh

CMD ["/home/capa/init.sh"]




