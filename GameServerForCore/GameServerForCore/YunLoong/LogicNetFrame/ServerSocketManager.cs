using System;

namespace YunLoong
{
    public abstract class ServerSocketManager : GSocketManager
    {
        public override void CloseSocket()
        {
            Console.WriteLine("ServerSocketManager" + "收到关闭连接处理");
            object[] keys = new object[connectSessions.Keys.Count];
            connectSessions.Keys.CopyTo(keys, 0);
            for (int i = 0; i < keys.Length; i++)
            {
                GSocket session;
                if(connectSessions.TryGetValue(keys[i], out session))
                {
                    session.socketHandle.SendByteMsgNow((int)BaseMsgCode.MSG_CODE_CLOSESOCKET_RET, 0);
                    CloseSession(keys[i], session);
                }
                //connectSessions.TryRemove(i, out session);
                //session.Close();
            }
            GSocket[] values = new GSocket[connectBag.Count];
            connectBag.CopyTo(values, 0);
            for (int i = 0; i < values.Length; i++)
            {
                connectBag.TryTake(ref values[i]);
                values[i].Dispose();
            }
        }
    }
}
