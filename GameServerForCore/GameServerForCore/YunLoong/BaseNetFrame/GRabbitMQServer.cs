using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace YunLoong
{
    public struct RabbitMQConfig
    {
        public string RabbitMQIP;
        public int RabbitMQPort;
        public string RabbitMQUser;
        public string RabbitMQPassword;
        public string RabbitMQvHost;
    }
    public abstract class GRabbitMQServer
    {
        public Dictionary<string, RabbitMQCtrl> RabbitMQCtrlDic;

        public static ConnectionFactory factory;

        public static IConnection connection;
        public void initConnect(RabbitMQConfig config)
        {
            RabbitMQCtrlDic = new Dictionary<string, RabbitMQCtrl>();
            factory = new ConnectionFactory();
            factory.HostName = config.RabbitMQIP;//RabbitMQ服务在本地运行
            factory.Port = config.RabbitMQPort;//RabbitMQ服务在本地运行
            factory.UserName = config.RabbitMQUser;//用户名
            factory.Password = config.RabbitMQPassword;//密码
            factory.VirtualHost = config.RabbitMQvHost;//RabbitMQ中的虚拟主机名
            Connect();
        }
        public void Connect()
        {
            if (connection != null)
            {
                connection.Dispose();
                connection = null;
            }
            connection = factory.CreateConnection();
        }
        public enum SubscribeType
        {
            ForList,
            ForMessage
        }

        //添加订阅
        public void addSubscribe(SubscribeType subscribeType, string exchange, string queue, Action<byte[]> ReceiveCallBack)
        {
            string key = exchange + "_" + queue;
            if (!RabbitMQCtrlDic.ContainsKey(key))
            {
                RabbitMQCtrl tRabbitMQCtrl = new RabbitMQCtrl();
                tRabbitMQCtrl.initRabbitMQCtrl(subscribeType, exchange, queue, ReceiveCallBack);
                RabbitMQCtrlDic.Add(key, tRabbitMQCtrl);
            }
        }
        //清除订阅
        public void removeSubscribe(string exchange, string queue)
        {
            string key = exchange + "_" + queue;
            if (RabbitMQCtrlDic.ContainsKey(key))
            {
                RabbitMQCtrlDic[key].Dispose();
                RabbitMQCtrlDic.Remove(key);
            }
        }
        //发布订阅(广播模式)
        public void sendPublish(string exchange, byte[] bytes)
        {
            using (var channel = connection.CreateModel())
            {
                channel.ExchangeDeclare(exchange, "fanout");
                channel.BasicPublish(exchange, "", null, bytes); //开始传递
            }
        }
        //发布消息队列(生产者消费者模式)
        public void sendPublish(string exchange, string queue, byte[] bytes)
        {
            using (var channel = connection.CreateModel())
            {
                channel.ExchangeDeclare(exchange, "direct");
                channel.QueueDeclare(queue, false, false, true, null);//创建消息队列
                channel.BasicPublish(exchange, queue, null, bytes); //开始传递
            }
        }
    }

    public class RabbitMQCtrl
    {
        public IConnection connection
        {
            get
            {
                return GRabbitMQServer.connection;
            }
        }
        public IModel channel;
        public Action<byte[]> ReceiveCallBack;
        public void initRabbitMQCtrl(GRabbitMQServer.SubscribeType subscribeType, string exchange, string queue, Action<byte[]> pReceiveCallBack)
        {
            ReceiveCallBack = pReceiveCallBack;
            channel = connection.CreateModel();
            channel.QueueDeclare(queue, false, false, true, null);

            switch (subscribeType)
            {
                case GRabbitMQServer.SubscribeType.ForList:
                    channel.ExchangeDeclare(exchange, "direct");
                    channel.QueueBind(queue, exchange, queue);
                    break;
                case GRabbitMQServer.SubscribeType.ForMessage:
                    channel.ExchangeDeclare(exchange, "fanout");
                    channel.QueueBind(queue, exchange, "");
                    break;
                default:
                    break;
            }
            var consumer = new EventingBasicConsumer(channel);
            channel.BasicConsume(queue, true, consumer);
            consumer.Received += Received_fun;
        }

        private void Received_fun(object sender, BasicDeliverEventArgs e)
        {
            if (ReceiveCallBack != null)
            {
                ReceiveCallBack(e.Body);
            }
        }

        public void Dispose(bool isEnd = false)
        {
            connection.Dispose();
            channel.Dispose();
            ReceiveCallBack = null;
        }

    }
    

}
