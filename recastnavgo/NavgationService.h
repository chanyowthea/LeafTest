#pragma once
#include "NavgationDataCache.h"

class dtNavMesh;
class dtNavMeshQuery;
class dtTileCache;

/// These are just sample areas to use consistent values across the samples.
/// The use should specify these base on his needs.
enum SamplePolyAreas
{
	SAMPLE_POLYAREA_GROUND,
	SAMPLE_POLYAREA_WATER,
	SAMPLE_POLYAREA_ROAD,
	SAMPLE_POLYAREA_DOOR,
	SAMPLE_POLYAREA_GRASS,
	SAMPLE_POLYAREA_JUMP,
};
enum SamplePolyFlags
{
	SAMPLE_POLYFLAGS_WALK = 0x01,		// Ability to walk (ground, grass, road)
	SAMPLE_POLYFLAGS_SWIM = 0x02,		// Ability to swim (water).
	SAMPLE_POLYFLAGS_DOOR = 0x04,		// Ability to move through doors.
	SAMPLE_POLYFLAGS_JUMP = 0x08,		// Ability to jump.
	SAMPLE_POLYFLAGS_DISABLED = 0x10,		// Disabled polygon
	SAMPLE_POLYFLAGS_ALL = 0xffff	// All abilities.
};

typedef std::map<int, unsigned int> GameObjectToObstacleMap;

class NavgationService
{
public:
	NavgationService();
	~NavgationService();

	int LoadNavMeshFromDataCache(NavgationDataCache* cacheManager, char* sceneName);


	bool FindStraightPath(float* vStart, float* vEnd, float* pathNodeArray, int* nodeCount);
	bool FindRandPointOnNearestPoly(float* vFindCenter, float fFindRadius, int iWantPointNum, float* fFindNodeArray, int* iFindPointNum);
	bool FindClosestPointOnNearestPoly(float* vFindCenter, float fFindRadius, float* fFindNode);

#pragma region obstacle
	bool AddObstacle(int gameObjectID, const float* pos, const float radius, const float height);

	bool RemoveObstacle(int gameObjectID);

	bool UpdateNavmesh(bool* upToDate);

	void ThreadEnd();
#pragma endregion

protected:
	class dtNavMesh* m_navMesh;
	class dtNavMeshQuery* m_navQuery;
	dtTileCache* m_tileCache;
	NavMeshType navType;

	struct LinearAllocator* m_talloc;
	struct FastLZCompressor* m_tcomp;
	struct MeshProcess* m_tmproc;

	GameObjectToObstacleMap gameObjToObstacleMap;

private:
	int loadTileMeshFromDataCache(OneSceneNavMeshData*  sceneData);

	int loadTileCacheFromDataCache(OneSceneTileCacheData* sceneData);
};

