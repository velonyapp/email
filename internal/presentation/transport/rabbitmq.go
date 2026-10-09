package transport

import (
	"context"
	"errors"
	"sync"

	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/application/command"
	"github.com/velonyapp/email/internal/conf"

	"github.com/Azure/go-amqp"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

const (
	emailTemplateResourcePattern = "emailTemplates/{email_template}"
)

var _ transport.Server = (*RabbitMQConsumer)(nil)

type RabbitMQConsumer struct {
	c          *conf.Transport
	commandBus *command.Bus

	conn      *rabbitmqamqp.AmqpConnection
	consumers map[string]*rabbitmqamqp.Consumer

	wg sync.WaitGroup
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	commandBus *command.Bus,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		c:          c,
		commandBus: commandBus,
	}
}

func (rc *RabbitMQConsumer) registerAllConsumers(ctx context.Context) error {
	var errs []error

	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queues.SendEmail); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	switch queue {
	case rc.c.Rabbitmq.Queues.SendEmail:
		req := new(v1.SendEmailRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		if content := req.Email.GetCustom(); content != nil {
			if _, err := command.Send(ctx, rc.commandBus, command.SendEmail{
				From:           req.Email.From,
				To:             req.Email.To,
				CC:             req.Email.Cc,
				BCC:            req.Email.Bcc,
				Subject:        content.Subject,
				HTML:           content.Html,
				Text:           content.Text,
				IdempotencyKey: req.IdempotencyKey,
			}); err != nil {
				return err
			}
		} else if content := req.Email.GetTemplate(); content != nil {
			var templateID string
			if err := resourcename.Sscan(content.Name, emailTemplateResourcePattern, &templateID); err != nil {
				return err
			}

			if _, err := command.Send(ctx, rc.commandBus, command.SendEmailFromTemplate{
				From:           req.Email.From,
				To:             req.Email.To,
				CC:             req.Email.Cc,
				BCC:            req.Email.Bcc,
				TemplateID:     templateID,
				Variables:      content.Variables,
				IdempotencyKey: req.IdempotencyKey,
			}); err != nil {
				return err
			}
		} else {
			return ErrInvalidMessage
		}

	default:
		return ErrUnsupportedMessage
	}

	return nil
}

func (rc *RabbitMQConsumer) registerConsumer(ctx context.Context, queue string) error {
	consumer, err := rc.conn.NewConsumer(ctx, queue, nil)
	if err != nil {
		return err
	}

	if rc.consumers == nil {
		rc.consumers = make(map[string]*rabbitmqamqp.Consumer)
	}

	rc.consumers[queue] = consumer

	return nil
}

func (rc *RabbitMQConsumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := rabbitmqamqp.Dial(
		ctx,
		rc.c.Rabbitmq.Address,
		&rabbitmqamqp.AmqpConnOptions{
			SASLType: amqp.SASLTypePlain(
				rc.c.Rabbitmq.Username,
				rc.c.Rabbitmq.Password,
			),
		},
	)
	if err != nil {
		return err
	}

	rc.conn = conn

	if err := rc.registerAllConsumers(ctx); err != nil {
		return err
	}

	errCh := make(chan error, len(rc.consumers))

	for queue, consumer := range rc.consumers {
		rc.wg.Go(func() {
			for {
				delivery, err := consumer.Receive(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}

					continue
				}

				message := delivery.Message()

				if err := rc.handleMessage(ctx, queue, message.GetData()); err != nil {
					switch {
					case errors.Is(err, ErrUnsupportedMessage),
						errors.Is(err, ErrInvalidMessage):
						delivery.Discard(ctx, nil)

					default:
						if err := delivery.RequeueWithAnnotationsAndDeliveryFailed(ctx, nil, true); err != nil {
							errCh <- err
							return
						}
					}

					continue
				}

				delivery.Accept(ctx)
			}
		})
	}

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:
		return err
	}
}

func (rc *RabbitMQConsumer) Stop(ctx context.Context) error {
	var errs []error

	for name, consumer := range rc.consumers {
		if consumer == nil {
			delete(rc.consumers, name)
			continue
		}

		if err := consumer.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		delete(rc.consumers, name)
	}

	if rc.conn != nil {
		if err := rc.conn.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		rc.conn = nil
	}

	rc.wg.Wait()

	return errors.Join(errs...)
}
