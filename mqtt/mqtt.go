package mqtt

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func New(broker string, username string, password string) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("could not connect to MQTT broker; %s, %w", broker, token.Error())
	}

	return c, nil
}
