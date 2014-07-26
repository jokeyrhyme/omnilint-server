FROM ubuntu:14.04
MAINTAINER Ron Waldon <jokeyrhyme@gmail.com> @jokeyrhyme

RUN apt-get update -y
RUN apt-get install php5-cli -y
#RUN apt-get install nodejs npm -y
#RUN apt-get install ruby2.0 -y


ADD omnilint-server_linux_amd64.tar.gz /opt
WORKDIR /opt

#ENV NEWRELIC_LICENSE ""
#ENV NEWRELIC_NAME ""

ENV PORT 3000
EXPOSE 3000

CMD /opt/omnilint-server_linux_amd64/omnilint-server
