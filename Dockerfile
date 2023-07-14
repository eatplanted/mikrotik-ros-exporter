# ------------- BUILD --------------- #
FROM golang:1.20 as build

RUN mkdir -p /src/build
WORKDIR /src/build

COPY . .

RUN make build

# -------------- RUN ---------------- #
FROM scratch

COPY --from=build /src/build/dist/mikrotik-ros-exporter ./

EXPOSE 9142

ENTRYPOINT ["./mikrotik-ros-exporter"]
CMD        [ "-config=/config.yml"]
