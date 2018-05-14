using GameServer.Model;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using YunLoong;

namespace GameServer.CsScript
{
    //指定协议号和回调方法名
    public enum MsgCode : int
    {
        MSG_CODE_LOGIN = 101,//登录协议请求
        MSG_CODE_LOGIN_RET = 102,//登录协议返回

        MSG_CODE_PLAYER = 1101,//玩家信息协议请求
        MSG_CODE_PLAYER_RET = 1102,//玩家信息协议返回

        GM_CODE_ADDHERO = 98001,
        GM_CODE_ADDWEAPON = 98002
    }
    //协议处理类
    public class MsgHandle
    {
        public Gbytes gBytes;
        public PlayerSocketHandle handle;

        public MsgHandle()
        {
            gBytes = new Gbytes();
        }
        public void MSG_CODE_LOGIN(byte[] bytes)
        {
            BinaryReader br = gBytes.giveBinaryReader(bytes);
            Console.WriteLine("MSG_CODE_LOGIN:" + br.ReadString());


            //Console.WriteLine(GNetManager.SendRequestForPost("http://localhost/test/testgetpost.php", GNetManager.CreateRequestParam("param", 9655)));
            //初始化用户缓存数据
            PlayerCacheLock tPlayerCacheLock;
            PlayerCacheManager.instance.TryFindOrCreateKey(1300000100001, out tPlayerCacheLock);
            //激活sokcethandle并保存
            handle.ActiveHandle(1300000100001);
            PlayerSocketManager.Instance.SaveSession(1300000100001, handle.GameSocket);

            MyPlayerData tMyPlayerData = tPlayerCacheLock.giveMyPlayerData();
            handle.SendRequest(MsgCode.MSG_CODE_PLAYER_RET, tMyPlayerData);

            handle.SendRequest(MsgCode.MSG_CODE_LOGIN_RET, 0);



        }
        public void MSG_CODE_PLAYER(byte[] bytes)
        {

        }
    }
    //协议方法参数
    public class OperationCode
    {
        //操作号
        public const int OPER_CODE_ALL = 0;//返回所有数据（客户端刷新所有）- 应用（英雄，萌动元素，布阵）
        public const int OPER_CODE_PART = 1;//返回指定数据（客户端刷新指定）- 应用（英雄，萌动元素，布阵）

    }
}