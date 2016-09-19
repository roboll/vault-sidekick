FROM scratch

ADD bin/vault-sidekick /vault-sidekick
ENTRYPOINT [ "/vault-sidekick" ]
