FROM ubuntu:14.04
MAINTAINER Ron Waldon <jokeyrhyme@gmail.com> @jokeyrhyme

RUN apt-get update -y
RUN apt-get install php5-cli -y
#RUN apt-get install nodejs npm -y
#RUN apt-get install ruby2.0 -y

WORKDIR /usr/local/bin
RUN apt-get install curl -y
RUN curl -sS https://getcomposer.org/installer | php

COPY php/composer.json /usr/src/php/composer.json
WORKDIR /usr/src/php
RUN composer.phar install
RUN ln -s /usr/src/php/vendor/bin/phpcs /usr/local/bin/phpcs

ADD omnilint-server_linux_amd64.tar.gz /opt
WORKDIR /opt

#ENV NEWRELIC_LICENSE ""
#ENV NEWRELIC_NAME ""

ENV PORT 3000
EXPOSE 3000

CMD /opt/omnilint-server_linux_amd64/omnilint-server
