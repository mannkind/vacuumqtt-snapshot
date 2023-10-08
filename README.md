# MQTT Snapshot

Post images from a rooted L10S Ultra (maybe other things) to MQTT.

On Dev Machine:
```sh
GOOS=linux GOARCH=arm64 go build
scp -O mqtt-snapshot root@VACUUMIP:/data

```

On Robot:
```
cat > /data/mqtt-snapshot.sh << EOF
while true; do
	/data/mqtt-snapshot send-latest --broker mqtt.lan:1883
	sleep 37
done
EOF

chmod +x /data/mqtt-snapshot.sh
cat >> /data/_root_postboot.sh << EOF
if [[ -f /data/mqtt-snapshot.sh ]]; then
	/data/mqtt-snapshot.sh > /dev/null 2>&1 &
fi
EOF

```