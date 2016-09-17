FROM scratch

ADD vault-sidekick /vault-sidekick
ENTRYPOINT [ "/vault-sidekick "]
