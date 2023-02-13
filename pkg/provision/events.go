package provision

import (
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

func sendSchedulingDoneEvent(podName, namespace, msg string) (*cloudevents.Event, cloudevents.Result) {
	event, err := newSchedulingDoneEvent(podName, namespace, Source)
	if err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to send scheduled done event: %s", err)
	}
	log.Printf(msg)

	return event, nil
}

func sendSchedulingFailedEvent(err error) (*cloudevents.Event, cloudevents.Result) {
	event, err := newSchedulingFailedEvent(err, Source)
	if err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to create Client scheduled failed event: %s", err)
	}
	log.Printf("A new Client scheduling failed: %s", err)

	return event, nil
}

func newSchedulingDoneEvent(podName, namespace, source string) (*event.Event, error) {
	data := ClientSchedulingDone{
		PodName:   podName,
		Namespace: namespace,
		Message:   ClientSchedulingDoneMessage,
	}

	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetSource(source)
	e.SetType(TypeClientSchedulingDone)
	if err := e.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return nil, err
	}

	return &e, nil
}

func newSchedulingFailedEvent(err error, source string) (*event.Event, error) {
	data := ClientSchedulingFailure{
		Message: err.Error(),
	}

	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetSource(source)
	e.SetType(TypeClientSchedulingFailure)
	if err := e.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return nil, err
	}

	return &e, nil
}
