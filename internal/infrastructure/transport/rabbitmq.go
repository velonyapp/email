package transport

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/notification/gen/api/v1"
	"github.com/velonyapp/notification/internal/conf"
	"github.com/velonyapp/notification/internal/presentation/api"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/proto"
)

const (
	sendEmailQueue = "notification.send_email"
)

type RabbitMQServer struct {
	address string
	service *api.Service

	conn      *rabbitmqamqp.AmqpConnection
	consumers []*rabbitmqamqp.Consumer
}

func NewRabbitMQServer(
	c *conf.Transport,
	service *api.Service,
) *RabbitMQServer {
	return &RabbitMQServer{
		address: c.Rabbitmq.Address,
		service: service,
	}
}

func (s *RabbitMQServer) Start(ctx context.Context) error {
	conn, err := rabbitmqamqp.Dial(ctx, s.address, nil)
	if err != nil {
		return err
	}
	s.conn = conn

	emailConsumer, err := conn.NewConsumer(ctx, sendEmailQueue, nil)
	if err != nil {
		return err
	}

	s.consumers = append(s.consumers, emailConsumer)

	return consume(
		ctx,
		emailConsumer,
		func() *v1.SendEmailRequest {
			return new(v1.SendEmailRequest)
		},
		func(ctx context.Context, req *v1.SendEmailRequest) error {
			_, err := s.service.SendEmail(ctx, req)
			return err
		},
	)
}

func consume[T proto.Message](
	ctx context.Context,
	consumer *rabbitmqamqp.Consumer,
	newMessage func() T,
	handler func(context.Context, T) error,
) error {
	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return nil
			}

			return err
		}

		req := newMessage()

		if err := proto.Unmarshal(
			delivery.Message().GetData(),
			req,
		); err != nil {
			if err := delivery.Discard(ctx, nil); err != nil {
				return err
			}

			continue
		}

		if err := handler(ctx, req); err != nil {
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

func (s *RabbitMQServer) Stop(ctx context.Context) error {
	for _, consumer := range s.consumers {
		if err := consumer.Close(ctx); err != nil {
			return err
		}
	}

	s.consumers = nil

	if s.conn != nil {
		if err := s.conn.Close(ctx); err != nil {
			return err
		}

		s.conn = nil
	}

	return nil
}
