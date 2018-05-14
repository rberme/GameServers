using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using YunLoong;

namespace GameServer.CsScript
{
    public class GameRabbitMQServer : GRabbitMQServer
    {
        string Exchange = "Game2Chat";
        public static GameRabbitMQServer Instance;
        static GameRabbitMQServer()
        {
            Instance = new GameRabbitMQServer();
        }
        public GameRabbitMQServer()
        {
            RabbitMQConfig pRabbitMQConfig = new RabbitMQConfig()
            {
                RabbitMQIP = GameDataConfig.Instance.RabbitMQIP,
                RabbitMQPort = GameDataConfig.Instance.RabbitMQPort,
                RabbitMQUser = GameDataConfig.Instance.RabbitMQUser,
                RabbitMQPassword = GameDataConfig.Instance.RabbitMQPassword,
                RabbitMQvHost = GameDataConfig.Instance.RabbitMQvHost
            };
            initConnect(pRabbitMQConfig);
            //addSubscribe(SubscribeType.ForList, "GMToolsMQ", "GMToolsHandle" + GameDataConfig.Instance.ServerId, GMToolsHandle_fun);//添加私有GM操作MQ
            //addSubscribe(SubscribeType.ForMessage, "GMSpreadMQ", "GMToolsHandle" + GameDataConfig.Instance.ServerId, GMToolsHandle_fun);//添加广播GM操作MQ
            //addSubscribe(SubscribeType.ForList, "GameServerMQ", "GameServer" + GameDataConfig.Instance.ServerId, GameServerHandle_fun);//添加私有通信MQ
            addSubscribe(SubscribeType.ForList, Exchange, "GameServer" + GameDataConfig.Instance.ServerId, GameServerHandle_fun);//添加私有通信MQ

            //sendPublish("GMSpreadMQ", Encoding.UTF8.GetBytes("GMSpreadMQ"));
            //sendPublish("GMToolsMQ", "GMToolsHandle" + GameDataConfig.Instance.ServerId, Encoding.UTF8.GetBytes("GMToolsMQ"));
            //sendPublish("GameServerMQ", "GameServer" + GameDataConfig.Instance.ServerId, Encoding.UTF8.GetBytes("GameServerMQ"));


        }


        private void GameServerHandle_fun(byte[] bytes)
        {
            BinaryReader br = new BinaryReader(new MemoryStream(bytes));
            int msgCode = br.ReadInt32();
            //接收来自游戏服务器
            Console.WriteLine("已接收游戏服务器：" + msgCode);
            try
            {
                //switch (msgCode)
                //{
                //    case MsgCode.MSG_GAMESERVER_PLAYERDATA_RET:
                //        
                //        break;
                //    default:
                //        break;
                //}
            }
            catch (Exception ex)
            {
                Console.WriteLine("参数异常: " + ex.ToString());
            }
        }

    }
}
