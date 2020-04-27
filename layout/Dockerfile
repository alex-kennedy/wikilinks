# Image which runs the LGL algorithm to produce a layout. The pagelinks.lgl file
# should be mounted at /input, and the output location to /output.

FROM alpine

ENV CPLUS_INCLUDE_PATH=/boost/boost_1_55_0

RUN apk update && apk upgrade && \
  apk add --no-cache bash git tar perl make g++ && \
  git clone https://github.com/TheOpteProject/LGL && cd LGL && git checkout 6ecc54f && \
  mkdir /boost && cd /boost && wget -qO- http://sourceforge.net/projects/boost/files/boost/1.55.0/boost_1_55_0.tar.gz | tar -xzf-

# Compiles the LGL engine.
WORKDIR /LGL
RUN ./setup.pl -i

# Runs the layout engine using the layout.conf file.
COPY . /config
ENV PERL5LIB=/LGL/perls
CMD ["/LGL/bin/lgl.pl", "-c", "/config/layout.conf"]
