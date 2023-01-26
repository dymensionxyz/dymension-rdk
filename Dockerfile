# ---------------------------------------------------------------------------- #
#                               Dymension builder                              #
# ---------------------------------------------------------------------------- #

FROM golang:1.18-bullseye as dymension-builder

ENV PACKAGES build-essential curl make git bash gcc wget
RUN apt-get update -y
RUN apt-get install -y $PACKAGES

WORKDIR /app

#Build dymd
RUN git clone https://github.com/dymensionxyz/dymension && cd dymension && make build

#Build rollappd
WORKDIR /app/dymension-rdk
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build


# ---------------------------------------------------------------------------- #
#                                    Rollapp                                   #
# ---------------------------------------------------------------------------- #

FROM debian:bullseye
COPY --from=dymension-builder /app/dymension/bin/dymd /usr/local/bin/
COPY --from=dymension-builder /app/dymension-rdk/rollappd /usr/local/bin/

ENV PACKAGES curl make bash jq
RUN apt-get update -y
RUN apt-get install -y $PACKAGES

WORKDIR /app
COPY scripts/* ./scripts/
RUN chmod +x ./scripts/*.sh

EXPOSE 26656 26667 1317 9090

ENTRYPOINT ["/usr/local/bin/rollappd"]
