using System;
using System.Collections;


namespace GameServer.CsScript
{
    public class GameDataConfig
    {
        public static GameDataConfig Instance;
        static GameDataConfig()
        {
            Instance = new GameDataConfig();
        }

        public int ServerId;
        public string IPPort;
        //public string AccountServer;
        //public string GameServer;
        //public string ProductSignKey;
        //public string ClientDesKey;

        public string Company;
        public string RedisIPPort;
        public string RabbitMQIPPort;
        public string RabbitMQvHost;
        public string RabbitMQUser;
        public string RabbitMQPassword;

        public string IP;
        public int Port;

        public string RabbitMQIP;
        public int RabbitMQPort;

        public void init()
        {
            Hashtable thash = GameResources.Instance.jsonhash["gameDataConfig"];
            ServerId = Convert.ToInt32(thash["ServerId"]);
            //AccountServer = Convert.ToString(thash["AccountServer"]);
            //GameServer = Convert.ToString(thash["GameServer"]);
            //string[] IPPort_arr = GameServer.Split(':');
            IPPort = Convert.ToString(thash["IPPort"]);
            string[] IPPort_arr = IPPort.Split(':');
            IP = IPPort_arr[0];
            Port = Convert.ToInt32(IPPort_arr[1]);

            ///ProductSignKey = Convert.ToString(thash["ProductSignKey"]);
            //ClientDesKey = Convert.ToString(thash["ClientDesKey"]);
            Company = Convert.ToString(thash["Company"]);
            RedisIPPort = Convert.ToString(thash["RedisIPPort"]);

            RabbitMQIPPort = Convert.ToString(thash["RabbitMQIPPort"]);
            IPPort_arr = RabbitMQIPPort.Split(':');
            RabbitMQIP = IPPort_arr[0];
            RabbitMQPort = Convert.ToInt32(IPPort_arr[1]);
            RabbitMQvHost = Convert.ToString(thash["RabbitMQvHost"]);
            RabbitMQUser = Convert.ToString(thash["RabbitMQUser"]);
            RabbitMQPassword = Convert.ToString(thash["RabbitMQPassword"]);
        }
    }
}