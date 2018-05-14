using YunLoong;

namespace GameServer.CsScript.Model
{
    public class PlayerRock
    {
        [EntityField(true, 0)]
        public long pid;
        [EntityField(true, 1)]
        public int rockid;
        [EntityField]
        public int num;
    }
}
