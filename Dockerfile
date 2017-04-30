FROM scratch
MAINTAINER asdine.elhrychy@gmail.com

EXPOSE 4000

COPY pulpe /pulpe
COPY dist /dist

CMD ["/pulpe", "server", "--mongo", "mongodb://mongodb:27017/pulpe"]
