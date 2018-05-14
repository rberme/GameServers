using System;
using System.Collections.Generic;
using GameServer.CsScript.Model;
using YunLoong;
namespace GameServer.CsScript
{
    public class PlayerCacheManager : BaseCacheManager<PlayerCacheLock, PlayerCacheStruct>
    {
        public static PlayerCacheManager instance;
        static PlayerCacheManager()
        {
            instance = new PlayerCacheManager();
        }
        public PlayerCacheManager() : base(GMySqlInfo.getGServerInfo())
        {
        }
        public override void initkeyHashSet()
        {
            gSqlManager.mySqlDataProvider.giveKeyForField<PlayerData>(out keyHashSet, "pid");
            foreach (var item in keyHashSet)
            {
                //Console.WriteLine(item);
            }

        }

    }
}
