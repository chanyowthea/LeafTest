package game

import (
	// "encoding/binary"
	"encoding/json"
	"io/ioutil"
	// "os"
	"strings"

	"golang.garena.com/dts/gameserver/gcore/math"
	"golang.garena.com/dts/gameserver/gcore/util"
	"server/recastnavgo"
)

const (
	_ byte = iota
	EPrimativeColliderType_Box
	EPrimativeColliderType_Sphere
	EPrimativeColliderType_Capsule
)

const (
	EPrimativeColliderLayer_CanFireThrough byte = 14
)

type PrimativeColliderData struct {
	ObjectName      string
	ObjectNameBytes []byte
	ColliderType    byte
	ColliderLayer   byte
	Location        math.Vector3
	Rotation        math.Quaternion
	Scale           math.Vector3
	Extent          math.Vector3
}

func (this *PrimativeColliderData) Create() *PrimativeColliderData {
	this.ObjectNameBytes = make([]byte, 128)
	return this
}

type MeshColliderData struct {
	ObjectName      string
	ObjectNameBytes []byte
	ColliderLayer   byte
	IsConvex        byte
	Location        math.Vector3
	Rotation        math.Quaternion
	Scale           math.Vector3
	Vertices        []float32
	Indices         []int32
}

func (this *MeshColliderData) Create() *MeshColliderData {
	this.ObjectNameBytes = make([]byte, 128)
	return this
}

type PlayerSpawnPoint struct {
	ID       int
	Postion  []string
	Forward []string
}

type ZombieSpawnPoint struct {
	ID       int
	Position  []string
}

type AISpawnPointData struct {
	Id                  int
	SpawnAIType         int
	SpecificPathGroupId int
	AttackSafeDoorId    int
	Position            []string
}

type LevelRoomsExtraData struct {
	Id    int
	Neighbors []struct {
		Id uint32
	}
	AISpawnPoints []AISpawnPointData
}

type LevelDoorsExtraData struct {
	Id       int
	CostCoin int
}

type PVEAISpecificPathData struct {
	PathGroupId int
	NaviSpots   []struct {
		NaviSpotPosition []string
	}
}

type PlayerSpawnRegionDataJson struct {
	Position []string
	Radius   float32
}

type ColliderData struct {
	Position     []string `json:"position"`
	Rotation     []string `json:"rotation"`
	Size         []string `json:"size"`
	ColliderType int      `json:"colliderType"`
	IsTrigger    bool     `json:"isTrigger"`
}

type LevelObjectData struct {
	Scene   string `json:"scene"`
	Objects []struct {
		Name      string   `json:"name"`
		Position  []string `json:"position"`
		Rotation  []string `json:"rotation"`
		Lossycale []string `json:"lossycale"`
		Colliders []struct {
			Position     []string `json:"position"`
			Rotation     []string `json:"rotation"`
			Lossycale    []string `json:"lossycale"`
			Size         []string `json:"size"`
			ColliderType int      `json:"colliderType"`
			IsTrigger    bool     `json:"isTrigger"`
		} `json:"colliders"`
		UserData string `json:"userData"`
	} `json:"objects"`

	PlayerSpawnPoints []PlayerSpawnPoint
	ZombieSpawnPoints []ZombieSpawnPoint

	LevelRoomsExtraDatas []LevelRoomsExtraData
	LevelDoorsExtraDatas []LevelDoorsExtraData
	PVEAISpecificPaths   []PVEAISpecificPathData
	PlayerSpawnRegion    PlayerSpawnRegionDataJson
}

const (
	ELevelObjectName_Door                   = "tdoor"
	ELevelObjectName_Container              = "container"
	ELevelObjectName_Vehicle                = "vehicle"
	ELevelObjectName_SpawnPoint             = "spawnpoint"
	ELevelObjectName_Airdrop                = "airdrop"
	ELevelObjectName_Shallow                = "shoal"
	ELevelObjectName_Projectile             = "projectile"
	ELevelObjectName_Fountain               = "fountainsafe"
	ELevelObjectName_UAV                    = "uav"
	ELevelObjectName_PveSafeDoor            = "pveSafeDoor"
	ELevelObjectName_Room                   = "troom"
	ELevelObjectName_Water                  = "water"
	ELevelObjectName_WaterLand              = "waterland"
	ELevelObjectName_KillZone               = "killzone"
	ELevelObjectName_NewContainer_Prefix    = "newcontainer"
	ELevelObjectName_NewContainerMushroom   = "newcontainerPickupMushroom"
	ELevelObjectName_NewContainerInstant    = "newcontainerPickupInstant"
	ELevelObjectName_NewContainerArmortools = "newcontainerPickupArmortools"
	ELevelObjectName_ClimbingTrigger        = "climbingTrigger"
	ELevelObjectName_Strop                  = "strop"
	ELevelObjectName_ForbiddenArea          = "forbiddenarea"
	ELevelObjectName_Oildrum                = "Oildrum"
	ELevelObjectName_CampFire               = "campFire"
	ELevelObjectName_Treasure               = "treasureContainer"
	ELevelObjectName_Football               = "football"
	ELevelObjectName_FootballGoal           = "footballGoal"
	ELevelObjectName_Missile                = "missile"
	ELevelObjectName_IceWall				= "IceWall"
	ELevelObjectName_MovePlatform		  = "MovePlatform"
	ELevelObjectName_Train		  		  = "Train"
	ELevelObjectName_WaterSurf				= "surftrigger"
	ELevelObjectName_Landmine				= "Landmine"
	ELevelObjectName_RedEnvelope			= "RedEnvelope"
	ELevelObjectName_DefenderPoint			= "DefenderPoint"
	ELevelObjectName_SmokeDrum				= "SmokeDrum"
)

const (
	PHYSICS_GRAVITY                        float32 = -9.8
	SCENE_DATA_PATH                        string  = "SceneData/"
	CONFIG_DATA_PATH                       string  = "Config/"
	SCENE_DATA_CONFIG_JSON                 string  = "Scene.json"
	SCENE_DATA_CONFIG_BIN                  string  = "Scene.bin"
	SCENE_NAV_DATA                         string  = "Scene.nav"
	AIRLINE_CONFIG_FILENAME                string  = "AirLines.json"
	BOT_FIGHT_BEHAVIORTREE_FILENAME        string  = "aitree.json"
	BOT_INSKY_BEHAVIORTREE_FILENAME        string  = "aitreeInSky.json"
	BOT_WAITINGROOM_BEHAVIORTREE_FILENAME  string  = "AITreeInWaitingRoom.json"
	AI_ZOMBIEGENERAL_BEHAVIORTREE_FILENAME string  = "aitree_zombiegeneral.json"
	AI_ZOMBIETHROW_BEHAVIORTREE_FILENAME   string  = "aitree_zombieThrow.json"
	PET_BEHAVIORTREE_FILENAME              string  = "aitree_pet.json"
)

type AirtransportLine struct {
	Lines []struct {
		InnerSphere     []string `json:"InnerSphere"`
		OuterSphere     []string `json:"OuterSphere"`
		InnerDiameter   string   `json:"InnerDiameter"`
		OuterDiameter   string   `json:"OuterDiameter"`
		EndJumpDeltaMin string   `json:"EndJumpDeltaMin"`
		EndJumpDeltaMax string   `json:"EndJumpDeltaMax"`
		Duration        string   `json:"Duration"`
	} `json:"Lines"`
}

func GetLevelObjectNameFromID(prefix string, id uint32) string {
	return prefix + util.UIntToString(uint(id))
}
func GetLevelObjectIDFromName(name string, prefix string) (uint32, bool) {
	if strings.Contains(name, prefix) {
		idStr := util.Substr(name, len(prefix), len(name)-len(prefix))
		return util.StringToUInt32(idStr)
	}
	return 0, false
}
func IsLevelObjectTypeOf(name string, prefix string) bool {
	return strings.Contains(name, prefix)
}

func IsColliderLayerCanFireThrough(layer byte) bool {
	if layer == EPrimativeColliderLayer_CanFireThrough {
		return true
	}
	return false
}

//singleton
var __SceneConfigInstance *SceneConfig

func SceneConfigGet() *SceneConfig {
	if __SceneConfigInstance == nil {
		__SceneConfigInstance = new(SceneConfig).Create()
	}
	return __SceneConfigInstance
}

type MapData struct {
	UniqueID				uint32
	MapID			uint32
	ModeID			uint32
	Name			string
	IsOpen			bool
	NavmeshType		uint32``

	WaitingMapName	string
}

type SceneData struct {
	Name                   string
	LevelObjectDatas       *LevelObjectData
	AirtransportLineDatas  *AirtransportLine
	PrimativeColliderDatas []*PrimativeColliderData
	MeshColliderDatas      []*MeshColliderData

	mapData *MapData

	WaitingMapSceneData *SceneData

	navMeshType     uint32
	autoPathService *recastnavgo.AutoPath
}

const (
	ENavmesh_No      uint32 = 0
	ENavmesh_Const   uint32 = 1
	ENavmesh_Dynaimc uint32 = 2
)

func (this *SceneData) GetSceneDataPath() string {
	return SCENE_DATA_PATH + this.Name + "/"
}

func (this *SceneData) GetAutoPath() *recastnavgo.AutoPath {
	return this.autoPathService
}

func (this *SceneData) Create(name string, navMeshType uint32, waitingmapName string) *SceneData {
	// 场景名称
	this.Name = name
	this.LevelObjectDatas = new(LevelObjectData)
	this.AirtransportLineDatas = new(AirtransportLine)
	this.PrimativeColliderDatas = make([]*PrimativeColliderData, 0)
	this.MeshColliderDatas = make([]*MeshColliderData, 0)

	// 两种navMesh有什么区别？
	this.navMeshType = navMeshType
	if waitingmapName != "" {
		this.WaitingMapSceneData = new(SceneData).Create(waitingmapName, ENavmesh_Const, "")
	} else {
		this.WaitingMapSceneData = nil
	}

	return this
}

func (this *SceneData) Load(config *SceneConfig) bool {
	//---------------------------
	// jsonDataPath := this.GetSceneDataPath() + this.Name + ".json"
	// util.LoggerPtr().Info("load", jsonDataPath)
	// // 从json中加载二进制数据
	// b, err := ioutil.ReadFile(jsonDataPath)
	// if err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }
	// // 将二进制数据转化成LevelObjectDatas
	// if err := json.Unmarshal(b, this.LevelObjectDatas); err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }

	// // 加载航线数据
	// //--------------------------
	// AirLineDataPath := this.GetSceneDataPath() + AIRLINE_CONFIG_FILENAME
	// util.LoggerPtr().Info("load", AirLineDataPath)
	// lines, err := ioutil.ReadFile(AirLineDataPath)
	// if err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }
	// if err := json.Unmarshal(lines, this.AirtransportLineDatas); err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }

	// // 加载binary数据，干什么的？
	// //--------------------------
	// binDataPath := this.GetSceneDataPath() + this.Name + ".bin"
	// util.LoggerPtr().Info("load", binDataPath)
	// f, err := os.Open(binDataPath)
	// if err != nil {
	// 	util.LoggerPtr().Error("Load", binDataPath, "Failed:", err)
	// 	return false
	// }
	// defer f.Close()
	// header := make([]byte, 4)
	// if err = binary.Read(f, binary.LittleEndian, header); err != nil {
	// 	return false
	// }
	// if string(header) != "NBSL" {
	// 	util.LoggerPtr().Error("corrupted file", binDataPath)
	// 	return false
	// }
	// version := uint32(0)
	// if err = binary.Read(f, binary.LittleEndian, &version); err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }
	// util.LoggerPtr().Debug("version =", version)
	// primativeColliderCount := uint32(0)
	// if err = binary.Read(f, binary.LittleEndian, &primativeColliderCount); err != nil {
	// 	util.LoggerPtr().Error(err)
	// 	return false
	// }
	// util.LoggerPtr().Debug("primativeColliderCount =", primativeColliderCount)
	// for i := 0; i < int(primativeColliderCount); i++ {
	// 	colliderData := new(PrimativeColliderData).Create()
	// 	binary.Read(f, binary.LittleEndian, colliderData.ObjectNameBytes)
	// 	colliderData.ObjectName = util.BytesToString(colliderData.ObjectNameBytes)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.ColliderType)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.ColliderLayer)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.Z)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.Z)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.W)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.Z)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Extent.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Extent.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Extent.Z)
	// 	//add collider datas
	// 	this.PrimativeColliderDatas = append(this.PrimativeColliderDatas, colliderData)
	// 	//util.LoggerPtr().Debug("primativeCollider colliderData.ObjectName=", string(colliderData.ObjectName))
	// }
	// meshColliderCount := uint32(0)
	// if err = binary.Read(f, binary.LittleEndian, &meshColliderCount); err != nil {
	// 	util.LoggerPtr().Debug(err)
	// 	return false
	// }
	// util.LoggerPtr().Debug("meshColliderCount =", meshColliderCount)
	// for i := 0; i < int(meshColliderCount); i++ {
	// 	colliderData := new(MeshColliderData).Create()
	// 	binary.Read(f, binary.LittleEndian, colliderData.ObjectNameBytes)
	// 	colliderData.ObjectName = util.BytesToString(colliderData.ObjectNameBytes)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.ColliderLayer)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.IsConvex)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Location.Z)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.Z)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Rotation.W)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.X)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.Y)
	// 	binary.Read(f, binary.LittleEndian, &colliderData.Scale.Z)

	// 	vertexCount := int32(0)
	// 	binary.Read(f, binary.LittleEndian, &vertexCount)
	// 	vertexValueCount := int(vertexCount * 3)
	// 	colliderData.Vertices = make([]float32, vertexValueCount)
	// 	for j := 0; j < vertexValueCount; j++ {
	// 		binary.Read(f, binary.LittleEndian, &colliderData.Vertices[j])
	// 	}
	// 	indiceCount := int32(0)
	// 	binary.Read(f, binary.LittleEndian, &indiceCount)
	// 	colliderData.Indices = make([]int32, indiceCount)
	// 	for j := 0; j < int(indiceCount); j++ {
	// 		binary.Read(f, binary.LittleEndian, &colliderData.Indices[j])
	// 	}
	// 	this.MeshColliderDatas = append(this.MeshColliderDatas, colliderData)
	// 	//util.LoggerPtr().Debug("meshCollider colliderData.ObjectName=", string(colliderData.ObjectName))
	// }

	// 读取地形数据
	path := this.GetSceneDataPath() + this.Name + ".nav"
	//util.LoggerPtr().Debug("Start NavMesh")
	//navmesh
	if this.navMeshType != ENavmesh_No {
		addResult := config.NavMeshCache.AddOneSceneDataFormFile(path, this.Name, int(this.navMeshType))
		if addResult != 0 {
			util.LoggerPtr().Warning("[LUOHAO]preload NavFile form disk Fail", path, addResult)
		} else {
			util.LoggerPtr().Debug("[LUOHAO]preload NavFile form disk Success", path)
		}
	}

	if this.navMeshType == ENavmesh_Const {
		var result int

		this.autoPathService, result = new(recastnavgo.AutoPath).Create(config.NavMeshCache, this.Name)

		if result != 0 {
			util.LoggerPtr().Error("[LUOHAO] const autoPathService create fail result = ", this.Name, result)
		} else {
			util.LoggerPtr().Debug("[LUOHAO] const autoPathService create Success", this.Name)
		}
	}

	if this.WaitingMapSceneData != nil {
		this.WaitingMapSceneData.Load(config)
	}

	return true
}

func (this *SceneData) LoadJsonData(jsonDataPath string, dataList *LevelObjectData) bool {
	util.LoggerPtr().Info("load", jsonDataPath)
	b, err := ioutil.ReadFile(jsonDataPath)
	if err != nil {
		util.LoggerPtr().Error(err)
		return false
	}
	if err := json.Unmarshal(b, dataList); err != nil {
		util.LoggerPtr().Error(err)
		return false
	}
	return true
}

type SceneConfig struct {
	sceneDatas   map[uint32]*SceneData
	NavMeshCache recastnavgo.NavgationDataCache
}

func (this *SceneConfig) Create() *SceneConfig {
	// 创建数据缓存
	this.NavMeshCache = recastnavgo.NewNavgationDataCache()

	// 创建场景数据字典
	//init scene data cached data
	this.sceneDatas = make(map[uint32]*SceneData)
	// for _, mData := range MapConfigGet().GetMapDataList() {
	// 	if mData.IsOpen == true {

	// 		this.sceneDatas[mData.UniqueID] = new(SceneData).Create(mData.Name, mData.NavmeshType, mData.WaitingMapName)
	// 	}
	// }

	return this
}

func (this *SceneConfig) Load() bool {
	for _, sData := range this.sceneDatas {
		if sData.Load(this) == false {
			return false
		}
	}
	return true
}
func (this *SceneConfig) Cleanup() {
	this.sceneDatas = nil
}
