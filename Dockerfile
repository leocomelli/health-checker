FROM leocomelli/oracle-client:11.2

COPY dist/health-checker /usr/local/bin/health-checker
COPY health.yml /usr/local/lib/health.yml

ENV HC_FILE /usr/local/lib/health.yml
ENV ORACLE_HOME /usr/lib/oracle/11.2/client64
ENV LD_LIBRARY_PATH /usr/local/lib/:/usr/lib/oracle/11.2/client64/lib/

ENTRYPOINT [ "/usr/local/bin/health-checker" ]