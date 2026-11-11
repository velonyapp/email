package transport

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	v1 "github.com/velonyapp/email/gen/job/v1"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/presentation/service"

	"github.com/Azure/go-amqp"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

const (
	sendEmailQueue = "email.send-email"
)

var _ transport.Server = (*RabbitMQConsumer)(nil)

type RabbitMQConsumer struct {
	c          *conf.Transport
	jobService *service.Job
	log        *slog.Logger

	conn        *rabbitmqamqp.AmqpConnection
	consumerMap map[string]*rabbitmqamqp.Consumer

	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	cond   *sync.Cond
	wg     sync.WaitGroup
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	jobService *service.Job,
	log *slog.Logger,
) *RabbitMQConsumer {
	srv := &RabbitMQConsumer{
		c:          c,
		jobService: jobService,
		log:        log,

		consumerMap: make(map[string]*rabbitmqamqp.Consumer),
	}
	srv.cond = sync.NewCond(&srv.mu)

	return srv
}

func (rc *RabbitMQConsumer) registerConsumerMap(ctx context.Context) error {
	var errs []error

	if err := rc.registerConsumer(ctx, sendEmailQueue); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, rc.c.Rabbitmq.Timeout.AsDuration())
	defer cancel()

	switch queue {
	case sendEmailQueue:
		payload := new(v1.SendEmailPayload)

		if err := protojson.Unmarshal(data, payload); err != nil {
			return ErrInvalidMessage
		}

		return rc.jobService.SendEmail(ctx, payload)

	default:
		return ErrUnsupportedMessage
	}
}

func (rc *RabbitMQConsumer) registerConsumer(ctx context.Context, queue string) error {
	consumer, err := rc.conn.NewConsumer(ctx, queue, nil)
	if err != nil {
		rc.log.ErrorContext(ctx,
			"failed to register RabbitMQ consumer",
			slog.String("queue", queue),
			slog.Any("error", err),
		)
		return err
	}

	rc.consumerMap[queue] = consumer

	rc.log.InfoContext(ctx,
		"RabbitMQ consumer registered",
		slog.String("queue", queue),
	)

	return nil
}

func (rc *RabbitMQConsumer) Start(ctx context.Context) error {
	rc.ctx, rc.cancel = context.WithCancel(ctx)
	defer rc.cancel()

	conn, err := rabbitmqamqp.Dial(ctx,
		rc.c.Rabbitmq.Address,
		&rabbitmqamqp.AmqpConnOptions{
			SASLType: amqp.SASLTypePlain(
				rc.c.Rabbitmq.Username,
				rc.c.Rabbitmq.Password,
			),
			RecoveryConfiguration: &rabbitmqamqp.RecoveryConfiguration{
				ActiveRecovery:           true,
				BackOffReconnectInterval: 5 * time.Second,
				MaxReconnectAttempts:     5,
			},
		},
	)
	if err != nil {
		return err
	}

	_, connected := conn.State().(*rabbitmqamqp.StateOpen)

	stateChanged := make(chan *rabbitmqamqp.StateChanged, 1)
	go func(ch chan *rabbitmqamqp.StateChanged) {
		for changed := range ch {
			switch changed.To.(type) {
			case *rabbitmqamqp.StateOpen:
				rc.mu.Lock()
				connected = true
				rc.cond.Broadcast()
				rc.mu.Unlock()

			case *rabbitmqamqp.StateClosed:
				rc.mu.Lock()
				connected = false
				rc.mu.Unlock()
			}
		}
	}(stateChanged)

	conn.NotifyStatusChange(stateChanged)

	rc.conn = conn

	if err := rc.registerConsumerMap(ctx); err != nil {
		return err
	}

	for queue, consumer := range rc.consumerMap {
		rc.wg.Go(func() {
			for {
				rc.mu.Lock()
				for !connected && rc.ctx.Err() == nil {
					rc.cond.Wait()
				}
				if rc.ctx.Err() != nil {
					rc.mu.Unlock()
					return
				}
				rc.mu.Unlock()

				delivery, err := consumer.Receive(rc.ctx)
				if err != nil {
					if rc.ctx.Err() != nil {
						return
					}

					rc.log.WarnContext(ctx,
						"failed to receive RabbitMQ message",
						slog.String("queue", queue),
						slog.Any("error", err),
					)
					continue
				}

				message := delivery.Message()

				if err := rc.handleMessage(ctx, queue, message.GetData()); err != nil {
					switch {
					case errors.Is(err, ErrUnsupportedMessage),
						errors.Is(err, ErrInvalidMessage):
						if err := delivery.Discard(ctx, nil); err != nil {
							rc.log.ErrorContext(ctx,
								"failed to discard RabbitMQ message",
								slog.String("queue", queue),
								slog.Any("error", err),
							)
						}

					default:
						rc.log.WarnContext(ctx,
							"failed to handle RabbitMQ message, requeueing",
							slog.String("queue", queue),
							slog.Any("error", err),
						)

						if err := delivery.RequeueWithAnnotationsAndDeliveryFailed(ctx, nil, true); err != nil {
							rc.log.ErrorContext(ctx,
								"failed to requeue RabbitMQ message",
								slog.String("queue", queue),
								slog.Any("error", err),
							)
						}
					}

					continue
				}

				if err := delivery.Accept(ctx); err != nil {
					rc.log.ErrorContext(ctx,
						"failed to accept RabbitMQ message",
						slog.String("queue", queue),
						slog.Any("error", err),
					)
				}
			}
		})
	}

	rc.wg.Wait()

	return nil
}

func (rc *RabbitMQConsumer) Stop(ctx context.Context) error {
	rc.cancel()
	rc.mu.Lock()
	rc.cond.Broadcast()
	rc.mu.Unlock()

	rc.wg.Wait()

	var errs []error

	for name, consumer := range rc.consumerMap {
		if err := consumer.Close(ctx); err != nil {
			rc.log.ErrorContext(ctx,
				"failed to close RabbitMQ consumer",
				slog.String("queue", name),
				slog.Any("error", err),
			)
			errs = append(errs, err)
		}

		delete(rc.consumerMap, name)
	}

	if rc.conn != nil {
		if err := rc.conn.Close(ctx); err != nil {
			rc.log.ErrorContext(ctx,
				"failed to close RabbitMQ connection",
				slog.Any("error", err),
			)
			errs = append(errs, err)
		}

		rc.conn = nil
	}

	if len(errs) == 0 {
		rc.log.InfoContext(ctx,
			"RabbitMQ consumer stopped",
		)
	}

	return errors.Join(errs...)
}
