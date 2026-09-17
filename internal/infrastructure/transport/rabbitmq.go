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

	conn     *rabbitmqamqp.AmqpConnection
	consumer *rabbitmqamqp.Consumer
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

	consumer, err := conn.NewConsumer(
		ctx,
		sendEmailQueue,
		nil,
	)
	if err != nil {
		_ = conn.Close(ctx)
		s.conn = nil
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

		req := new(v1.SendEmailRequest)

		if err := proto.Unmarshal(
			delivery.Message().GetData(),
			req,
		); err != nil {
			if err := delivery.Discard(ctx, nil); err != nil {
				return err
			}

			continue
		}

		if _, err := s.service.SendEmail(ctx, req); err != nil {
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
