FROM alpine

COPY image-sync-server /opt/image-sync-server
COPY config.yaml /opt/config.yaml

WORKDIR /opt

ENTRYPOINT [ "/opt/image-sync-server" ]
CMD [ "-f","/opt/config.yaml" ]