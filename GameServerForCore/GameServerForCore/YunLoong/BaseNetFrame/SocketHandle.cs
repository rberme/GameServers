using System;
using System.Collections.Concurrent;
using System.IO;
using System.Threading;

namespace YunLoong
{

    public static class BinaryReadWriteExtend
    {
        public static int[] ReadIntArr(this BinaryReader reader)
        {
            //int len = reader.ReadInt32();
            //byte[] bytes = reader.ReadBytes(len);
            //return ProtoBufUtils.Deserialize<int[]>(bytes);
            return reader.ReadObject<int[]>();
        }

        public static T ReadObject<T>(this BinaryReader reader)
        {
            int len = reader.ReadInt32();
            byte[] bytes = reader.ReadBytes(len);
            return ProtoBufUtils.Deserialize<T>(bytes);
        }


        public static void Write<T>(this BinaryWriter bw, T obj)
        {
            byte[] temp = ProtoBufUtils.Serialize(obj);
            bw.Write(temp.Length);
            bw.Write(temp);
        }
    }


    public abstract class SocketHandle
    {
        public GSocket GameSocket;
        public Gbytes gBytes
        {
            get
            {
                return GameSocket.gBytes;
            }
        }

        public abstract void ReadMsgHandle(BinaryReader br);
        public abstract void SendCallback(int code);
        public abstract void CloseSocket();

        //socket客户端连接成功时回调
        public virtual void ConnectedHandle() { }
        //socket客户端连接失败时回调
        public virtual void ConnectFail() { }
        //同步发送
        public void SendByteMsgNow(int code,byte state)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(13);
            bw.Write(code);
            bw.Write(state);
            bw.Write(0);
            GameSocket.sendbytesNow(gBytes.wbytes);
        }
        public void SendByteMsgNow(int code,byte state, byte[] sendbuffer)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(13 + sendbuffer.Length);
            bw.Write(code);
            bw.Write(state);
            bw.Write(sendbuffer.Length);
            bw.Write(sendbuffer);
            GameSocket.sendbytesNow(gBytes.wbytes);
        }
        //异步发送
        public void SendByteMsg(int code, byte state)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(13);
            bw.Write(code);
            bw.Write(state);
            bw.Write(0);
            GameSocket.Sendbytes(code, gBytes.wbytes);
        }
        public void SendByteMsg(int code, byte state, byte[] sendbuffer)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(13 + sendbuffer.Length);
            bw.Write(code);
            bw.Write(state);
            bw.Write(sendbuffer.Length);
            bw.Write(sendbuffer);
            GameSocket.Sendbytes(code, gBytes.wbytes);
        }

        public void SendRequest(int msgcode, byte state)
        {
            SendByteMsg(msgcode, state);
        }

        public void SendRequest(int code, params object[] tparams)
        {
            BinaryWriter bw = gBytes.writer();
            for (int i = 0; i < tparams.Length; i++)
            {
                object param = tparams[i];
                if (param is int)
                    bw.Write((int)param);
                else if (param is string)
                    bw.Write((string)param);
                else if (param is bool)
                    bw.Write((bool)param);
                else if (param is long)
                    bw.Write((long)param);
                else if (param is byte)
                    bw.Write((byte)param);
                else
                {
                    byte[] bytes = ProtoBufUtils.Serialize(param);
                    bw.Write(bytes.Length);
                    bw.Write(bytes);
                }
            }
            SendByteMsg(code, 0, gBytes.wbytes);
        }

        /*

        public void SendRequest(int msgcode, int param1)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(param1);
            SendByteMsg(msgcode, gBytes.wbytes);
        }


        public void SendRequest(int msgcode, int param1, int param2)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(param1);
            bw.Write(param2);
            SendByteMsg(msgcode, gBytes.wbytes);
        }

        public void SendRequest(int msgcode, int param1, byte[] sendbuffer)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(param1);
            bw.Write(sendbuffer);
            SendByteMsg(msgcode, gBytes.wbytes);
        }

        public void SendRequest(int msgcode, int param1, int param2, byte[] sendbuffer)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(param1);
            bw.Write(param2);
            bw.Write(sendbuffer);
            SendByteMsg(msgcode, gBytes.wbytes);
        }
        public void SendRequest(int msgcode, int param1, int param2, int param3, byte[] sendbuffer)
        {
            BinaryWriter bw = gBytes.writer();
            bw.Write(param1);
            bw.Write(param2);
            bw.Write(param3);
            bw.Write(sendbuffer);
            SendByteMsg(msgcode, gBytes.wbytes);
        }

        public void SendRequest(int msgcode, params object[] @params)
        {
            BinaryWriter bw = gBytes.writer();
            for (int i = 0; i < @params.Length; i++)
            {
                object param = @params[i];
                if (param is int)
                    bw.Write((int)param);
                else if (param is string)
                    bw.Write((string)param);
                else if (param is bool)
                    bw.Write((bool)param);
            }
            SendByteMsg(msgcode, gBytes.wbytes);
        }
        */

        public void ResponseObject(int code, Object obj)
        {
            byte[] bytes = ProtoBufUtils.Serialize(obj);
            int len = bytes.Length;
            SendRequest(code, len, bytes);
        }
        public void ResponseObject(int code, int param1, Object obj)
        {
            byte[] bytes = ProtoBufUtils.Serialize(obj);
            int len = bytes.Length;
            SendRequest(code, param1, len, bytes);
        }
        public void ResponseObject(int code, int param1, int param2, Object obj)
        {
            byte[] bytes = ProtoBufUtils.Serialize(obj);
            int len = bytes.Length;
            SendRequest(code, param1, param2, len, bytes);
        }

        //public abstract void changeUserTime();

        public abstract void SendHeartBeat();

        /*
        public override void SendCallback(int code)
        {
            if (code == MsgCode.MSG_CODE_HEARTBEAT_RET)
            {
                changeUserTime();//修改指定player的最后通信时间戳（下线时间戳），id
            }
        }
        */

    }
}
