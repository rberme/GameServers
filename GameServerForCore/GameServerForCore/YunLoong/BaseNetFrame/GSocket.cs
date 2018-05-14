using System.Net.Sockets;
using System.Text;
using System;
using System.IO;
using System.Net;

namespace YunLoong
{

    public class GSocket : IDisposable
    {
        public SocketHandle socketHandle;
        public Socket wSocket;
        public Gbytes gBytes;
        public long id;

        public byte[] netmsg;
        public int netmsglen;

        public GSocket()
        {
            gBytes = new Gbytes();
        }
        public GSocket(Socket socket) : this()
        {
            wSocket = socket;
        }

        public void Connect(string ip, int port)
        {
            try
            {
                wSocket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                wSocket.Connect(IPAddress.Parse(ip), port);
                ReceiveAsync();
                socketHandle.ConnectedHandle();
            }
            catch
            {
                //Console.WriteLine("主动连接服务器失败");
                socketHandle.ConnectFail();
            }
        }

        public void ReceiveAsync()
        {
            gBytes.init(4);
            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadCallback), gBytes);
        }
        private void ReadCallback(IAsyncResult ar)
        {

            if (wSocket.Connected)
            {
                try
                {
                    int bytesRead = wSocket.EndReceive(ar);
                    if (bytesRead > 0)
                    {
                        if (bytesRead == 4)
                        {

                            BinaryReader br = gBytes.reader();
                            int msglen = br.ReadInt32();
                            if (msglen < 65536)
                            {
                                gBytes.init(msglen - 4);
                                wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadMsgCallback), gBytes);
                            }
                            else
                            {
                                gBytes.init(1024);
                                wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadExceptionCallback), gBytes);
                            }
                        }
                        else
                        {
                            gBytes.init(1024);
                            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadExceptionCallback), gBytes);
                        }
                    }
                    else
                    {
                        Console.WriteLine(id + "连接断开");
                        closeSocket();
                    }
                }
                catch
                {
                    Console.WriteLine(id + "连接服务器异常");
                }
            }
        }
        private void ReadMsgCallback(IAsyncResult ar)
        {
            if (wSocket.Connected)
            {
                try
                {
                    int bytesRead = wSocket.EndReceive(ar);
                    if (bytesRead > 0)
                    {
                        if (bytesRead == gBytes.bufferSize)
                        {
                            BinaryReader br = gBytes.reader();
                            socketHandle.ReadMsgHandle(br);
                            gBytes.init(4);
                            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadCallback), gBytes);
                        }
                        else
                        {
                            gBytes.init(1024);
                            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadExceptionCallback), gBytes);
                        }
                    }
                    else
                    {
                        Console.WriteLine(id + "连接断开");
                        closeSocket();
                    }
                }
                catch
                {
                    Console.WriteLine(id + "连接服务器异常");
                }
            }
        }
        private void ReadExceptionCallback(IAsyncResult ar)
        {
            if (wSocket.Connected)
            {
                try
                {
                    int bytesRead = wSocket.EndReceive(ar);
                    if (bytesRead > 0)
                    {

                        if (bytesRead < gBytes.bufferSize)
                        {
                            gBytes.init(4);
                            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadCallback), gBytes);
                        }
                        else
                        {
                            wSocket.BeginReceive(gBytes.buffer, 0, gBytes.bufferSize, 0, new AsyncCallback(ReadExceptionCallback), gBytes);
                        }
                    }
                    else
                    {
                        Console.WriteLine(id + "连接断开");
                        closeSocket();
                    }
                }
                catch
                {
                    Console.WriteLine(id + "连接服务器异常");
                }
            }
        }

        private void closeSocket()
        {
            socketHandle.CloseSocket();
        }
        
        public void Dispose()
        {
            wSocket.Close();
            wSocket.Dispose();
        }

        public void destroy()
        {
            wSocket = null;
        }
        //同步发送
        public void sendbytesNow(byte[] wbytes)
        {
            if (wSocket.Connected)
            {
                wSocket.Send(wbytes);
            }
        }
        //异步发送
        public void Sendbytes(int code, byte[] wbytes)
        {
            if (wSocket.Connected)
            {
                wSocket.BeginSend(wbytes, 0, wbytes.Length, 0, new AsyncCallback(SendCallback), code);
            }
        }
        private void SendCallback(IAsyncResult ar)
        {
            try
            {
                int bytesSent = wSocket.EndSend(ar);
                //Console.WriteLine("协议发送成功："+id+";code:"+ ar.AsyncState+"isHeartBeat"+ (Convert.ToInt32(ar.AsyncState) == MsgCode.MSG_CODE_HEARTBEAT_RET));
                socketHandle.SendCallback(Convert.ToInt32(ar.AsyncState));

            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
                Console.WriteLine("协议发送失败：" + id);
            }
        }
    }

}