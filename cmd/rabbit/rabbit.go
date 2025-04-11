package rabbit

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/mitchellh/mapstructure"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/smartrecruiters/rabbitr/cmd/commons"
)

const (
	xDeathHeader     = "x-death"
	xDeathTimeHeader = "time"
)

// Consumer gathers consumer RabbitMQ related objects
type Consumer struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Tag        string
	Done       chan error
}

// Publisher gathers publisher RabbitMQ related objects
type Publisher struct {
	Connection               *amqp.Connection
	Channel                  *amqp.Channel
	Tag                      string
	MessagesToPublish        int
	NumberOfPublished        int
	NumberOfConfirmed        int
	AllPublished             bool
	AllPublishedAndConfirmed chan bool
	Publishes                chan uint64
	Confirms                 chan amqp.Confirmation
}

// MessageExtras holds data about additional message headers and properties
type MessageExtras struct {
	Headers    amqp.Table `json:"headers"`
	Properties Properties `json:"properties"`
}

// Properties represents RabbitMQ message properties used during message publishing
type Properties struct {
	AppID           string `json:"app_id,omitempty"`           // application use - creating application
	ContentEncoding string `json:"content_encoding,omitempty"` // MIME content encoding
	ContentType     string `json:"content_type,omitempty"`     // MIME content type
	CorrelationID   string `json:"correlation_id,omitempty"`   // application use - correlation identifier
	DeliveryMode    uint8  `json:"delivery_mode,omitempty"`    // queue implementation use - Transient (1) or Persistent (2)
	Expiration      string `json:"expiration,omitempty"`       // implementation use - message expiration spec
	MessageID       string `json:"message_id,omitempty"`       // application use - message identifier
	Priority        uint8  `json:"priority,omitempty"`         // queue implementation use - 0 to 9
	ReplyTo         string `json:"reply_to,omitempty"`         // application use - address to reply to (ex: RPC)
	Type            string `json:"type,omitempty"`             // application use - message type name
	UserID          string `json:"user_id,omitempty"`          // application use - creating user id
	Timestamp       int64  `json:"timestamp,omitempty"`        // application use - message timestamp
	Exchange        string `json:"exchange,omitempty"`
	RoutingKey      string `json:"routing_key,omitempty"`
}

// GetRabbitClient returns rabbit client initialized with a provided server coordinates
func GetRabbitClient(serverName string) *rabbithole.Client {
	var exists bool
	var coordinates *commons.ServerCoordinates

	cfg, err := commons.GetApplicationConfig(serverName)
	commons.AbortIfError(err)

	if coordinates, exists = cfg.Servers[serverName]; !exists {
		fmt.Printf("Configuration for server %s has not been found, please add it first via: `rabbitr server add` command", serverName)
		os.Exit(1)
	}
	client, err := rabbithole.NewClient(coordinates.APIURL, coordinates.Username, coordinates.Password)
	commons.AbortIfError(err)
	return client
}

// GetAmqpRabbitConnection returns amqp rabbit connection initialized with a provided server coordinates
func GetAmqpRabbitConnection(serverName, vhost string) *amqp.Connection {
	var exists bool
	var coordinates *commons.ServerCoordinates

	cfg := commons.GetCachedApplicationConfig(serverName)

	if coordinates, exists = cfg.Servers[serverName]; !exists {
		fmt.Printf("Configuration for server %s has not been found, please add it first via: `rabbitr server add` command", serverName)
		os.Exit(1)
	}
	amqpURI := constructAmqpURI(coordinates, vhost)
	commons.Debugf("Connecting to %s", amqpURI)

	config := amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Properties: amqp.Table{
			"connection_name":        "rabbitr",
			"consumer_cancel_notify": true,
		},
	}

	if strings.HasPrefix(amqpURI, "amqps://") {
		tlsConfig := new(tls.Config)
		tlsConfig.InsecureSkipVerify = true
		conn, err := amqp.DialTLS(amqpURI, tlsConfig)
		commons.AbortIfError(err)
		setupCloseListener(conn)
		return conn
	}

	conn, err := amqp.DialConfig(amqpURI, config)
	commons.AbortIfError(err)
	setupCloseListener(conn)
	return conn
}

func setupCloseListener(connection *amqp.Connection) {
	go func() {
		commons.Debugf("Closing: %s", <-connection.NotifyClose(make(chan *amqp.Error)))
	}()
}

func constructAmqpURI(coordinates *commons.ServerCoordinates, vhost string) string {
	host := strings.ReplaceAll(coordinates.AmqpURL, "amqp://", fmt.Sprintf("amqp://%s:%s@", coordinates.Username, coordinates.Password))
	host = strings.ReplaceAll(host, "amqps://", fmt.Sprintf("amqps://%s:%s@", coordinates.Username, coordinates.Password))
	return fmt.Sprintf("%s/%s", host, vhost)
}

func getProperties(msg amqp.Delivery) Properties {
	props := Properties{
		AppID:           msg.AppId,
		ContentEncoding: msg.ContentEncoding,
		ContentType:     msg.ContentType,
		DeliveryMode:    msg.DeliveryMode,
		Priority:        msg.Priority,
		CorrelationID:   msg.CorrelationId,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		MessageID:       msg.MessageId,
		Type:            msg.Type,
		UserID:          msg.UserId,
		Exchange:        msg.Exchange,
		RoutingKey:      msg.RoutingKey,
	}

	if !msg.Timestamp.IsZero() {
		props.Timestamp = msg.Timestamp.Unix()
	}

	return props
}

// CreateMessageToPublish constructs message to sent from the provided contents
func CreateMessageToPublish(msgBody []byte, extras MessageExtras) amqp.Publishing {
	msg := amqp.Publishing{
		Body:    msgBody,
		Headers: convertXDeathToAmqpTable(extras.Headers),
	}
	fillMessageProperties(&msg, &extras)
	return msg
}

func fillMessageProperties(msg *amqp.Publishing, extras *MessageExtras) {
	msg.AppId = extras.Properties.AppID
	msg.ContentEncoding = extras.Properties.ContentEncoding
	msg.ContentType = extras.Properties.ContentType
	msg.CorrelationId = extras.Properties.CorrelationID
	msg.DeliveryMode = extras.Properties.DeliveryMode
	msg.Expiration = extras.Properties.Expiration
	msg.MessageId = extras.Properties.MessageID
	msg.Priority = extras.Properties.Priority
	msg.ReplyTo = extras.Properties.ReplyTo
	msg.Type = extras.Properties.Type
	msg.UserId = extras.Properties.UserID
	if msg.ContentEncoding == "" {
		msg.ContentEncoding = "UTF-8"
	}
}

func convertXDeathToAmqpTable(headers amqp.Table) amqp.Table {
	deaths := make([]interface{}, 0)
	for _, v := range headers[xDeathHeader].([]interface{}) {
		var singleDeath amqp.Table

		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &singleDeath,
			TagName:  "json",
		}
		decoder, err := mapstructure.NewDecoder(cfg)
		commons.AbortIfError(err)

		err = decoder.Decode(v)
		deaths = append(deaths, singleDeath)
		commons.AbortIfErrorWithMsg("unable to decode %s", err)
	}
	headers[xDeathHeader] = deaths
	return headers
}

// InitMessageExtras initializes message extras object based on the amqp delivery
func InitMessageExtras(msg amqp.Delivery) MessageExtras {
	return MessageExtras{
		Headers:    convertTimeHeadersToUnixTime(msg.Headers),
		Properties: getProperties(msg),
	}
}

func convertTimeHeadersToUnixTime(headers amqp.Table) amqp.Table {
	if headers == nil || headers[xDeathHeader] == nil {
		return headers
	}
	xDeaths := headers[xDeathHeader].([]interface{})
	for _, v := range xDeaths {
		singleDeath := v.(amqp.Table)
		if singleDeath[xDeathTimeHeader] != nil {
			_, ok := singleDeath[xDeathTimeHeader].(time.Time)
			if ok {
				unixTs := singleDeath[xDeathTimeHeader].(time.Time).Unix()
				commons.Debugf("Converting xDeath time: %s to unix ts: %d", singleDeath["time"], unixTs)
				singleDeath[xDeathTimeHeader] = unixTs
			}
		}
	}
	return headers
}

// InitConsumer initializes consumer object
func InitConsumer() *Consumer {
	return &Consumer{
		Connection: nil,
		Channel:    nil,
		Tag:        "rabbitr",
		Done:       make(chan error),
	}
}

// InitPublisher initializes publisher object
func InitPublisher() *Publisher {
	return &Publisher{
		Connection:               nil,
		Channel:                  nil,
		Tag:                      "rabbitr",
		AllPublishedAndConfirmed: make(chan bool),
		Publishes:                make(chan uint64, 8), // We'll allow for a few outstanding publisher confirms
	}
}

// SetupCloseHandler setups consumer close handler
func SetupCloseHandler(consumer *Consumer, verbose bool) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		commons.PrintIfTrue(verbose, "Ctrl+C pressed in Terminal")
		if err := consumer.Shutdown(verbose); err != nil {
			commons.Debugf("error during shutdown: %s", err)
		}
		os.Exit(0)
	}()
}

// Shutdown performs publisher shutdown
func (p *Publisher) Shutdown(verbose bool) error {
	time.Sleep(time.Second)
	if err := p.Channel.Close(); err != nil {
		return fmt.Errorf("publisher channel close failed: %s", err)
	}

	if err := p.Connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer func() {
		commons.PrintIfTrue(verbose, "AMQP shutdown OK")
	}()

	return nil
}

// Shutdown performs consumer shutdown
func (c *Consumer) Shutdown(verbose bool) error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.Tag, false); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.Connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer func() {
		commons.PrintIfTrue(verbose, "AMQP shutdown OK")
	}()

	// wait for handle() to exit
	return <-c.Done
}
