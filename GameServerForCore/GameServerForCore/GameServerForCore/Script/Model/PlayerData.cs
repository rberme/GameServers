using System;
using YunLoong;

namespace GameServer.CsScript.Model
{
    public class PlayerData : BaseEntity 
    {
        [EntityField(true)]
        public long pid;
        [EntityField]
        public int userid;
        [EntityField]
        public int serverid;
        [EntityField]
        public string pname;
        [EntityField]
        public string picon;
        [EntityField]
        public string piconborder;
        [EntityField]
        public int pexp;
        [EntityField]
        public int pdiamond;
        [EntityField]
        public int pgold;
        [EntityField]
        public int pspirit;
        [EntityField]
        public int pmedal;
        [EntityField]
        public int phonor;
        [EntityField]
        public DateTime pcreatetime;
        [EntityField]
        public DateTime plastlogintime;
        [EntityField]
        public int prandseed;
        [EntityField]
        public int pguildid;
        [EntityField]
        public DateTime plastspiritsddtime;
    }
}
