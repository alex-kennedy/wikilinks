FROM golang:1.13 as base

COPY . /app

# Preparation image (downloads, processes files)
FROM base as prep
WORKDIR /app
RUN cd prep && go build .
EXPOSE 8080
CMD ["/app/prep/prep"]

