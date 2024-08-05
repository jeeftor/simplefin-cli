FROM scratch
COPY simplefin-cli /usr/bin/sfcli
ENTRYPOINT ["/usr/bin/sfcli"]