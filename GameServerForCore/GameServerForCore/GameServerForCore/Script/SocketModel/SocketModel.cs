using System;
using ProtoBuf;
using System.Collections.Generic;

namespace GameServer.Model
{
    [ProtoContract]
    [ProtoInclude(101, typeof(MyPlayerData))]
    [ProtoInclude(102, typeof(GPlayerData))]
    public class BPlayerData
    {
        [ProtoMember(1)]
        public long pid;
        [ProtoMember(2)]
        public string pname;
        [ProtoMember(3)]
        public string picon;
        [ProtoMember(4)]
        public string piconborder;
        [ProtoMember(5)]
        public int pexp;
    }
    [ProtoContract]
    public class MyPlayerData : BPlayerData
    {
        [ProtoMember(11)]
        public int pdiamond;
        [ProtoMember(12)]
        public int pgold;
        [ProtoMember(13)]
        public int pspirit;
        [ProtoMember(14)]
        public int pmedal;
        [ProtoMember(15)]
        public int phonor;
        [ProtoMember(21)]
        public int pguildid;
        [ProtoMember(31)]
        public int equipBagCount;
    }

    [ProtoContract]
    public class GPlayerData : BPlayerData
    {
        [ProtoMember(15)]
        public int phonor;
        //公会ID GUILD
        [ProtoMember(21)]
        public int pguildid;
    }

    [ProtoContract]
    [ProtoInclude(101, typeof(QHeroData))]
    [ProtoInclude(102, typeof(PHeroData))]
    public class BHeroData
    {
        [ProtoMember(1)]
        public int heroid;//对应表中ID
        [ProtoMember(2)]
        public int hexp;//角色总经验
        [ProtoMember(3)]
        public int hquality;//品阶
        [ProtoMember(4)]
        public int hawaken;//觉醒
        [ProtoMember(5)]
        public int hskin;//皮肤ID
        [ProtoMember(6)]
        public Dictionary<int, int> hRockSlot;
        /*
        hRockSlot1
        hRockSlot2
        ......
        hRockSlot8

        qEquipSlot1
        qEquipSlot2
        .......
        qEquipSlot6
        */
    }

    [ProtoContract]
    public class QHeroData : BHeroData
    {
        [ProtoMember(11)]
        public Dictionary<int, int> qEquipSlot;//equip的引用列表:sid
    }
    [ProtoContract]
    public class PHeroData : BHeroData
    {
        [ProtoMember(11)]
        public Dictionary<int, BEquipData> bEquipSlot;//equip的对象列表
        [ProtoMember(21)]
        public long pid;//英雄所属玩家id
    }


    [ProtoContract]
    [ProtoInclude(101, typeof(GEquipData))]
    public class BEquipData
    {
        [ProtoMember(1)]
        public int equipid;//对应表中ID//propMode+equiptype+star;//equipid = (((propMode << 4) + equiptype) << 4) + star;
        [ProtoMember(2)]
        public int eexp;//装备总经验
        [ProtoMember(3)]
        public int erand;//属性随机种子参数
        [ProtoMember(4)]
        public int eawakenrand;//觉醒随机种子参数
        [ProtoMember(5)]
        public int eawakencount;//觉醒计数

        public int propMode
        {
            get
            {
                return equipid >> 8;
            }
        }
        public int equiptype//装备类型
        {
            get
            {
                return (equipid >> 4) & 0xf;
            }
        }
        public int star//装备星级
        {
            get
            {
                return equipid & 0xf;
            }
        }


        public int equipTypeKey
        {
            get
            {
                return propMode * 100 + equiptype;
            }
        }
        public int equipStarKey
        {
            get
            {
                return propMode * 100 + star;
            }
        }

        public int equipIcon { get { return star != 2 ? equiptype * 10 + star : equiptype * 10 + 1; } }
    }
    [ProtoContract]
    public class GEquipData : BEquipData
    {
        [ProtoMember(11)]
        public int sid;//唯一ID
        [ProtoMember(12)]
        public int eowner;//拥有者（英雄ID+槽ID），default：0未装裱
        [ProtoMember(13)]
        public bool eislock;//锁定
        [ProtoMember(14)]
        public int ecreateticks;//创建时的时间戳

    }
    [ProtoContract]
    public class GRockData
    {
        [ProtoMember(1)]
        public int rockid;//rockid+element+lv;//rid = (((id << 4) + element) << 4) + lv;id六位数,element和lv最大15
        [ProtoMember(2)]
        public int num;
        public int Id
        {
            get
            {
                return rockid >> 8;
            }
        }
        public int Element
        {
            get
            {
                return (rockid >> 4) & 0xf;
            }
        }
        public int Lv
        {
            get
            {
                return rockid & 0xf;
            }
        }

    }
    [ProtoContract]
    public class GItemData
    {
        [ProtoMember(1)]
        public int itemid;
        [ProtoMember(2)]
        public int num;
    }
    //单一布阵对象（用于PVP对方布阵）
    [ProtoContract]
    public class PLineupData
    {
        [ProtoMember(1)]
        public long pid;
        [ProtoMember(2)]
        public Dictionary<int, PHeroData> lineupDic;
    }


    //单一布阵对象（用于自身布阵）
    [ProtoContract]
    public class CLineupData
    {
        [ProtoMember(1)]
        public int sid;
        [ProtoMember(2)]
        public Dictionary<int, CHeroData> lineupDic;
    }
    //引用英雄对象（用于布阵，记录传递己方布阵信息）
    [ProtoContract]
    public class CHeroData
    {
        [ProtoMember(1)]
        public long pid;
        [ProtoMember(2)]
        public int heroid;
    }

}