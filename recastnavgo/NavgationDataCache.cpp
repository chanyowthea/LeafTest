#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "NavgationDataCache.h"
#include "Detour/DetourNavMesh.h"

using namespace std;


NavgationDataCache::NavgationDataCache()
{
}

NavgationDataCache::~NavgationDataCache()
{
	SceneDataMap::iterator iter;

	for (iter = sceneNavDataMap.begin(); iter != sceneNavDataMap.end(); iter++)
	{
		for (unsigned int i = 0; i < iter->second->dataPtrList.size(); i++)
		{
			if (iter->second->dataPtrList[i] != NULL)
			{
				delete iter->second->dataPtrList[i];
				iter->second->dataPtrList[i] = NULL;
			}
		}
	}

	sceneNavDataMap.clear();

	SceneDataTileCacheMap::iterator tileIter;

	for (tileIter = sceneTileCacheMap.begin(); tileIter != sceneTileCacheMap.end(); iter++)
	{
		for (unsigned int i = 0; i < tileIter->second->dataPtrList.size(); i++)
		{
			if (tileIter->second->dataPtrList[i] != NULL)
			{
				delete tileIter->second->dataPtrList[i];
				tileIter->second->dataPtrList[i] = NULL;
			}
		}
	}

	sceneTileCacheMap.clear();
}

static const int NAVMESHSET_MAGIC = 'M' << 24 | 'S' << 16 | 'E' << 8 | 'T'; //'MSET';
static const int NAVMESHSET_VERSION = 1;

static const int TILECACHESET_MAGIC = 'T' << 24 | 'S' << 16 | 'E' << 8 | 'T'; //'TSET';
static const int TILECACHESET_VERSION = 1;

int NavgationDataCache::AddOneSceneDataFormFile(const char* path, char* sceneName, int NavType)
{

	if ( NavType == ENavType_TileMesh)
	{
		return AddTileMeshFromFile(path, sceneName);
	}
	else
	{
		return AddTileCacheFromFile(path, sceneName);
	}

}

int NavgationDataCache::AddTileMeshFromFile(const char* path, char* sceneName)
{
	FILE* fp = fopen(path, "rb");
	if (!fp) return 1;

	// Read header.
	NavMeshSetHeader header;
	size_t readLen = fread(&header, sizeof(NavMeshSetHeader), 1, fp);

	if (readLen != 1)
	{
		fclose(fp);
		return 2;
	}
	if (header.magic != NAVMESHSET_MAGIC)
	{
		fclose(fp);
		return 3;
	}
	if (header.version != NAVMESHSET_VERSION)
	{
		fclose(fp);
		return 4;
	}

	OneSceneNavMeshData* oneSceneData = new OneSceneNavMeshData();
	oneSceneData->meshParam = header.params;


	// Read tiles.
	for (int i = 0; i < header.numTiles; ++i)
	{
		NavMeshTileHeader tileHeader;
		readLen = fread(&tileHeader, sizeof(tileHeader), 1, fp);
		if (readLen != 1)
		{
			fclose(fp);
			return 5;
		}

		if (!tileHeader.tileRef || !tileHeader.dataSize)
			break;

		oneSceneData->tileRefList.push_back(tileHeader.tileRef);
		oneSceneData->tileDataSizeList.push_back(tileHeader.dataSize);

		unsigned char* data = (unsigned char*)dtAlloc(tileHeader.dataSize, DT_ALLOC_PERM);
		if (!data) break;
		memset(data, 0, tileHeader.dataSize);
		readLen = fread(data, tileHeader.dataSize, 1, fp);

		oneSceneData->dataPtrList.push_back(data);

		if (readLen != 1)
		{
			fclose(fp);
			return 6;
		}
	}

	string strName = sceneName;
	sceneNavDataMap[strName] = oneSceneData;

	fclose(fp);

	return 0;
}

int NavgationDataCache::AddTileCacheFromFile(const char* path, char* sceneName)
{
	FILE* fp = fopen(path, "rb");
	if (!fp) return 1;

	// Read header.
	TileCacheSetHeader header;
	size_t headerReadReturnCode = fread(&header, sizeof(TileCacheSetHeader), 1, fp);

	if (headerReadReturnCode != 1)
	{
		// Error or early EOF
		fclose(fp);
		return 2;
	}
	if (header.magic != TILECACHESET_MAGIC)
	{
		fclose(fp);
		return 3;
	}
	if (header.version != TILECACHESET_VERSION)
	{
		fclose(fp);
		return 4;
	}

	OneSceneTileCacheData* oneTileCacheData = new OneSceneTileCacheData();
	oneTileCacheData->meshParam = header.meshParams;
	oneTileCacheData->cacheParams = header.cacheParams;

	// Read tiles.
	for (int i = 0; i < header.numTiles; ++i)
	{
		TileCacheTileHeader tileHeader;
		size_t tileHeaderReadReturnCode = fread(&tileHeader, sizeof(tileHeader), 1, fp);
		if (tileHeaderReadReturnCode != 1)
		{
			// Error or early EOF
			fclose(fp);
			return 5;
		}
		if (!tileHeader.tileRef || !tileHeader.dataSize)
			break;

		oneTileCacheData->tileDataSizeList.push_back(tileHeader.dataSize);

		unsigned char* data = (unsigned char*)dtAlloc(tileHeader.dataSize, DT_ALLOC_PERM);
		if (!data) break;
		memset(data, 0, tileHeader.dataSize);
		size_t tileDataReadReturnCode = fread(data, tileHeader.dataSize, 1, fp);
		if (tileDataReadReturnCode != 1)
		{
			// Error or early EOF
			dtFree(data);
			fclose(fp);
			return 6;
		}

		oneTileCacheData->dataPtrList.push_back(data);
	}

    sceneTileCacheMap[sceneName] = oneTileCacheData;
	fclose(fp);
	return 0;
}
