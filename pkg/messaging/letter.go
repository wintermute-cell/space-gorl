package messaging

import rl "github.com/gen2brain/raylib-go/raylib"

/*
* letter.go provides a very simple and generic messaging system. It allows for
* sending of arbitrary "Message" structs to multiple receivers. For each
* message type, any number of receivers can be created. When a message of a
* type is sent, it will be appended to every receivers inbox that matches the
* message type.
 */

/* Usage:
*
* 1. Create a new message struct
* 2. Create a new receiver using NewReceiver[MyMessageType](myDomain)
* 3. Use SendMessage(myDomain, myMessage) to send a message
* 4. Fetch messages with receiver.GetNextMessage() or receiver.GetMessageAt()
*
 */


type GainScoreMessage struct {
    Amount int32
    FromKill bool
    KilledPosition rl.Vector2
}

type DamageFromEnemyMessage struct {
    Amount int32
}





type Receiver[T any] struct {
    Inbox []T
}

// NewReceiver creates and returns a new Receiver instance for a specific
// message type and a domain. Only messages of the given type and domain will
// be received.
func NewReceiver[T any](domain string) *Receiver[T] {
    rec := &Receiver[T]{Inbox: make([]T, 0)}
    ms.Receivers[domain] = append(ms.Receivers[domain], rec)
    return rec
}

// GetMessageAt retrieves a message from the inbox at the specified index.
// Returns the message and a boolean indicating whether the retrieval was
// successful.
func (r *Receiver[T]) GetMessageAt(index int) (T, bool) {
    if index >= 0 && index < len(r.Inbox) {
        return r.Inbox[index], true
    }
    var zero T
    return zero, false
}

// GetNextMessage retrieves and removes the last message from the inbox.
// Returns the message and a boolean indicating whether a message was
// available.
func (r *Receiver[T]) GetNextMessage() (T, bool) {
    l := len(r.Inbox)
    if l > 0 {
        msg := r.Inbox[l-1]
        if l == 1 {
            r.Inbox = []T{}
        } else {
            r.Inbox = r.Inbox[:l-1]
        }
        return msg, true
    }
    var zero T
    return zero, false
}


type messageSystem struct {
    Receivers map[string][]any // Map of message type to receivers
}

var ms *messageSystem

// InitMessageSystem initializes the global MessageSystem.
func InitMessageSystem() {
    ms =  &messageSystem{Receivers: make(map[string][]any)}
}

// SendMessage delivers a message to all receivers of the message type
// registered under the specified domain.
func SendMessage[T any](domain string, message T) {
    if receivers, ok := ms.Receivers[domain]; ok {
        for _, rec := range receivers {
            if receiver, ok := rec.(*Receiver[T]); ok {
                receiver.Inbox = append(receiver.Inbox, message)
            }
        }
    }
}

