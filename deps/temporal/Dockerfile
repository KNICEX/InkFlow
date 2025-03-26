FROM alpine:3.21.3
ARG TARFILE=temporal_cli_latest_linux_amd64.tar.gz
RUN wget "https://temporal.download/cli/archive/latest?platform=linux&arch=amd64" -O $TARFILE
RUN tar -xzf $TARFILE
RUN rm $TARFILE
RUN mv temporal /usr/local/bin/temporal
EXPOSE 7233 8233
CMD ["temporal", "server", "start-dev", "--db-filename", "temporal.db", "--ip","0.0.0.0"]