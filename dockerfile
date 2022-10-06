FROM alpine

COPY main /barkbridge

CMD ["barkbridge", "server"]