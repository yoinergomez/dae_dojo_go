FROM golang

ARG app_env
ENV APP_ENV $app_env

WORKDIR /go/src/github.com/camiloperezv/dojo_go/app
COPY ./app .

RUN go get ./ \
	&& go get -u github.com/gorilla/mux \
	&& go get gopkg.in/mgo.v2 \
	&& go get gopkg.in/mgo.v2/bson \
	&& go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi

EXPOSE 8080