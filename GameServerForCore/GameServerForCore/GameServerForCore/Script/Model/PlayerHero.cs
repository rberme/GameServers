using YunLoong;

namespace GameServer.CsScript.Model
{
    public class PlayerHero : BaseEntity
    {
        [EntityField(true, 0)]
        public long pid;
        [EntityField(true, 1)]
        public int heroid;
        [EntityField]
        public int hexp;
        [EntityField]
        public int hquality;
        [EntityField]
        public int hawaken;
        [EntityField]
        public int hequip1;
        [EntityField]
        public int hequip2;
        [EntityField]
        public int hequip3;
        [EntityField]
        public int hequip4;
        [EntityField]
        public int hequip5;
        [EntityField]
        public int hequip6;
        [EntityField]
        public int hskillrock1;
        [EntityField]
        public int hskillrock2;
        [EntityField]
        public int hskillrock3;
        [EntityField]
        public int hskillrock4;
        [EntityField]
        public int hskillrock5;
        [EntityField]
        public int hskillrock6;
        [EntityField]
        public int hskillrock7;
        [EntityField]
        public int hskillrock8;
        [EntityField]
        public int hskin;

    }
}
