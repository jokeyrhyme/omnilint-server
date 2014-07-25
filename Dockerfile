FROM ubuntu:14.04
MAINTAINER Ron Waldon <jokeyrhyme@gmail.com> @jokeyrhyme

RUN apt-get update -y
RUN apt-get install php5-cli -y
#RUN apt-get install nodejs npm -y
#RUN apt-get install ruby2.0 -y

#ENV NEWRELIC_LICENSE
#ENV NEWRELIC_NAME
