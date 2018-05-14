using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Timers;

namespace YunLoong
{
    public abstract class GSocketManager
    {

        public static Action CloseSocketEvent;
        private static bool _isCloseSocket;
        public static bool isCloseSocket
        {
            get
            {
                return _isCloseSocket;
            }
            set
            {
                _isCloseSocket = value;
                if (CloseSocketEvent != null)
                {
                    CloseSocketEvent();
                }

            }
        }
        public int port;
        public string host;
        ///创建终结点（EndPoint）
        public IPAddress ip;
        public IPEndPoint ipe;
        private Timer timer;
        private int timerCount = 0;
        public List<GSocket> socketlist;

        /// <summary>
        /// 客户端连接上的Session组合
        /// </summary>
        protected readonly ConcurrentDictionary<object, GSocket> connectSessions = new ConcurrentDictionary<object, GSocket>();
        protected readonly BagWithExpire<GSocket> connectBag = new BagWithExpire<GSocket>();

        public GSocketManager()
        {
            CloseSocketEvent += CloseSocket;
            timer = new System.Timers.Timer(10000);
            timer.Enabled = true;
            timer.AutoReset = true;
            timer.Elapsed += new System.Timers.ElapsedEventHandler(checksocket);
        }
        public abstract void CloseSocket();
        public void initSocket()
        {
            socketlist = new List<GSocket>();
            ip = IPAddress.Parse(host);//把ip地址字符串转换为IPAddress类型的实例
            ipe = new IPEndPoint(ip, port);//用指定的端口和ip初始化IPEndPoint类的新实例
            createsocketlistener();
        }
        public void createsocketlistener()
        {
            ///创建socket并开始监听
            Socket socketlistener = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);//创建一个socket对像，如果用udp协议，则要用SocketType.Dgram类型的套接字
            socketlistener.Bind(ipe);//绑定EndPoint对像（2000端口和ip地址）
            socketlistener.Listen(1024);//开始监听
            socketlistener.BeginAccept(new AsyncCallback(AcceptCallback), socketlistener);
        }
        public void AcceptCallback(IAsyncResult ar)
        {
            if (!isCloseSocket)
            {
                // Get the socket that handles the client request.
                Socket socketlistener = (Socket)ar.AsyncState;
                Socket handler = socketlistener.EndAccept(ar);

                GSocket session = new GSocket(handler);
                connectBag.Add(session);
                Console.WriteLine("新连接,等待链接验证");
                CreateSocketHandle(session);
                session.ReceiveAsync();
                socketlistener.BeginAccept(new AsyncCallback(AcceptCallback), socketlistener);

            }
        }
        /// <summary>
        /// 用于给指定的GSocket对象绑定SocketHandle对象
        /// </summary>
        /// <param name="session"></param>
        public abstract void CreateSocketHandle(GSocket session);
        //异常连接，清除并关闭
        public void CloseSession(GSocket session)
        {
            if (connectBag.TryTake(ref session))
            {
                Console.WriteLine("connectBag清除连接成功"+connectBag.Count);
                session.Dispose();
            }
        }
        public void CloseSession(object sid, GSocket session)
        {
            if (connectSessions.ContainsKey(sid))
            {
                connectSessions.TryRemove(sid, out session);
                Console.WriteLine("connectSessions清除连接成功"+sid+"数量:"+connectSessions.Count);
                session.Dispose();
            }
        }
        public virtual void SaveSession(long sid, GSocket session)
        {
            if (connectBag.TryTake(ref session))
            {
                if (!connectSessions.ContainsKey(sid))
                {
                    connectSessions.GetOrAdd(sid, session);
                }
                else
                {
                    connectSessions[sid].wSocket.Close();
                    connectSessions[sid] = session;
                }
                connectSessions[sid].id = sid;
                Console.WriteLine("connectSessions添加连接成功" + sid + "数量:" + connectSessions.Count);
            }
        }

        public void checksocket(object sender, System.Timers.ElapsedEventArgs e)
        {
            Console.WriteLine("检查Socket列表,广播心跳包:"+ connectSessions.Count);
            foreach(var item in connectSessions)
            {
                try
                {
                    item.Value.socketHandle.SendHeartBeat();
                }
                catch
                {
                    CloseSession(item.Key, item.Value);
                }
                //key.Value.socketHandle.SendRequest(MsgCode.MSG_CODE_HEARTBEAT_RET, GameData.gettimer());
            }
            timerCount++;
            if (timerCount == 6)
            {
                timerCount = 0;
                CheckForMinuteEvent();
            }
        }
        public abstract void CheckForMinuteEvent();
        public SocketHandle giveSocketHandle(int sid)
        {
            GSocket value;
            if(connectSessions.TryGetValue(sid, out value))
                return connectSessions[sid].socketHandle;
            return null;
        }

        public bool IsPlayerOnline(int userId)
        {
            return connectSessions.ContainsKey(userId);
        }

        //public void Foreach(Predicate<GSocket> match,Action<GSocket> action)
        //{
        //    ICollection<int> keys = connectSessions.Keys;
        //    var er = keys.GetEnumerator();
        //    while (er.MoveNext())
        //    {
        //        GSocket gs = null;
        //        if(connectSessions.TryGetValue(er.Current, out gs) && match(gs))
        //            action(gs);
        //    }
        //}
    }
}
