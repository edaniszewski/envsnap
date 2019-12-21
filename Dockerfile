FROM scratch

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="edaniszewski/envsnap" \
      org.label-schema.vcs-url="https://github.com/edaniszewski/envsnap" \
      org.label-schema.vendor="Erick Daniszewski"

ADD envsnap /bin

WORKDIR envsnap

ENTRYPOINT ["/bin/envsnap"]
