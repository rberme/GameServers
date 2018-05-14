using YunLoong;

namespace GameServer.CsScript.Model
{
    public class PlayerItem
    {
        [EntityField(true, 0)]
        public long pid;
        [EntityField(true, 1)]
        public int itemid;
        [EntityField]
        public int num;
    }
}
