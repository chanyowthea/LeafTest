#pragma once

#include <vector>
#include <map>
#include <string>

#define DT_CachPOLYREF64 1

#ifdef DT_CachPOLYREF64
#include <stdint.h>
typedef uint64_t dtCachePolyRef;
#else
typedef unsigned int dtCachePolyRef;
#endif


enum NavMeshType
{
	ENavType_TileMesh = 1,
	ENavType_TileCache = 2
};

#pragma region TileMesh
struct NavMeshParams
{
	float orig[3];					///< The world space origin of the navigation mesh's tile space. [(x, y, z)]
	float tileWidth;				///< The width of each tile. (Along the x-axis.)
	float tileHeight;				///< The height of each tile. (Along the z-axis.)
	int maxTiles;					///< The maximum number of tiles the navigation mesh can contain.
	int maxPolys;					///< The maximum number of polygons each tile can contain.
};

struct NavMeshSetHeader
{
	int magic;
	int version;
	int numTiles;

	NavMeshParams params;
};

struct NavMeshTileHeader
{
	dtCachePolyRef tileRef;
	int dataSize;
};


struct  OneSceneNavMeshData
{
	NavMeshParams meshParam;

	std::vector<dtCachePolyRef> tileRefList;

	std::vector<int> tileDataSizeList;

	std::vector<unsigned char*> dataPtrList;
};

typedef std::map<std::string, OneSceneNavMeshData*> SceneDataMap;
#pragma endregion

#pragma region TileCache
struct TileCacheParams
{
	float orig[3];
	float cs, ch;
	int width, height;
	float walkableHeight;
	float walkableRadius;
	float walkableClimb;
	float maxSimplificationError;
	int maxTiles;
	int maxObstacles;
};

struct TileCacheSetHeader
{
	int magic;
	int version;
	int numTiles;
	NavMeshParams meshParams;
	TileCacheParams cacheParams;
};

struct TileCacheTileHeader
{
	dtCachePolyRef tileRef;
	int dataSize;
};

struct  OneSceneTileCacheData
{
	NavMeshParams meshParam;

	TileCacheParams cacheParams;

	std::vector<int> tileDataSizeList;

	std::vector<unsigned char*> dataPtrList;
};

typedef std::map<std::string, OneSceneTileCacheData*> SceneDataTileCacheMap;

#pragma endregion

class NavgationDataCache
{
public:
	NavgationDataCache();
	~NavgationDataCache();

	int AddOneSceneDataFormFile(const char* path, char* sceneName, int NavType);

	SceneDataMap sceneNavDataMap;
	int AddTileMeshFromFile(const char* path, char* sceneName);

	SceneDataTileCacheMap sceneTileCacheMap;
	int AddTileCacheFromFile(const char* path, char* sceneName);
};

