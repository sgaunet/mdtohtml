FROM alpine:3.16.0 AS final
LABEL maintainer="Sylvain Gaunet <sgaunet@gmail.com>"

RUN addgroup -S mdtohtml_group -g 1000 && adduser -S mdtohtml -G mdtohtml_group --uid 1000

WORKDIR /usr/bin/
COPY mdtohtml .

USER mdtohtml

CMD [ "/usr/bin/mdtohtml"] 