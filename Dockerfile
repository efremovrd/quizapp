FROM golang:1.18.1

RUN groupadd --gid 5000 minotauro \
&& useradd --home-dir /home/minotauro --create-home --uid 5000 \
--gid 5000 --shell /bin/sh --skel /dev/null minotauro

RUN mkdir /home/minotauro/quizapp

RUN chown minotauro /home/minotauro/quizapp

USER minotauro

COPY ./ /home/minotauro/quizapp/

RUN export GOPATH=/home/minotauro/quizapp

WORKDIR /home/minotauro/quizapp

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7

RUN swag init -d cmd/api,internal

RUN go build cmd/api/main.go

ENTRYPOINT [ "./main" ]
