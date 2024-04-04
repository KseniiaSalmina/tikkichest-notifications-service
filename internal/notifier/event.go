package notifier

import "errors"

type Event struct {
	Object   Object `json:"object"`
	ObjectID int    `json:"object_id"`
	Change   Change `json:"change"`
}

// Object can be Portfolio, Craft or Content
type Object string

const (
	Portfolio Object = "portfolio"
	Craft     Object = "craft"
	Content   Object = "content"
)

// Change can be CreateObj, UpdateObj or DeleteObj
type Change string

const (
	CreateObj Change = "created"
	UpdateObj Change = "changed"
	DeleteObj Change = "deleted"
)

var incorrectObjectErr = errors.New("unknown object: should be portfolio, craft or content")
var incorrectChangeErr = errors.New("unknown change: should be create, update or delete")

func validateEvent(event Event) error {
	if event.Object != Portfolio && event.Object != Craft && event.Object != Content {
		return incorrectObjectErr
	}
	if event.Change != CreateObj && event.Change != UpdateObj && event.Change != DeleteObj {
		return incorrectChangeErr
	}
	return nil
}
