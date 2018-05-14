using System;
using YunLoong;

namespace GameServer.CsScript.Model
{
    public class PlayerEquip
    {
        [EntityField(true, 0)]
        public long pid;
        [EntityField(true, 1)]
        public int sid;
        [EntityField]
        public int equipid;
        [EntityField]
        public int eexp;
        [EntityField]
        public int erand;
        [EntityField]
        public int eawakenrand;
        [EntityField]
        public int eawakencount;
        [EntityField]
        public DateTime ecreateticks;
        [EntityField]
        public int eowner;
        [EntityField]
        public bool eislock;
        [EntityField]
        public byte estate;
    }
}
