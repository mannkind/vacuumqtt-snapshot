# VacuuMQTT Snapshot

Post images from a rooted L10S Ultra (maybe other things) to MQTT.


## Setup

### Binary

* Download a release from GitHub or build one locally.
```sh
GOOS=linux GOARCH=arm64 go build
scp -O vacuumqtt-snapshot root@VACUUMIP:/data
```

### Vacuum

```sh
# Run vacuumqtt-snapshot at startup
cat >> /data/_root_postboot.sh << EOF
if [[ -f /data/vacuumqtt-snapshot ]]; then
    /data/vacuumqtt-snapshot send-latest --broker mqtt.lan:1883 --topic whatever/your/topic/is > /dev/null 2>&1 &
fi
EOF

# Reboot
reboot
```