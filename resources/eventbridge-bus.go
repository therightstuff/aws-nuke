package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type EventBridgeBus struct {
	svc                *eventbridge.EventBridge
	eventBridgeBusName *string
}

func init() {
	register("EventBridgeBus", ListEventBridgeBuses)
}

func ListEventBridgeBuses(sess *session.Session) ([]Resource, error) {
	svc := eventbridge.New(sess)
	eventBridgeBuses := []Resource{}

	params := &eventbridge.ListEventBusesInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListEventBuses(params)
		if err != nil {
			return nil, err
		}

		for _, eventBus := range output.EventBuses {
			eventBridgeBuses = append(eventBridgeBuses, &EventBridgeBus{
				svc:                svc,
				eventBridgeBusName: eventBus.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return eventBridgeBuses, nil
}

func (b *EventBridgeBus) Remove() error {
	_, err := b.svc.DeleteEventBus(&eventbridge.DeleteEventBusInput{
		Name: b.eventBridgeBusName,
	})

	return err
}

func (b *EventBridgeBus) String() string {
	return *b.eventBridgeBusName
}
