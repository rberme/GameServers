using System;
using System.Collections.Generic;
using YunLoong;
namespace GameServer.CsScript
{
    public class PlayerSocketManager : ServerSocketManager
    {
        public static PlayerSocketManager Instance;

        //public static Action minuteEvent;
        static PlayerSocketManager()
        {
            Instance = new PlayerSocketManager();
        }
        
        public PlayerSocketManager()
        {
            init();
        }
        public void init()
        {
            host = GameDataConfig.Instance.IP;
            port = GameDataConfig.Instance.Port;
            initSocket();
        }
        public override void CreateSocketHandle(GSocket session)
        {
            session.socketHandle = new PlayerSocketHandle();
            session.socketHandle.GameSocket = session;
        }

        public override void CheckForMinuteEvent()
        {
            //当前在线玩家ID
            List<object> keys = new List<object>(connectSessions.Keys);
            //当前缓存玩家ID
            List<object> cacheKeys = new List<object>(PlayerCacheManager.instance.CacheLockDic.Keys);

            //清除过期的玩家数据对象
            //if (minuteEvent != null)
            //{
            //    minuteEvent();
            //}
            //GameMain.Instance.CheckAndRemovePlayer();
        }
    }

}
