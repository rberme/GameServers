using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Threading;
using YunLoong;
namespace GameServer.CsScript
{
    public class PlayerSocketHandle : ServerSocketHandle<MsgHandle>
    {
        public PlayerSocketHandle()
        {
            handleObject = new MsgHandle();
            handleObject.handle = this;
        }
        public void ActiveHandle(object id)
        {
            goodSocket = true;
            userid = id;
        }
        public override void changeUserTime()
        {
            //修改指定player的最后通信时间戳（下线时间戳），userid
            //Player myplayer = GameMain.Instance.givePlayer(userid);
            //myplayer.lastlogintime = DateTime.Now;
        }
        public void SendRequest(MsgCode msgcode, byte param1)
        {
            SendRequest((int)msgcode, param1);
        }
        public void SendRequest(MsgCode msgcode, params object[] tparams)
        {
            SendRequest((int)msgcode, tparams);
        }
        
        public override string getMsgHandleFuncName(int msgid)
        {
            return ((MsgCode)msgid).ToString();
        }
        public override void CloseSocket()
        {
            try
            {
                if (goodSocket)
                {
                    changeUserTime();//修改指定player的最后通信时间戳（下线时间戳），userid
                    PlayerSocketManager.Instance.CloseSession(userid, GameSocket);
                }
                else
                {
                    PlayerSocketManager.Instance.CloseSession(GameSocket);
                }
            }
            catch (Exception ex)
            {
                //Console.WriteLine("清理数据:" + userid);
            }
        }
    }
}