FROM golang

# if left blank app will run with dev settings
# to build production image run:
# $ docker build ./api --build-args app_env=production
ARG app_env
ENV APP_ENV $app_env

# 
COPY ./src /go/src
WORKDIR /go/src/api

# install libraries
RUN go get github.com/jinzhu/gorm
RUN go get github.com/julienschmidt/httprouter
RUN go get github.com/lib/pq
RUN go get github.com/lib/pq/hstore
RUN go get github.com/pilu/fresh

# after this /go/bin/api binary will apear
RUN go install

# if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
# if production setting will build binary
CMD if test "${NODE_ENV}" = "production" ; then \
	/go/bin/api; \
	else \
	fresh; \
	fi
	
EXPOSE 8080
