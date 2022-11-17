package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type EventBridgeRule struct {
	svc                 *eventbridge.EventBridge
	eventBridgeBusName  *string
	eventBridgeRuleName *string
}

func init() {
	register("EventBridgeRule", ListEventBridgeRules)
}

func ListEventBridgeRules(sess *session.Session) ([]Resource, error) {
	svc := eventbridge.New(sess)
	eventBridgeRules := []Resource{}

	params := &eventbridge.ListRulesInput{
		Limit: aws.Int64(25),
	}

	eventBridgeBuses, err := ListEventBridgeBuses(sess)
	if err != nil {
		return nil, err
	}

	for _, eventBridgeBus := range eventBridgeBuses {
		params.EventBusName = eventBridgeBus.(*EventBridgeBus).eventBridgeBusName
		params.NextToken = nil

		for {
			output, err := svc.ListRules(params)
			if err != nil {
				return nil, err
			}

			for _, rule := range output.Rules {
				eventBridgeRules = append(eventBridgeRules, &EventBridgeRule{
					svc:                 svc,
					eventBridgeBusName:  rule.EventBusName,
					eventBridgeRuleName: rule.Name,
				})
			}

			if output.NextToken == nil {
				break
			}

			params.NextToken = output.NextToken
		}
	}

	return eventBridgeRules, nil
}

func (r *EventBridgeRule) ListTargetIds() ([]*string, error) {
	targetIds := []*string{}

	params := &eventbridge.ListTargetsByRuleInput{
		EventBusName: r.eventBridgeBusName,
		Rule:         r.eventBridgeRuleName,
		Limit:        aws.Int64(25),
	}

	for {
		output, err := r.svc.ListTargetsByRule(params)
		if err != nil {
			return nil, err
		}

		for _, target := range output.Targets {
			targetIds = append(targetIds, target.Id)
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return targetIds, nil
}

func (r *EventBridgeRule) Remove() error {
	targetIds, listTargetIdsErr := r.ListTargetIds()
	if listTargetIdsErr != nil {
		return listTargetIdsErr
	}

	if len(targetIds) > 0 {
		_, removeTargetsErr := r.svc.RemoveTargets(&eventbridge.RemoveTargetsInput{
			Ids:          targetIds,
			EventBusName: r.eventBridgeBusName,
			Rule:         r.eventBridgeRuleName,
		})
		if removeTargetsErr != nil {
			return removeTargetsErr
		}
	}

	_, deleteRuleErr := r.svc.DeleteRule(&eventbridge.DeleteRuleInput{
		EventBusName: r.eventBridgeBusName,
		Name:         r.eventBridgeRuleName,
	})

	return deleteRuleErr
}

func (r *EventBridgeRule) String() string {
	return *r.eventBridgeRuleName
}
