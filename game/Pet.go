package game

import (
	"golang.garena.com/dts/gameserver/gcore/math"
)
type IGetPosition interface {
	GetPosition() *math.Vector3
}
//pet transform
type IPositionForward interface {
	IGetPosition
	SetPosition(pos *math.Vector3)
	GetForward() *math.Vector3
	SetForward(fw *math.Vector3)
}
type PetEntity struct {
	//IPositionForward
	position *math.Vector3
	forward  *math.Vector3
}
func (this *PetEntity) Create() {
	this.position = math.Vector3_0()
	this.forward = math.Vector3_Z()
}
func (this *PetEntity) GetPosition() *math.Vector3 {
	return this.position.Clone()
}
func (this *PetEntity) SetPosition(v *math.Vector3) {
	this.position = v.Clone()
}
func (this *PetEntity) GetForward() *math.Vector3 {
	return this.forward.Clone()
}
func (this *PetEntity) SetForward(fw *math.Vector3) {
	this.forward.CopyFrom(fw)
}

//pet logic
type IPet interface {
	IPositionForward
	GetHeight() float32
	GetRadius() float32
	GetPathComp() *AutoPathComponent
}
type Pet struct {
	PetEntity
	pathComp *AutoPathComponent
}

//init
func (this *Pet) Create(matchContextData *MatchContextData) *Pet {
	this.pathComp = new(AutoPathComponent).Create(this, matchContextData.AutoPathUtil)
	this.pathComp.height = this.GetHeight()
	this.pathComp.radius = this.GetRadius()
	if true {
		this.pathComp.CreateQueryCache(matchContextData)
	}

	return this
}

func (this *Pet) GetPathComp() *AutoPathComponent {
	return this.pathComp
}
func (this *Pet) GetHeight() float32 {
	return 0.4
}
func (this *Pet) GetRadius() float32 {
	return 0.3
}