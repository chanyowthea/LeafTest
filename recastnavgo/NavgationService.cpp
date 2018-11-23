#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "NavgationService.h"
#include "Detour/DetourNavMeshQuery.h"
#include "Detour/DetourCommon.h"
#include "Detour/DetourNavMeshBuilder.h"
#include "DetourTileCache/DetourTileCache.h"
#include "DetourTileCache/DetourTileCacheBuilder.h"
#include "fastlz.h"

using namespace std;

#pragma region tilecahce tool struct
struct FastLZCompressor : public dtTileCacheCompressor
{
	virtual int maxCompressedSize(const int bufferSize)
	{
		return (int)(bufferSize* 1.05f);
	}

	virtual dtStatus compress(const unsigned char* buffer, const int bufferSize,
		unsigned char* compressed, const int /*maxCompressedSize*/, int* compressedSize)
	{
		*compressedSize = fastlz_compress((const void *const)buffer, bufferSize, compressed);
		return DT_SUCCESS;
	}

	virtual dtStatus decompress(const unsigned char* compressed, const int compressedSize,
		unsigned char* buffer, const int maxBufferSize, int* bufferSize)
	{
		*bufferSize = fastlz_decompress(compressed, compressedSize, buffer, maxBufferSize);
		return *bufferSize < 0 ? DT_FAILURE : DT_SUCCESS;
	}
};

struct LinearAllocator : public dtTileCacheAlloc
{
	unsigned char* buffer;
	size_t capacity;
	size_t top;
	size_t high;

	LinearAllocator(const size_t cap) : buffer(0), capacity(0), top(0), high(0)
	{
		resize(cap);
	}

	~LinearAllocator()
	{
		dtFree(buffer);
	}

	void resize(const size_t cap)
	{
		if (buffer) dtFree(buffer);
		buffer = (unsigned char*)dtAlloc(cap, DT_ALLOC_PERM);
		capacity = cap;
	}

	virtual void reset()
	{
		high = dtMax(high, top);
		top = 0;
	}

	virtual void* alloc(const size_t size)
	{
		if (!buffer)
			return 0;
		if (top + size > capacity)
			return 0;
		unsigned char* mem = &buffer[top];
		top += size;
		return mem;
	}

	virtual void free(void* /*ptr*/)
	{
		// Empty
	}
};

struct MeshProcess : public dtTileCacheMeshProcess
{
	inline MeshProcess()
	{
	}

	virtual void process(struct dtNavMeshCreateParams* params,
		unsigned char* polyAreas, unsigned short* polyFlags)
	{
		// Update poly flags from areas.
		for (int i = 0; i < params->polyCount; ++i)
		{
			if (polyAreas[i] == DT_TILECACHE_WALKABLE_AREA)
				polyAreas[i] = SAMPLE_POLYAREA_GROUND;

			if (polyAreas[i] == SAMPLE_POLYAREA_GROUND ||
				polyAreas[i] == SAMPLE_POLYAREA_GRASS ||
				polyAreas[i] == SAMPLE_POLYAREA_ROAD)
			{
				polyFlags[i] = SAMPLE_POLYFLAGS_WALK;
			}
			else if (polyAreas[i] == SAMPLE_POLYAREA_WATER)
			{
				polyFlags[i] = SAMPLE_POLYFLAGS_SWIM;
			}
			else if (polyAreas[i] == SAMPLE_POLYAREA_DOOR)
			{
				polyFlags[i] = SAMPLE_POLYFLAGS_WALK | SAMPLE_POLYFLAGS_DOOR;
			}
		}
	}
};

#pragma endregion


NavgationService::NavgationService()
{
	m_navMesh = NULL;
	m_navQuery = NULL;
	m_tileCache = NULL;
	navType = ENavType_TileMesh;

	m_talloc = new LinearAllocator(32000);
	m_tcomp = new FastLZCompressor;
	m_tmproc = new MeshProcess;
}


NavgationService::~NavgationService()
{
	if (m_navQuery)
	{
		dtFreeNavMeshQuery(m_navQuery);
	}
	m_navQuery = NULL;

	if (m_navMesh)
	{
		dtFreeNavMesh(m_navMesh);
	}
	m_navMesh = NULL;

	if (m_tileCache)
	{
		dtFreeTileCache(m_tileCache);
	}
	m_tileCache = NULL;

	if (m_talloc)
	{
		delete m_talloc;
		m_talloc = NULL;
	}

	if (m_tcomp)
	{
		delete m_tcomp;
		m_tcomp = NULL;
	}

	if (m_tmproc)
	{
		delete m_tmproc;
		m_tmproc = NULL;
	}

	gameObjToObstacleMap.clear();
}



int NavgationService::LoadNavMeshFromDataCache(NavgationDataCache* cacheManager, char* sceneName)
{
	if (cacheManager == NULL)
	{
		return -1;
	}

	string strName = sceneName;

	if (cacheManager->sceneNavDataMap.find(strName) != cacheManager->sceneNavDataMap.end())
	{
		return loadTileMeshFromDataCache(cacheManager->sceneNavDataMap[strName]);
	}
	else if (cacheManager->sceneTileCacheMap.find(strName) != cacheManager->sceneTileCacheMap.end())
	{
		return loadTileCacheFromDataCache(cacheManager->sceneTileCacheMap[strName]);
	}

	return -2;

}

#pragma region load const navmesh
int NavgationService::loadTileMeshFromDataCache(OneSceneNavMeshData*  sceneData)
{
	if (sceneData == NULL)
	{
		return 1;
	}

	navType = ENavType_TileMesh;

	dtFreeNavMesh(m_navMesh);

	m_navMesh = dtAllocNavMesh();

	dtNavMeshParams tmpParams;
	tmpParams.orig[0] = sceneData->meshParam.orig[0];
	tmpParams.orig[1] = sceneData->meshParam.orig[1];
	tmpParams.orig[2] = sceneData->meshParam.orig[2];
	tmpParams.tileWidth = sceneData->meshParam.tileWidth;
	tmpParams.tileHeight = sceneData->meshParam.tileHeight;
	tmpParams.maxTiles = sceneData->meshParam.maxTiles;
	tmpParams.maxPolys = sceneData->meshParam.maxPolys;

	dtStatus status = m_navMesh->init(&tmpParams);

	if (dtStatusFailed(status))
	{
		return 2;
	}

	for (unsigned int i = 0; i < sceneData->tileRefList.size(); i++)
	{
		if (i >= sceneData->tileDataSizeList.size()
			|| i >= sceneData->dataPtrList.size())
		{
			return 3;
		}
		int dataSize = sceneData->tileDataSizeList[i];

		unsigned char* data = (unsigned char*)dtAlloc(dataSize, DT_ALLOC_PERM);
		if (!data) break;
		memset(data, 0, dataSize);
		memcpy(data, sceneData->dataPtrList[i], dataSize);

		m_navMesh->addTile(data, dataSize, DT_TILE_FREE_DATA, sceneData->tileRefList[i], 0);
	}

	m_navQuery = dtAllocNavMeshQuery();
	dtStatus nResult = m_navQuery->init(m_navMesh, 2048);

	bool bResult = dtStatusSucceed(nResult);
	if (!bResult)
	{
		return 4;
	}

	return 0;
}
#pragma endregion


#pragma region load tileCache


int NavgationService::loadTileCacheFromDataCache(OneSceneTileCacheData* sceneData)
{
	if (sceneData == NULL)
	{
		return 1;
	}

	navType = ENavType_TileCache;

	//init navMesh
	dtFreeNavMesh(m_navMesh);
	m_navMesh = dtAllocNavMesh();

	dtNavMeshParams tmpParams;
	tmpParams.orig[0] = sceneData->meshParam.orig[0];
	tmpParams.orig[1] = sceneData->meshParam.orig[1];
	tmpParams.orig[2] = sceneData->meshParam.orig[2];
	tmpParams.tileWidth = sceneData->meshParam.tileWidth;
	tmpParams.tileHeight = sceneData->meshParam.tileHeight;
	tmpParams.maxTiles = sceneData->meshParam.maxTiles;
	tmpParams.maxPolys = sceneData->meshParam.maxPolys;

	dtStatus status = m_navMesh->init(&tmpParams);

	if (dtStatusFailed(status))
	{
		return 2;
	}

	//init tilecache
	dtFreeTileCache(m_tileCache);

	m_tileCache = dtAllocTileCache();

	if (!m_tileCache)
	{
		return 3;
	}

	dtTileCacheParams tmpTileCacheParams;

	tmpTileCacheParams.orig[0] = sceneData->cacheParams.orig[0];
	tmpTileCacheParams.orig[1] = sceneData->cacheParams.orig[1];
	tmpTileCacheParams.orig[2] = sceneData->cacheParams.orig[2];

	tmpTileCacheParams.cs = sceneData->cacheParams.cs;
	tmpTileCacheParams.ch = sceneData->cacheParams.ch;
	tmpTileCacheParams.width = sceneData->cacheParams.width;
    tmpTileCacheParams.height = sceneData->cacheParams.height;
	tmpTileCacheParams.walkableHeight = sceneData->cacheParams.walkableHeight;
	tmpTileCacheParams.walkableRadius = sceneData->cacheParams.walkableRadius;
	tmpTileCacheParams.walkableClimb = sceneData->cacheParams.walkableClimb;
	tmpTileCacheParams.maxSimplificationError = sceneData->cacheParams.maxSimplificationError;
	tmpTileCacheParams.maxTiles = sceneData->cacheParams.maxTiles;
	tmpTileCacheParams.maxObstacles = sceneData->cacheParams.maxObstacles;

	status = m_tileCache->init(&tmpTileCacheParams, m_talloc, m_tcomp, m_tmproc);
	if (dtStatusFailed(status))
	{
		return 4;
	}

	for (unsigned int i = 0; i < sceneData->dataPtrList.size(); i++)
	{
		if (i >= sceneData->tileDataSizeList.size()
			|| i >= sceneData->dataPtrList.size())
		{
			return 5;
		}
		int dataSize = sceneData->tileDataSizeList[i];

		unsigned char* data = (unsigned char*)dtAlloc(dataSize, DT_ALLOC_PERM);
		if (!data) break;
		memset(data, 0, dataSize);
		memcpy(data, sceneData->dataPtrList[i], dataSize);

		dtCompressedTileRef tile = 0;
		dtStatus addTileStatus = m_tileCache->addTile(data, dataSize, DT_COMPRESSEDTILE_FREE_DATA, &tile);
		if (dtStatusFailed(addTileStatus))
		{
			dtFree(data);
		}

		if (tile)
			m_tileCache->buildNavMeshTile(tile, m_navMesh);
	}

	m_navQuery = dtAllocNavMeshQuery();
	dtStatus nResult = m_navQuery->init(m_navMesh, 2048);

	bool bResult = dtStatusSucceed(nResult);
	if (!bResult)
	{
		return 6;
	}

	return 0;
}
#pragma endregion



#pragma region findpath
bool NavgationService::FindStraightPath(float* vStart, float* vEnd, float* pathArray, int* nodeCount)
{
	if (m_navQuery == 0 || m_navMesh == 0)
		return false;

	dtQueryFilter m_filter;
	m_filter.setIncludeFlags(SAMPLE_POLYFLAGS_WALK); // TODO

	static const int MAX_POLYS = 256;
	float  m_aStraightNavPath[MAX_POLYS * 3];


	float aPolyPickExt[3] = { 0 };
	aPolyPickExt[0] = 2;
	aPolyPickExt[1] = 2;
	aPolyPickExt[2] = 2;

	dtPolyRef cStartRef = 0;
	dtPolyRef cEndRef = 0;

	dtStatus nResult = m_navQuery->findNearestPoly(vStart, aPolyPickExt, &m_filter, &cStartRef, 0);
	if (dtStatusFailed(nResult)
		|| !m_navMesh->isValidPolyRef(cStartRef))
		return false;

	nResult = m_navQuery->findNearestPoly(vEnd, aPolyPickExt, &m_filter, &cEndRef, 0);
	if (dtStatusFailed(nResult)
		|| !m_navMesh->isValidPolyRef(cEndRef))
		return false;

	dtPolyRef aNavPolys[MAX_POLYS] = { 0 };
	int nNavPolys = 0;
	nResult = m_navQuery->findPath(cStartRef, cEndRef, vStart, vEnd, &m_filter, aNavPolys, &nNavPolys, MAX_POLYS);
	if (dtStatusFailed(nResult))
		return false;

	if (nNavPolys)
	{
		// In case of partial path, make sure the end point is clamped to the last polygon.
		float epos[3] = { 0 };
		dtVcopy(epos, vEnd);
		if (aNavPolys[nNavPolys - 1] != cEndRef)
			m_navQuery->closestPointOnPoly(aNavPolys[nNavPolys - 1], vEnd, epos, NULL);

		unsigned char aStraightPathFlags[MAX_POLYS] = { 0 };
		dtPolyRef     aStraightPathPolys[MAX_POLYS] = { 0 };
		int nPathNodeCount = 0;

		nResult = m_navQuery->findStraightPath(
			vStart,
			epos,
			aNavPolys,
			nNavPolys,
			m_aStraightNavPath,
			aStraightPathFlags,
			aStraightPathPolys,
			(int*)&nPathNodeCount,
			MAX_POLYS);


		for (int i = 0; i < nPathNodeCount * 3; i++)
		{
			pathArray[i] = m_aStraightNavPath[i];
		}
		*nodeCount = nPathNodeCount;

		return bool(nResult);
	}
	return false;
}

static float frand()
{
	int iMax = RAND_MAX > 32767 ? 32767 : RAND_MAX;
	int iRand = rand() % iMax;
	return (float)iRand / iMax;
}

bool NavgationService::FindRandPointOnNearestPoly(float* vFindCenter, float fFindRadius, int iWantPointNum, float* fFindNodeArray, int* iFindPointNum)
{
	if (m_navQuery == 0 || m_navMesh == 0)
	{
		return false;
	}

	dtPolyRef cStartRef = 0;
	float aPolyPickExt[3] = { 2, 2, 2 };
	dtQueryFilter queryFilter;
	queryFilter.setIncludeFlags(SAMPLE_POLYFLAGS_WALK); // TODO

	dtStatus nResult = m_navQuery->findNearestPoly(vFindCenter, aPolyPickExt, &queryFilter, &cStartRef, 0);
	if (dtStatusFailed(nResult)
		|| !m_navMesh->isValidPolyRef(cStartRef))
		return false;

	static const int MAX_POINT = 256;

	iWantPointNum = iWantPointNum > MAX_POINT ? MAX_POINT : iWantPointNum;

	int iMaxTry = iWantPointNum * 16;

	for (int i = 0; i < iMaxTry && *iFindPointNum < iWantPointNum; ++i)
	{
		dtPolyRef cRandRef = 0;
		float findRadomPoint[3];
		float resDis = 0;

		nResult = m_navQuery->findRandomPointAroundCircle(cStartRef, vFindCenter, fFindRadius,
			&queryFilter, &frand,
			&cRandRef, findRadomPoint);

		if (dtStatusSucceed(nResult))
		{
			resDis = dtVdist(findRadomPoint, vFindCenter);

			if (resDis <= fFindRadius)
			{
				dtVcopy(&fFindNodeArray[*iFindPointNum * 3], findRadomPoint);
				(*iFindPointNum)++;
			}

		}
	}

	if (*iFindPointNum > 0)
	{
		return true;
	}

	return false;
}
bool NavgationService::FindClosestPointOnNearestPoly(float* vFindCenter, float fFindRadius, float* fFindNode)
{
	if (m_navQuery == 0 || m_navMesh == 0)
	{
		return false;
	}

	dtPolyRef cStartRef = 0;
	float aPolyPickExt[3] = { 2, 2, 2 };
	dtQueryFilter queryFilter;
	queryFilter.setIncludeFlags(SAMPLE_POLYFLAGS_WALK); // TODO

	dtStatus nResult = m_navQuery->findNearestPoly(vFindCenter, aPolyPickExt, &queryFilter, &cStartRef, 0);
	if (dtStatusFailed(nResult)
		|| !m_navMesh->isValidPolyRef(cStartRef))
		return false;
	
	float findRadomPoint[3];
	bool posOverPoly;
	nResult = m_navQuery->closestPointOnPoly(cStartRef, vFindCenter, findRadomPoint, &posOverPoly);

	if (dtStatusSucceed(nResult))
	{
		float resDis = dtVdist(findRadomPoint, vFindCenter);

		if (resDis <= fFindRadius)
		{
			dtVcopy(fFindNode, findRadomPoint);
			return true;
		}

	}

	return false;
}
#pragma endregion

#pragma region obstacle
bool NavgationService::AddObstacle(int gameObjectID, const float* pos, const float radius, const float height)
{
	if (!m_tileCache)
		return false;

	if (gameObjToObstacleMap.find(gameObjectID) != gameObjToObstacleMap.end())
	{
		return false;
	}

	float p[3];
	dtVcopy(p, pos);
	p[1] -= height / 2;

	dtObstacleRef obID = 0;

	dtStatus result = m_tileCache->addObstacle(p, radius, height, &obID);

	gameObjToObstacleMap[gameObjectID] = obID;

	return bool(result);
}

bool NavgationService::RemoveObstacle(int gameObjectID)
{
	if (!m_tileCache)
	{
		return false;
	}

	if (gameObjToObstacleMap.find(gameObjectID) == gameObjToObstacleMap.end())
	{
		return false;
	}

	dtStatus result = m_tileCache->removeObstacle(gameObjToObstacleMap[gameObjectID]);

	gameObjToObstacleMap.erase(gameObjectID);

	return bool(result);
}

bool NavgationService::UpdateNavmesh(bool* upToDate)
{
	if (navType != ENavType_TileCache)
	{
		return true;
	}

	if (!m_navMesh)
		return false;
	if (!m_tileCache)
		return false;

	float dt = 0;

	return m_tileCache->update(dt, m_navMesh, upToDate);
}

#pragma endregion


void NavgationService::ThreadEnd()
{
    if(m_navQuery)
    {
        m_navQuery->ThreadEnd();
    }
}