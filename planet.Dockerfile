FROM golang:1.16

WORKDIR /src
ENV PATH="/go/bin:${PATH}"

COPY . .
RUN go mod tidy

CMD ["tail", "-f", "/dev/null"]