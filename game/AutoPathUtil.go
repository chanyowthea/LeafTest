package game

import (
	"golang.garena.com/dts/gameserver/gcore/math"
	"golang.garena.com/dts/gameserver/gcore/util"
	// "golang.garena.com/dts/gameserver/physxgo"
	"server/recastnavgo"
)

//copied from AIEntityAgent_AutoPath to be used in non-player struct
type AutoPathComponent struct {
	autoPathNodes        []*math.Vector3
	curAutoPathNodeIndex int
	isAutoPathing        bool
	startPathTimes       int32
	targetDirection      *math.Vector3
	entity               IGetPosition
	util                 *AutoPathUtil
	//for MoveByRaycast()
	matchContextData  *MatchContextData
	queryCacheCreated bool
	queryLayers       uint
	hitLayers         uint
	height            float32
	radius            float32
}

func (this *AutoPathComponent) Create(e IGetPosition, u *AutoPathUtil) *AutoPathComponent {
	this.StopAutoPath()
	this.entity = e
	this.util = u
	return this
}
func (this *AutoPathComponent) CreateQueryCache(matchContextData *MatchContextData) {
	this.matchContextData = matchContextData
	this.queryCacheCreated = true
	dynamicLayers := uint((1 << 2) | (1 << 0))
	staticLayers := uint((1 << 1)) 
	this.hitLayers = dynamicLayers | staticLayers               //static for airdrop boxes
	this.queryLayers = staticLayers
}
func (this *AutoPathComponent) Destroy() {
	if this.queryCacheCreated {
		this.queryCacheCreated = false
	}
}

func (this *AutoPathComponent) IsAutoPathing() bool {
	return this.isAutoPathing
}
func (this *AutoPathComponent) StartPath(pathNodeArray []*math.Vector3) bool {
	this.ClearPath()
	if len(pathNodeArray) == 0 {
		return false
	}
	this.autoPathNodes = append(pathNodeArray)
	this.isAutoPathing = true
	this.startPathTimes++
	return true
}
func (this *AutoPathComponent) StopAutoPath() {
	this.ClearPath()
	this.targetDirection = math.Vector3_0()
}
func (this *AutoPathComponent) ClearPath() {
	this.autoPathNodes = make([]*math.Vector3, 0)
	this.curAutoPathNodeIndex = 0
	this.isAutoPathing = false
}
func (this *AutoPathComponent) GetCurrPathLastNode() *math.Vector3 {
	if !this.isAutoPathing || len(this.autoPathNodes) == 0 {
		return nil
	}
	return this.autoPathNodes[len(this.autoPathNodes)-1]
}

func (this *AutoPathComponent) StartAutoPath(vTarget *math.Vector3) bool {
	bFindPath, pathNodeArray := this.util.FindAutoPath(this.entity.GetPosition(), vTarget)
	if !bFindPath || len(pathNodeArray) == 0 {
		util.LoggerPtr().Debug("@@@ StartAutoPath fail", this.entity.GetPosition(), vTarget)
		return false
	}
	// if PET_DEBUG_LOG {
	// 	util.LoggerPtr().Debug("@@@ StartAutoPath---------------->", this.entity.GetPosition(), vTarget)
	// 	for i, p := range pathNodeArray {
	// 		util.LoggerPtr().Debug(i, p)
	// 	}
	// }
	return this.StartPath(pathNodeArray)
}
func (this *AutoPathComponent) CheckAndStartAutoPath(vTarget *math.Vector3) bool {
	bFindPath, _, pathNodeArray := this.util.CheckTargetReachable(this.entity.GetPosition(), vTarget)
	if !bFindPath || len(pathNodeArray) == 0 {
		//util.LoggerPtr().Debug("@zjs, CheckAndStartAutoPath fail", vTarget)
		return false
	}
	//util.LoggerPtr().Info("@zjs, CheckAndStartAutoPath---------------->", vTarget, "lastNode", pathNodeArray[0])
	return this.StartPath(pathNodeArray)
}
func (this *AutoPathComponent) UpdateAutoPath(gameTime util.TimeAbsMS, deltaTime util.TimeRelMS) {
	if this.isAutoPathing == false || len(this.autoPathNodes) == 0 {
		return
	}

	bReachDest := false
	vDest2Pos := this.entity.GetPosition().DecreaseBy(this.autoPathNodes[this.curAutoPathNodeIndex])
	vDest2Pos.Y = 0
	fDest2PosDist := vDest2Pos.Magnitude()

	if fDest2PosDist < 0.01 {
		bReachDest = true
	} else {
		fDot := this.targetDirection.Dot(vDest2Pos)

		if fDot > 0 {
			bReachDest = true
		}
	}

	if bReachDest {
		this.curAutoPathNodeIndex++
	}

	if this.curAutoPathNodeIndex >= len(this.autoPathNodes) {
		this.StopAutoPath()
	} else {
		nextPathNodeClone := this.autoPathNodes[this.curAutoPathNodeIndex].Clone()
		var startPos *math.Vector3
		if this.curAutoPathNodeIndex > 0 {
			startPos = this.autoPathNodes[this.curAutoPathNodeIndex-1]
		} else {
			startPos = this.entity.GetPosition()
		}
		vNextDest2Pos := nextPathNodeClone.DecreaseBy(startPos)
		//vNextDest2Pos.Y = math.MinF32(vNextDest2Pos.Y, 0)
		this.targetDirection = vNextDest2Pos.Normalize()
	}
}

//throws exceptions if service is nil
type AutoPathUtil struct {
	service *recastnavgo.AutoPath
}

func (this *AutoPathUtil) Create(service *recastnavgo.AutoPath) *AutoPathUtil {
	this.service = service
	return this
}

func (this *AutoPathUtil) FindAutoPath(vStart *math.Vector3, vTarget *math.Vector3) (bool, []*math.Vector3) {
	bFindPath, pathNodeArray := this.service.FindStraightPath(vStart, vTarget)
	fTargetDis := math.GetDistance(vStart, vTarget)

	for bFindPath && len(pathNodeArray) > 0 {
		fFirstNodeDist := math.GetDistance(vStart, pathNodeArray[0])
		if fFirstNodeDist < 0.01 && fTargetDis > 0.01 {
			pathNodeArray = append(pathNodeArray[:0], pathNodeArray[1:]...)
			bFindPath = len(pathNodeArray) != 0
		} else {
			break
		}
	}
	return bFindPath, pathNodeArray
}
func (this *AutoPathUtil) FindRandPointOnNearestPoly(vFindCenter *math.Vector3, fFindRadius float32, iWantPointNum int) (bFind bool, pFindPoints []*math.Vector3) {
	return this.service.FindRandPointOnNearestPoly(vFindCenter, fFindRadius, iWantPointNum)
}
func (this *AutoPathUtil) CheckTargetReachable(vStart *math.Vector3, vTarget *math.Vector3) (bool, float32, []*math.Vector3) {
	pathFound, pathNodeArray := this.FindAutoPath(vStart, vTarget)
	var pathLength float32 = 0
	if pathFound {
		if len(pathNodeArray) > 0 {
			lastNode := pathNodeArray[len(pathNodeArray)-1]
			distToTarget := math.GetDistance(vTarget, lastNode)
			if distToTarget > 0.3 {
				util.LoggerPtr().Debug("@@@ CheckTargetReachable fail, last node is not destination", distToTarget, "target", vTarget, "lastNode", lastNode)
				pathFound = false
			}
		}
		if pathFound && len(pathNodeArray) > 1 {
			for i := 1; i < len(pathNodeArray); i++ {
				pathLength += math.GetDistance(pathNodeArray[i-1], pathNodeArray[i])
			}
		}
	}
	util.LoggerPtr().Info("@@@, pathFound", pathFound)
	return pathFound, pathLength, pathNodeArray
}
