#!/bin/bash

while true; do
  cat <<EOF > /tmp/metrics.prom
# HELP power_exporter_watts CPU power in watts (estimated)
# TYPE power_exporter_watts gauge
EOF

  if [ -f /sys/class/powercap/intel-rapl:0/energy_uj ]; then
    ENERGY1=$(cat /sys/class/powercap/intel-rapl:0/energy_uj)
    sleep 1
    ENERGY2=$(cat /sys/class/powercap/intel-rapl:0/energy_uj)
    DELTA=$((ENERGY2 - ENERGY1))
    POWER=$(echo "$DELTA / 1000000" | bc -l)
  else
    # Fallback using powertop CSV
    powertop --time=1 --csv=/tmp/powertop.csv > /dev/null 2>&1
    POWER=$(grep "Power est." /tmp/powertop.csv | grep W | awk -F',' '{sum += $2} END {print sum}')
  fi

  echo "power_exporter_watts ${POWER:-0}" >> /tmp/metrics.prom
  sleep 5
done &

# Serve metrics
while true; do
  echo -e "HTTP/1.1 200 OK\nContent-Type: text/plain\n\n$(cat /tmp/metrics.prom)" | \
    nc -l -p 9109 -q 1
done
