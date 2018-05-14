using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Threading;

namespace YunLoong
{
    
    public abstract class ServerSocketHandle<T> : SocketHandle
    {
        public bool goodSocket;
        public T handleObject;
        public object userid;
        public Timer tempTimer = null;

        
        public override void SendCallback(int code)
        {
            if (code == (int)BaseMsgCode.MSG_CODE_HEARTBEAT_RET)
            {
                changeUserTime();//修改指定player的最后通信时间戳（下线时间戳），id
            }
        }
        public abstract void changeUserTime();
        public abstract string getMsgHandleFuncName(int msgid);
        public override void SendHeartBeat()
        {
            SendRequest((int)BaseMsgCode.MSG_CODE_HEARTBEAT_RET, GameData.gettimer());
        }

        public override void ReadMsgHandle(BinaryReader br)
        {

            int msgid = br.ReadInt32();
            int state = br.ReadByte();
            if (state == 0)
            {
                string msgHandleFuncName = getMsgHandleFuncName(msgid);
                MethodInfo tMethodInfo = handleObject.GetType().GetMethod(msgHandleFuncName);
                if (tMethodInfo != null)
                {
                    int datalen = br.ReadInt32();
                    if (datalen > 0)
                    {
                        byte[] databytes = br.ReadBytes(datalen);
                        try
                        {
                            //需要传参的协议处理
                            tMethodInfo.Invoke(handleObject, new object[] { databytes });
                        }
                        catch (Exception e)
                        {
                            Console.WriteLine("收到客户端协议:" + msgid + "出错,处理方法异常:" + msgHandleFuncName + ",Error:" + e.ToString());
                        }
                    }
                    else
                    {
                        try
                        {
                            //无参数的协议处理
                            tMethodInfo.Invoke(handleObject, null);
                        }
                        catch (Exception e)
                        {
                            Console.WriteLine("收到客户端协议:" + msgid + "出错,处理方法异常:" + msgHandleFuncName + ",Error:" + e.ToString());
                        }
                    }
                }
                else
                {
                    Console.WriteLine("收到客户端协议:" + msgid + "出错,未实现的处理方法:" + msgHandleFuncName);
                }
            }
            else
            {
                Console.WriteLine("收到客户端协议:" + msgid + "出错,ErrorCode:" + state);
            }
        }

    }
}
