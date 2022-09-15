FROM golang
WORKDIR /app
COPY . .
RUN cd cmd/formdress && go build
EXPOSE 8000
CMD ["/app/cmd/formdress/formdress", "-d", "/app/docs", "-l", "0.0.0.0:8000"]
