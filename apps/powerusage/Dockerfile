FROM debian:bullseye-slim

RUN apt-get update && \
	apt-get install -y curl powertop bc  netcat && \
	apt-get clean

COPY power-exporter.sh /usr/local/bin/power-exporter
RUN chmod +x /usr/local/bin/power-exporter

EXPOSE 9109
CMD ["/usr/local/bin/power-exporter"]
