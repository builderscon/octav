FROM debian

ENV PATH=/local-lib/bin:/opt/perl/bin:$PATH
RUN apt-get update && \
    apt-get install -qy --no-install-recommends libssl-dev 
ADD perl.tar.gz /
ADD cpanfile /
ADD adminweb.tar.gz /
ADD Octav.pm /adminweb/lib/WebService/Octav.pm
ADD adminweb-local-lib.tar.gz /
ENV PERL5OPT="-Mlib=/local-lib/lib/perl5 -Mlib=/adminweb/lib"
CMD [ "start_server", "--port", "5000", "--signal-on-hup=USR1", "/adminweb/script/run-plackup.sh" ]
