# MQTT Snapshot

Post images from a rooted L10S Ultra (maybe other things) to MQTT.


## Setup

### Binary

* Download a release from GitHub or build one locally.
```sh
GOOS=linux GOARCH=arm64 go build
scp -O mqtt-snapshot root@VACUUMIP:/data
```

### Vacuum

```sh
# Run mqtt-snapshot at startup
cat >> /data/_root_postboot.sh << EOF
if [[ -f /data/mqtt-snapshot ]]; then
    /data/mqtt-snapshot send-latest --broker mqtt.lan:1883 > /dev/null 2>&1 &
fi
EOF

# Reboot
reboot
```