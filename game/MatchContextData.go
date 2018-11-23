package game

import (
	"server/recastnavgo"
)

type MatchContextData struct {
	autoPathService            *recastnavgo.AutoPath
	AutoPathUtil               *AutoPathUtil
}

func (this *MatchContextData) Create() *MatchContextData {
	
	if true {
		// TODO set data for this variable
		// curSceneData := SceneData()
		// this.CreateAndLoadAutoPath(curSceneData)
	}

	this.AutoPathUtil = new(AutoPathUtil).Create(this.autoPathService) //after CreateAndLoadAutoPath

	return this
}

func (this *MatchContextData) CreateAndLoadAutoPath(sceneDate *SceneData) {

	switch sceneDate.navMeshType {
	case ENavmesh_Const:
		this.autoPathService = sceneDate.autoPathService
	case ENavmesh_Dynaimc:
		var result int
		this.autoPathService, result = new(recastnavgo.AutoPath).Create(SceneConfigGet().NavMeshCache, sceneDate.Name)
		if result != 0 {
			// util.LoggerPtr().Error("[LUOHAO] Dynaimc autoPathService create fail result = ", sceneDate.Name, result)
		} else {
			// util.LoggerPtr().Debug("[LUOHAO] Dynaimc autoPathService create Success", sceneDate.Name)
		}
	}
}

func (this *MatchContextData) Destroy() {
	// curSceneData := SceneData()
	// this.PetManager.Destroy()

	// if this.autoPathService != nil && curSceneData.navMeshType == ENavmesh_Const {
	// 	this.autoPathService.ThreadEnd()
	// }

	// if this.autoPathService != nil && curSceneData.navMeshType == ENavmesh_Dynaimc {
	// 	this.autoPathService.Destroy()
	// 	this.autoPathService = nil
	// }
}