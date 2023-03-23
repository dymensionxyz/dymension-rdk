# ---------------------------------------------------------------------------- #
#                               Dymension builder                              #
# ---------------------------------------------------------------------------- #

FROM golang:1.18-bullseye as dymension-builder

RUN apt-get update -y && apt-get install -y build-essential

ENV PACKAGES curl make git bash gcc python3 wget

RUN apt-get install -y $PACKAGES

WORKDIR /app

RUN git clone https://github.com/dymensionxyz/dymension && cd dymension && make install

# ---------------------------------------------------------------------------- #
#                                    Rollapp                                   #
# ---------------------------------------------------------------------------- #

FROM golang:1.18-bullseye
COPY --from=dymension-builder /go/bin/dymd /go/bin

RUN apt-get update -y && apt-get install -y build-essential

ENV PACKAGES curl make git bash gcc python3 wget

RUN apt-get install -y $PACKAGES

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN go install ./cmd/rollappd

RUN chmod +x ./scripts/*.sh

EXPOSE 26656 26657 1317 9090
