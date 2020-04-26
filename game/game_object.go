package game

import (
	"math"

	"github.com/donbattery/bnj/model"
)

type gameObject struct {
	id       string
	parentId string
	objType  string
	anim     int
	x        float64
	y        float64
	flipX    bool
	flipY    bool
	vector   *vector
}

func newGameObject(id, parentId, objType string, x, y float64) *gameObject {
	return &gameObject{
		id:       id,
		parentId: parentId,
		objType:  objType,
		x:        x,
		y:        y,
		vector:   newVector(0, 0),
	}
}

func (obj *gameObject) dump() model.GameObjectDump {
	return model.GameObjectDump{
		ObjType: obj.objType,
		Anim:    obj.anim,
		X:       int(math.Round(obj.x)),
		Y:       int(math.Round(obj.y)),
		FlipX:   obj.flipX,
		FlipY:   obj.flipY,
	}
}
