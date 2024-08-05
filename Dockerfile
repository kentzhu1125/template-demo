FROM harbor.galaksiodatanubo.work/base/centos:7.9.2009
WORKDIR /
COPY app app
RUN chmod +X /app
ENTRYPOINT ["/app"]
EXPOSE 8888