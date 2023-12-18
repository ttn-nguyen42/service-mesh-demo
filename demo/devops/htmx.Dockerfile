FROM alpine:3.18 AS runner

COPY templates templates
COPY configs.json configs.json
COPY main main

CMD [ "./main" ]