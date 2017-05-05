FROM scratch
MAINTAINER asdine.elhrychy@gmail.com

EXPOSE 4000

COPY dist/ /dist
COPY dist/pulpe /pulpe

CMD ["/pulpe", "server", "--mongo", "mongodb://mongodb:27017/pulpe"]
