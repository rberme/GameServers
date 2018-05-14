using System;
using System.Collections.Generic;
using GameServer.CsScript.Model;
using GameServer.Model;
using YunLoong;
namespace GameServer.CsScript
{
    public class PlayerCacheLock : BaseCacheLock<PlayerCacheStruct>
    {
        public static Type mainEntity { get { return typeof(PlayerData); } }

        public PlayerCacheLock()
        {
            //changeDataSet()
        }
        //逻辑处理
        public void readDataHandle(Action handleAction)
        {
            using (Read())
            {
                handleAction();
            }
        }
        public void writeDataHandle(Action handleAction)
        {
            using (Write())
            {
                handleAction();
            }
        }
        public MyPlayerData giveMyPlayerData() {
            using (Read())
            {
                MyPlayerData tMyPlayerData = GameData.CloneModel<MyPlayerData>(Data.playerData);
                return tMyPlayerData;
            }
        }
    }
    public class PlayerCacheStruct : BaseCacheStruct
    {
        //PlayerData(表记录)
        //HeroGroup(Dic)--HeroData(表记录)
        //ItemGroup(Dic)--ItemData(表记录)
        public PlayerData playerData;
        public Dictionary<int, PlayerHero> playerHeroDic;
        public Dictionary<int, PlayerItem> playerItemDic;
        public Dictionary<int, PlayerRock> playerRockDic;
        public Dictionary<int, PlayerEquip> playerEquipDic;

        public long pCacheID { get { return Convert.ToInt64(cacheID); } }
        public override GSqlManager gSqlManager { get; set; }

        public PlayerCacheStruct()
        {
            gSqlManager = new GSqlManager(GMySqlInfo.getGServerInfo());

            playerData = new PlayerData();
            playerHeroDic = new Dictionary<int, PlayerHero>();
            playerItemDic = new Dictionary<int, PlayerItem>();
            playerRockDic = new Dictionary<int, PlayerRock>();
            playerEquipDic = new Dictionary<int, PlayerEquip>();
        }
        //初始化用户
        public override void initCacheData()
        {
            gSqlManager.TryGiveResultSet(cacheID, out playerData);

            Console.WriteLine($"{cacheID}:初始化数据{ playerData.pid}");
            //初始化数据对象并加载
            base.initCacheData();
        }
        //创建用户
        public override void createCacheData()
        {
            playerData = DefaultModelData.defaultPlayerData((int)(pCacheID/100000), GameDataConfig.Instance.ServerId);
            gSqlManager.mySqlDataProvider.saveData(playerData);
            //创建数据对象并保存
            base.createCacheData();
        }



    }
}
