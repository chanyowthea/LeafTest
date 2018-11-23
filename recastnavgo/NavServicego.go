package recastnavgo

import (
	"golang.garena.com/dts/gameserver/gcore/math"
	"golang.garena.com/dts/gameserver/gcore/util"
)

const MAX_PLOYS = 256

type NavService struct {
	cNavService NavgationService
}

func (this *NavService) Create(cacheManager NavgationDataCache, sceneName string) (*NavService, int) {
	this.cNavService = NewNavgationService()
	return this, this.cNavService.LoadNavMeshFromDataCache(cacheManager, sceneName)
}

func (this *NavService) Destroy() {
	DeleteNavgationService(this.cNavService)
}

func (this *NavService) ThreadEnd() {
	this.cNavService.ThreadEnd()
}
func (this *NavService) FindStraightPath(vStart *math.Vector3, vEnd *math.Vector3) (bFind bool, pPath []*math.Vector3) {

	var fsStart [3]float32
	var fsEnd [3]float32

	fsStart[0] = vStart.X
	fsStart[1] = vStart.Y
	fsStart[2] = vStart.Z

	fsEnd[0] = vEnd.X
	fsEnd[1] = vEnd.Y
	fsEnd[2] = vEnd.Z

	pathArray := make([]float32, MAX_PLOYS*3, MAX_PLOYS*3)
	nodeCount := 0

	bFindPath := this.cNavService.FindStraightPath(&fsStart[0], &fsEnd[0], &pathArray[0], &nodeCount)

	var path []*math.Vector3

	for i := 0; i < nodeCount; i++ {
		node := math.Vector3_0()
		node.X = pathArray[i*3]
		node.Y = pathArray[i*3+1]
		node.Z = pathArray[i*3+2]
		path = append(path, node)
	}

	return bFindPath, path
}

func (this *NavService) FindRandPointOnNearestPoly(vFindCenter *math.Vector3, fFindRadius float32, iWantPointNum int) (bFind bool, pFindPoints []*math.Vector3) {

	var fCenter [3]float32

	fCenter[0] = vFindCenter.X
	fCenter[1] = vFindCenter.Y
	fCenter[2] = vFindCenter.Z

	wantPointArray := make([]float32, iWantPointNum*3, iWantPointNum*3)

	iFindPointNum := 0
	bFindPoint := this.cNavService.FindRandPointOnNearestPoly(&fCenter[0], fFindRadius, iWantPointNum, &wantPointArray[0], &iFindPointNum)

	var pFindResult []*math.Vector3
	for i := 0; i < iFindPointNum; i++ {
		node := math.Vector3_0()
		node.X = wantPointArray[i*3]
		node.Y = wantPointArray[i*3+1]
		node.Z = wantPointArray[i*3+2]
		pFindResult = append(pFindResult, node)
	}

	return bFindPoint, pFindResult
}
func (this *NavService) FindClosestPointOnNearestPoly(vFindCenter *math.Vector3, fFindRadius float32) (bFind bool, pFindPoint *math.Vector3) {

	var fCenter [3]float32

	fCenter[0] = vFindCenter.X
	fCenter[1] = vFindCenter.Y
	fCenter[2] = vFindCenter.Z

	var fResult [3]float32
	bFindPoint := this.cNavService.FindClosestPointOnNearestPoly(&fCenter[0], fFindRadius, &fResult[0])

	pFindPoint = &math.Vector3{
		fResult[0], fResult[1], fResult[2],
	}

	return bFindPoint, pFindPoint
}

func (this *NavService) AddObstacle(gameObjectID int, pos *math.Vector3, radius float32, height float32) bool {
	var fPos [3]float32
	fPos[0] = pos.X
	fPos[1] = pos.Y
	fPos[2] = pos.Z

	return this.cNavService.AddObstacle(gameObjectID, &fPos[0], radius, height)
}

func (this *NavService) RemoveObstacle(gameObjectID int) bool {
	return this.cNavService.RemoveObstacle(gameObjectID)
}

func (this *NavService) Update() bool {
	bUptodate := true

	result := this.cNavService.UpdateNavmesh(&bUptodate)

	if bUptodate == false {
		util.LoggerPtr().Warning("[LUOHAO] ---updating-------------")
	}

	return result
}
