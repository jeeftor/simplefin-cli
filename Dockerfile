# Dockerfile
FROM scratch
COPY ./dist/{{ .Os }}_{{ .Arch }}/simplefin-cli /mybin
ENTRYPOINT ["/mybin"]