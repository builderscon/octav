FROM debian

ENV PATH=/local-lib/bin:/opt/perl/bin:$PATH
ADD cpanfile /
RUN apt-get update && \
    apt-get install -qy --no-install-recommends wget build-essential ca-certificates libssl-dev && \
    wget http://www.cpan.org/src/5.0/perl-5.20.3.tar.gz -O - | tar xz && \
    cd perl-5.20.3 && \
    ./Configure -de -Dprefix=/opt/perl -Dman1dir=none -Dman3dir=none -Doptimize='-O3' && \
    make all install distclean && \
    cd .. && \
    rm -rf perl-5.20.3 && \
    wget -O cpanm http://cpanmin.us && \
    chmod +x cpanm && \
    ./cpanm -L local-lib -n --installdeps . && \
    apt-get autoremove -qy wget ca-certificates build-essential

ADD adminweb.tar.gz /
ADD Octav.pm /adminweb/lib/WebService/Octav.pm
ENV PERL5OPT="-Mlib=/local-lib/lib/perl5 -Mlib=/adminweb/lib"
RUN apt-get install -qy --no-install-recommends build-essential && \
    cd /adminweb && /cpanm -L /local-lib -n --installdeps . && \
    apt-get autoremove -qy build-essential

CMD start_server --port 5000 --signal-on-hup=USR1 /adminweb/script/run-plackup.sh
