FROM golang:1.22.2-alpine AS start

WORKDIR /Forum/

RUN apk add gcc musl-dev

COPY . .

RUN go build -o forum .

FROM alpine

WORKDIR /myProject

COPY --from=start /Forum/forum /myProject/forum
COPY --from=start /Forum/static /myProject/static
COPY --from=start /Forum/templates /myProject/templates
COPY --from=start /Forum/database.db /myProject/database.db


LABEL version="0.0.1"
LABEL projectname="FORUM"

CMD ["./forum"]