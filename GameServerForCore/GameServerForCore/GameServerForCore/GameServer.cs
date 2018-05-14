using System;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Diagnostics;
using System.Collections.Generic;
using System.Collections.Concurrent;
using GameServer.CsScript;
using YunLoong;

namespace GameServerForCore
{

    class GameServer
    {
        static void Main(string[] args)
        {
            Console.OutputEncoding = System.Text.Encoding.UTF8;//设置当前环境代码页为UTF8，解决中文乱码

            GameResources.Instance.startLoadRes();
            GameDataConfig.Instance.init();
            TableData.Instance.initTableData();

            PlayerCacheManager tPlayerCacheManager = PlayerCacheManager.instance;
            GameRabbitMQServer tGameRabbitMQServer = GameRabbitMQServer.Instance;
            PlayerSocketManager tPlayerSocketManager = PlayerSocketManager.Instance;


            GNetManager.TestGNetManager();

            while (true)
            {
                string operation = Console.ReadLine();
                if (operation == "quit")
                {
                    Environment.Exit(0);
                }
                if (operation == "close")
                {
                    GSocketManager.isCloseSocket = true;
                }
                if (operation == "cls")
                {
                    Console.Clear();
                }

            }


        }
    }

}
