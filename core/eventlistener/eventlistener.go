package eventlistener

import (
	"errors"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

const ()

type EventListener struct {
	EventId  string
	NotifyCh chan interface{}
	TimeOut  time.Duration
}

func NewEventListener(eventID string,
	timeOut time.Duration) *EventListener {

	return &EventListener{
		EventId:  eventID,
		NotifyCh: make(chan interface{}),
		TimeOut:  timeOut,
	}
}

func (el *EventListener) DeleteEventListener() {
	if el.NotifyCh != nil {
		close(el.NotifyCh)
	}
}

func (el *EventListener) SendEventNotify(content interface{}) {
	notifyCh := el.NotifyCh
	notifyCh <- content
}

func (el *EventListener) WaitEvent() (interface{}, error) {
	if el.TimeOut > 0 {
		select {
		case <-time.After(el.TimeOut):
			return nil, errors.New("timeout!")
		case v, ok := <-el.NotifyCh:
			if !ok {
				return nil, errors.New("channel has been closed!")
			}

			return v, nil
		}
	}

	v, ok := <-el.NotifyCh
	if !ok {
		return nil, errors.New("channel has been closed!")
	}

	return v, nil
}

/*
* Event Listener manager
* how to use ?
* 1. create event listener manager.
*  mgr := NewEventListenerManager
* 2. send message with ID
* 3. wait message reply
* mgr.WatchEvent(msgId, timeOut)
* 4. call mgr.MatchEventAndDispatch when ID matched!
* 5. recieved the message or timeout at mgr.WatchEvent sides.
 */
type EventListenerManager struct {
	EventListenerMap *sync.Map
}

func NewEventListenerManager() *EventListenerManager {
	var listenerMap sync.Map

	return &EventListenerManager{
		EventListenerMap: &listenerMap,
	}
}

func (elm *EventListenerManager) GetEventListener(eventID string) *EventListener {
	v, exist := elm.EventListenerMap.Load(eventID)
	if !exist {
		return nil
	}

	eventListener, isThisType := v.(*EventListener)
	if !isThisType {
		return nil
	}

	return eventListener
}

func (elm *EventListenerManager) PutEventListener(listener *EventListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}

	eventID := listener.EventId
	if eventID == "" {
		return errors.New("listener eventID is empty")
	}

	if elm.GetEventListener(eventID) != nil {
		return errors.New("listener has exists")
	}

	elm.EventListenerMap.Store(eventID, listener)

	return nil
}

func (elm *EventListenerManager) DeleteEventListener(listener *EventListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}

	eventID := listener.EventId
	if eventID == "" {
		return errors.New("listener eventID is empty")
	}

	if elm.GetEventListener(eventID) == nil {
		return errors.New("listener not exists")
	}

	elm.EventListenerMap.Delete(eventID)
	listener.DeleteEventListener()

	return nil
}

/* register the event listener. */
func (elm *EventListenerManager) RegisterEventListener(eventID string,
	timeOut time.Duration) (*EventListener, error) {

	eventListener := NewEventListener(eventID, timeOut)

	/* Add event listener into ListenerManager */
	err := elm.PutEventListener(eventListener)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return eventListener, nil
}

/* unregister the event listener.*/
func (elm *EventListenerManager) UnregisterEventListener(listener *EventListener) error {
	err := elm.DeleteEventListener(listener)
	//*listener = nil

	return err
}

/* Match the event and dispatch it. */
func (elm *EventListenerManager) MatchEventAndDispatch(eventID string, content interface{}) error {

	listener := elm.GetEventListener(eventID)
	if listener == nil {
		//No matched event, we just return.
		return nil
	}

	/* matched, then we dispatch the event */
	listener.SendEventNotify(content)

	return nil
}

/* Watch the event.*/
func (elm *EventListenerManager) WatchEvent(eventID string, timeOut time.Duration) (interface{}, error) {

	listener, err := elm.RegisterEventListener(eventID, timeOut)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	content, err := listener.WaitEvent()
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	err = elm.UnregisterEventListener(listener)
	return content, err
}
