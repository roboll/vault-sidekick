FROM alpine:3.4

ADD bin/vault-sidekick /vault-sidekick
ENTRYPOINT [ "/vault-sidekick" ]
