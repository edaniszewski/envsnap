FROM scratch

ADD envsnap /bin

WORKDIR envsnap

ENTRYPOINT ["/bin/envsnap"]