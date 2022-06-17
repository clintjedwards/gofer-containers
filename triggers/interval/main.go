// Trigger interval simply triggers the subscribed pipeline repeatedly at the given interval.
package main

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/clintjedwards/gofer/gofer_sdk/go"
	"github.com/clintjedwards/gofer/gofer_sdk/go/proto"
	"github.com/clintjedwards/polyfmt"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

const (
	// "every" is the time between pipeline runs.
	// Supports golang native duration strings: https://pkg.go.dev/time#ParseDuration
	//
	// Examples: "1m", "60s", "3h", "3m30s"
	ParameterEvery = "every"
)

const (
	// The minimum duration pipelines can set for the "every" parameter.
	ConfigMinDuration = "min_duration"
)

type subscription struct {
	pipelineTriggerLabel string
	pipeline             string
	namespace            string
	quit                 context.CancelFunc
}

type subscriptionID struct {
	pipelineTriggerLabel string
	pipeline             string
	namespace            string
}

type trigger struct {
	minDuration time.Duration
	// in-memory store to be passed to the main program through the watch function
	quitAllSubscriptions context.CancelFunc
	events               chan *proto.WatchResponse
	parentContext        context.Context
	// mapping of subscription id to quit channel so we can reap the goroutines.
	subscriptions map[subscriptionID]*subscription
}

func newTrigger() (*trigger, error) {
	minDurationStr := sdk.GetConfig(ConfigMinDuration)
	minDuration := time.Minute * 1
	if minDurationStr != "" {
		parsedDuration, err := time.ParseDuration(minDurationStr)
		if err != nil {
			return nil, err
		}
		minDuration = parsedDuration
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &trigger{
		minDuration:          minDuration,
		events:               make(chan *proto.WatchResponse, 100),
		quitAllSubscriptions: cancel,
		parentContext:        ctx,
		subscriptions:        map[subscriptionID]*subscription{},
	}, nil
}

func (t *trigger) startInterval(ctx context.Context, pipeline, namespace, pipelineTriggerLabel string, duration time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			t.events <- &proto.WatchResponse{
				Details:              "Triggered due to the passage of time.",
				PipelineTriggerLabel: pipelineTriggerLabel,
				NamespaceId:          namespace,
				PipelineId:           pipeline,
				Result:               proto.WatchResponse_SUCCESS,
				Metadata:             map[string]string{},
			}
			log.Debug().Str("namespaceID", namespace).Str("pipelineID", pipeline).
				Str("trigger_label", pipelineTriggerLabel).Msg("new tick for specified interval; new event spawned")
		}
	}
}

func (t *trigger) Subscribe(ctx context.Context, request *proto.SubscribeRequest) (*proto.SubscribeResponse, error) {
	interval, exists := request.Config[ParameterEvery]
	if !exists {
		return nil, fmt.Errorf("could not find required configuration parameter %q", ParameterEvery)
	}

	duration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("could not parse interval string: %w", err)
	}

	if duration < t.minDuration {
		return nil, fmt.Errorf("durations cannot be less than %s", t.minDuration)
	}

	subctx, quit := context.WithCancel(t.parentContext)
	t.subscriptions[subscriptionID{
		request.PipelineTriggerLabel,
		request.PipelineId,
		request.NamespaceId,
	}] = &subscription{request.PipelineTriggerLabel, request.NamespaceId, request.PipelineId, quit}

	go t.startInterval(subctx, request.PipelineId, request.NamespaceId, request.PipelineTriggerLabel, duration)

	log.Debug().Str("namespace_id", request.NamespaceId).Str("trigger_label", request.PipelineTriggerLabel).Str("pipeline_id", request.PipelineId).Msg("subscribed pipeline")
	return &proto.SubscribeResponse{}, nil
}

func (t *trigger) Watch(ctx context.Context, request *proto.WatchRequest) (*proto.WatchResponse, error) {
	select {
	case <-ctx.Done():
		return &proto.WatchResponse{}, nil
	case event := <-t.events:
		return event, nil
	}
}

func (t *trigger) Unsubscribe(ctx context.Context, request *proto.UnsubscribeRequest) (*proto.UnsubscribeResponse, error) {
	subscription, exists := t.subscriptions[subscriptionID{
		pipelineTriggerLabel: request.PipelineTriggerLabel,
		pipeline:             request.PipelineId,
		namespace:            request.NamespaceId,
	}]
	if !exists {
		return &proto.UnsubscribeResponse{},
			fmt.Errorf("could not find subscription for trigger %s pipeline %s namespace %s",
				request.PipelineTriggerLabel, request.PipelineId, request.NamespaceId)
	}

	subscription.quit()
	delete(t.subscriptions, subscriptionID{
		pipelineTriggerLabel: request.PipelineTriggerLabel,
		pipeline:             request.PipelineId,
		namespace:            request.NamespaceId,
	})
	return &proto.UnsubscribeResponse{}, nil
}

func (t *trigger) Info(ctx context.Context, request *proto.InfoRequest) (*proto.InfoResponse, error) {
	return sdk.InfoResponse("https://clintjedwards.com/gofer/docs/triggers/interval/overview")
}

func (t *trigger) ExternalEvent(ctx context.Context, request *proto.ExternalEventRequest) (*proto.ExternalEventResponse, error) {
	return &proto.ExternalEventResponse{}, nil
}

func (t *trigger) Shutdown(ctx context.Context, request *proto.ShutdownRequest) (*proto.ShutdownResponse, error) {
	t.quitAllSubscriptions()
	close(t.events)

	return &proto.ShutdownResponse{}, nil
}

func mustInitFormatter() polyfmt.Formatter {
	fmtter, err := polyfmt.NewFormatter(polyfmt.Pretty, false)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start formatter")
	}
	return fmtter
}

func installer() {
	fmtter := mustInitFormatter()
	headerColor := color.New(color.Underline, color.Bold, color.FgBlue)

	fmtter.Println(headerColor.Sprintf("Interval Trigger Setup\n"))
	fmtter.Println(":: The interval trigger allows users to trigger their pipelines on the " +
		"passage of time by setting a particular duration.\n")
	fmtter.Println("First, let's prevent users from setting too low of an interval by setting a minimum duration. " +
		"Durations are set via Golang duration strings. For example, entering a duration of '10h' would be 10 hours. " +
		"You can find more documentation on valid strings here: https://pkg.go.dev/time#ParseDuration.\n")

	fmtter.Finish()
	var minDuration string
	fmt.Print("> Set a minimum duration for all pipelines: ")
	fmt.Scanln(&minDuration)
	fmtter = mustInitFormatter()

	_, err := time.ParseDuration(minDuration)
	if err != nil {
		fmtter.PrintSuccess(fmt.Sprintf("minimum duration string %q is invalid: %v", minDuration, err))
		fmtter.Finish()
		return
	}
	fmtter.PrintSuccess(fmt.Sprintf("Valid minimum duration %q set.", minDuration))
	fmtter.PrintSuccess("Trigger configuration complete.")
	fmtter.Println("")

	fmtter.Print("Registering Trigger")

	config := map[string]string{
		"MIN_DURATION": minDuration,
	}

	err = sdk.InstallTrigger(config)
	if err != nil {
		fmtter.PrintErr(fmt.Sprintf("could not register trigger: %v", err))
		fmtter.Finish()
		return
	}

	fmtter.PrintSuccess("Registered Trigger")
	fmtter.Finish()
}

func main() {
	trigger, err := newTrigger()
	if err != nil {
		panic(err)
	}
	sdk.NewTrigger(trigger, installer)
}
