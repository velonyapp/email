package transport

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/presentation/api"

	"github.com/Azure/go-amqp"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/proto"
)

const (
	emailJobsQueue = "email.jobs"

	sendEmailRoutingKey = "email.send-email"
)

type RabbitMQServer struct {
	address  string
	username string
	password string
	service  *api.Service

	conn     *rabbitmqamqp.AmqpConnection
	consumer *rabbitmqamqp.Consumer
}

func NewRabbitMQServer(
	c *conf.Transport,
	service *api.Service,
) *RabbitMQServer {
	return &RabbitMQServer{
		address:  c.Rabbitmq.Address,
		username: c.Rabbitmq.Username,
		password: c.Rabbitmq.Password,
		service:  service,
	}
}

func (s *RabbitMQServer) Start(ctx context.Context) error {
	conn, err := rabbitmqamqp.Dial(
		ctx,
		s.address,
		&rabbitmqamqp.AmqpConnOptions{
			SASLType: amqp.SASLTypePlain(
				s.username,
				s.password,
			),
		},
	)
	if err != nil {
		return err
	}

	s.conn = conn

	consumer, err := conn.NewConsumer(ctx, emailJobsQueue, nil)
	if err != nil {
		return err
	}

	s.consumer = consumer

	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return nil
			}

			return err
		}

		message := delivery.Message()

		routingKey, ok := message.Annotations["x-routing-key"].(string)
		if !ok {
			if err := delivery.Discard(ctx, nil); err != nil {
				return err
			}

			continue
		}

		if err := s.handle(ctx, routingKey, message.GetData()); err != nil {
			if err := delivery.Requeue(ctx); err != nil {
				return err
			}

			continue
		}

		if err := delivery.Accept(ctx); err != nil {
			return err
		}
	}
}

func (s *RabbitMQServer) handle(
	ctx context.Context,
	routingKey string,
	data []byte,
) error {
	switch routingKey {
	case sendEmailRoutingKey:
		req := new(v1.SendEmailRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return err
		}

		_, err := s.service.SendEmail(ctx, req)
		return err

	default:
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}
}

func (s *RabbitMQServer) Stop(ctx context.Context) error {
	if s.consumer != nil {
		if err := s.consumer.Close(ctx); err != nil {
			return err
		}

		s.consumer = nil
	}

	if s.conn != nil {
		if err := s.conn.Close(ctx); err != nil {
			return err
		}

		s.conn = nil
	}

	return nil
}
