RUN go mod download

COPY . .

RUN go build -o main .
expose 3000

CMD["./main"]