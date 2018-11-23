package recastnavgo

import (
	"golang.garena.com/dts/gameserver/gcore/math"
)

type AutoPath struct {
	navService *NavService
}

func (this *AutoPath) Create(cacheManager NavgationDataCache, sceneName string) (*AutoPath, int) {
	bResult := 0

	if len(sceneName) == 0 {
		return this, -1
	}
	if this.navService != nil {
		return this, -1
	}

	this.navService, bResult = new(NavService).Create(cacheManager, sceneName)
	return this, bResult * 10
}

func (this *AutoPath) Destroy() {
	if this.navService != nil {
		this.navService.Destroy()
		this.navService = nil
	}
}

func (this *AutoPath) ThreadEnd() {
	if this.navService != nil {
		this.navService.ThreadEnd()
	}
}

func (this *AutoPath) FindStraightPath(vStart *math.Vector3, vEnd *math.Vector3) (bFind bool, pPath []*math.Vector3) {
	if this.navService != nil {
		return this.navService.FindStraightPath(vStart, vEnd)
	}
	return false, nil
}

func (this *AutoPath) FindRandPointOnNearestPoly(vFindCenter *math.Vector3, fFindRadius float32, iWantPointNum int) (bFind bool, pFindPoints []*math.Vector3) {
	if this.navService != nil {
		return this.navService.FindRandPointOnNearestPoly(vFindCenter, fFindRadius, iWantPointNum)
	}
	return false, nil
}
func (this *AutoPath) FindClosestPointOnNearestPoly(vFindCenter *math.Vector3, fFindRadius float32) (bFind bool, pFindPoint *math.Vector3) {
	if this.navService != nil {
		return this.navService.FindClosestPointOnNearestPoly(vFindCenter, fFindRadius)
	}
	return false, nil
}

func (this *AutoPath) AddGameObjectObstaclToNavMesh(gameObjectID int, pos *math.Vector3, radius float32, height float32) bool {
	if this.navService != nil {
		return this.navService.AddObstacle(gameObjectID, pos, radius, height)
	}
	return false
}

func (this *AutoPath) RemoveGameObjectObstacl(gameObjectID int) bool {
	if this.navService != nil {
		return this.navService.RemoveObstacle(gameObjectID)
	}
	return false
}

func (this *AutoPath) UpdateService() bool {
	if this.navService != nil {
		return this.navService.Update()
	}
	return false
}
