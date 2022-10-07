FROM alpine

COPY barkbridge /barkbridge

ENTRYPOINT [ "/barkbridge" ]