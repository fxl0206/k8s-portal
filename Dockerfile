FROM docker.io/jgowdy/centos7-zero 

ADD portal-server /usr/local/bin/portal-server
ADD html /html
ENTRYPOINT ["/usr/local/bin/portal-server"]
CMD ["/html"]




