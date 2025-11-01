FROM alpine:3.22.2 AS final
RUN addgroup -S mdtohtml_group -g 1000 && adduser -S mdtohtml -G mdtohtml_group --uid 1000
WORKDIR /usr/bin/
COPY mdtohtml .
USER mdtohtml
CMD [ "/usr/bin/mdtohtml"] 