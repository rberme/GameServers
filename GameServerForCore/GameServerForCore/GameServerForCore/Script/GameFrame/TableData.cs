using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using YunLoong;

namespace GameServer.CsScript
{
    using System;
    using System.Collections;
    using System.Collections.Generic;
    using System.Reflection;

    public class TableData
    {
        public static TableData Instance;
        public static int[] qualIconID_arr = new int[] { 0, 1, 2, 4, 7, 11 };
        public static int[] fqualIconID_arr = new int[] { 0, 1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5 };
        static TableData()
        {
            Instance = new TableData();
        }

        public Dictionary<string, Hashtable> tablehash;

        public Dictionary<int, ItemTable> ItemTablehash;
        public Dictionary<int, HeroTable> HeroTablehash;
        public Dictionary<int, EquipStarTable> EquipStarTablehash;
        public Dictionary<int, EquipTypeTable> EquipTypeTablehash;

        public Dictionary<int, BuffInfoTable> BuffInfoTablehash;
        public Dictionary<int, BattleBuffTable> BattleBuffTablehash;
        public Dictionary<int, BattleElementTable> BattleElementTablehash;
        public Dictionary<int, HeroSkillTable> HeroSkillTablehash;

        public Dictionary<int, SkillDisplayModeTable> SkillDisplayModeTablehash;
        public Dictionary<int, SkillLVBuffLVTable> SkillLVBuffLVTablehash;
        public Dictionary<int, SkillParamTable> SkillParamTablehash;

        public Dictionary<int, RockInfoTable> RockInfoTablehash;
        public Dictionary<int, SkillRockTable> SkillRockTablehash;

        public TableData()
        {

        }
        public void initTableData()
        {
            tablehash = new Dictionary<string, Hashtable>();

            cteatetablehashforxml<ItemTable>("item", "id", "id", ref ItemTablehash, ItemTable.fieldarr);
            cteatetablehashforxml<HeroTable>("hero", "id", "id", ref HeroTablehash, HeroTable.fieldarr);
            cteatetablehashforxml<EquipStarTable>("equipstar", "id", "id", ref EquipStarTablehash, EquipStarTable.fieldarr);
            cteatetablehashforxml<EquipTypeTable>("equiptype", "id", "id", ref EquipTypeTablehash, EquipTypeTable.fieldarr);

            cteatetablehashforxml<BuffInfoTable>("buffinfo", "id", "id", ref BuffInfoTablehash, BuffInfoTable.fieldarr);
            cteatetablehashforxml<BattleBuffTable>("battlebuff", "id", "id", ref BattleBuffTablehash, BattleBuffTable.fieldarr);
            cteatetablehashforxml<BattleElementTable>("battleelement", "id", "id", ref BattleElementTablehash, BattleElementTable.fieldarr);
            cteatetablehashforxml<HeroSkillTable>("heroskill", "id", "id", ref HeroSkillTablehash, HeroSkillTable.fieldarr);

            cteatetablehashforxml<SkillDisplayModeTable>("skilldisplaymode", "skillmodeid", "id", ref SkillDisplayModeTablehash, SkillDisplayModeTable.fieldarr);
            cteatetablehashforxml<SkillLVBuffLVTable>("skilllvbufflv", "id", "id", ref SkillLVBuffLVTablehash, SkillLVBuffLVTable.fieldarr);
            cteatetablehashforxml<SkillParamTable>("skillparam", "skillid", "id", ref SkillParamTablehash, SkillParamTable.fieldarr);


            cteatetablehashforxml<RockInfoTable>("rockinfo", "rocktype", "id", ref RockInfoTablehash, RockInfoTable.fieldarr);
            cteatetablehashforxml<SkillRockTable>("skillrock", "id", "id", ref SkillRockTablehash, SkillRockTable.fieldarr);

            //givehashfortablename("skillmode", "id");

        }



        public Hashtable givehashfortablename(string tablename, string tablekey)
        {
            Hashtable xmlhash;
            if (tablehash.ContainsKey(tablename))
            {
                xmlhash = tablehash[tablename];
            }
            else
            {
                //UnityEngine.Debug.Log(tablename);
                xmlhash = XMLHandle.givexmlforkeyhash(GameResources.Instance.xmlhash[tablename], tablekey);
                tablehash.Add(tablename, xmlhash);
            }
            return xmlhash;
        }
        public void cteatetablehashforxml<T>(string tablename, string tablekey, string keyfield, ref Dictionary<int, T> outhash, string[] fieldarr) where T : new()
        {
            Hashtable xmlhash = givehashfortablename(tablename, tablekey);

            outhash = new Dictionary<int, T>();
            ICollection keys = xmlhash.Keys;

            foreach (var key in keys)
            {
                int tid = Convert.ToInt32(key);
                T tcobj = new T();
                GameData.valuetofield(tcobj, keyfield, tid);
                outhash.Add(tid, tcobj);
                Hashtable shash = (Hashtable)xmlhash[key];
                for (int i = 0; i < fieldarr.Length; i++)
                {
                    if (shash.ContainsKey(fieldarr[i]))
                    {
                        GameData.valuetofield(tcobj, fieldarr[i], shash[fieldarr[i]]);
                    }
                }
                checkAndCreateParamDic<T, int>(tcobj, shash, "int");
                checkAndCreateParamDic<T, string>(tcobj, shash, "string");
                checkAndCreateParamDic<T, bool>(tcobj, shash, "bool");

                if (tcobj.GetType().GetMethod("initData") != null)
                {
                    tcobj.GetType().GetMethod("initData").Invoke(tcobj, null);
                }
            }
        }
        public void checkAndCreateParamDic<T, U>(T tcobj, Hashtable shash, string stype)
        {
            FieldInfo propNameField = typeof(T).GetField(stype + "fieldarr");
            if (propNameField != null)
            {
                string[] fieldName_arr = propNameField.GetValue(null) as string[];
                Dictionary<string, U> paramDic = typeof(T).GetField(stype + "ParamDic").GetValue(tcobj) as Dictionary<string, U>;
                for (int i = 0; i < fieldName_arr.Length; i++)
                {
                    if (shash.ContainsKey(fieldName_arr[i]))
                    {
                        paramDic.Add(fieldName_arr[i], (U)Convert.ChangeType(shash[fieldName_arr[i]], typeof(U)));
                        //UnityEngine.Debug.Log(typeof(U).ToString());
                    }
                }
            }

        }
        public void cteatetablehashforxml<T>(string tablename, string tablekey, string keyfield, ref Dictionary<string, T> outhash, string[] fieldarr) where T : new()
        {
            Hashtable xmlhash = givehashfortablename(tablename, tablekey);

            outhash = new Dictionary<string, T>();
            ICollection keys = xmlhash.Keys;
            foreach (var key in keys)
            {
                T tcobj = new T();
                string kid = null;
                Hashtable shash = (Hashtable)xmlhash[key];
                for (int i = 0; i < fieldarr.Length; i++)
                {
                    if (shash.ContainsKey(fieldarr[i]))
                    {
                        GameData.valuetofield(tcobj, fieldarr[i], shash[fieldarr[i]]);
                        if (fieldarr[i] == keyfield)
                        {
                            kid = Convert.ToString(shash[fieldarr[i]]);
                        }
                    }
                }
                if (!string.IsNullOrEmpty(kid))
                {
                    outhash.Add(kid, tcobj);
                }
            }
        }

        public object givePropFromTable(string tablename, int key, string fieldname)
        {
            if (tablehash[tablename].ContainsKey(Convert.ToString(key)))
            {
                Hashtable datahash = (Hashtable)(tablehash[tablename][Convert.ToString(key)]);
                if (datahash.ContainsKey(fieldname))
                {
                    return datahash[fieldname];
                }
            }
            return null;
        }

        public T givePropFromDisplayMode<T>(int modelid, string fieldname)
        {
            object tobj = givePropFromTable("displaymode", modelid, fieldname);
            if (tobj != null)
            {
                try
                {
                    return (T)Convert.ChangeType(tobj, typeof(T));
                }
                catch
                {
                    return default(T);
                }
            }
            else
            {
                return default(T);
            }
        }
        public List<int> giveDisplayModeForList(int skillid, string fieldname)
        {
            List<int> tlist = new List<int>();
            string tstr = givePropFromDisplayMode<string>(skillid, fieldname);
            if (tstr != null && tstr != string.Empty)
            {
                string[] tstr_arr = tstr.Split(',');
                for (int i = 0; i < tstr_arr.Length; i++)
                {
                    tlist.Add(Convert.ToInt32(tstr_arr[i]));
                }
            }
            return tlist;
        }
        /// <summary>
        /// 从skillprop表获得指定skillid项的指定字段fieldname的指定T类型值;
        /// </summary>
        /// <param name="skillid">指定技能ID</param>
        /// <param name="fieldname">指定获取的字段</param>
        /// <returns></returns>
        public T givePropFromSkillProp<T>(int skillid, string fieldname)
        {
            object tobj = givePropFromTable("skillmode", skillid, fieldname);
            if (tobj != null)
            {
                try
                {
                    return (T)Convert.ChangeType(tobj, typeof(T));
                }
                catch
                {
                    return default(T);
                }
            }
            else
            {
                return default(T);
            }
        }

        public List<int> givePropForList(int skillid, string fieldname)
        {
            List<int> tlist = new List<int>();
            string tstr = givePropFromSkillProp<string>(skillid, fieldname);
            if (tstr != null && tstr != string.Empty)
            {
                string[] tstr_arr = tstr.Split(',');
                for (int i = 0; i < tstr_arr.Length; i++)
                {
                    tlist.Add(Convert.ToInt32(tstr_arr[i]));
                }
            }
            return tlist;
        }

        public void setPropToObject(string tablename, int key, string[] fieldarr, object tobj)
        {
            Hashtable datahash = (Hashtable)(tablehash[tablename][Convert.ToString(key)]);
            for (int i = 0; i < fieldarr.Length; i++)
            {
                if (datahash.ContainsKey(fieldarr[i]))
                {
                    GameData.valuetofield(tobj, fieldarr[i], datahash[fieldarr[i]]);
                }
            }
        }
    }









    //数据表结构
    public class HeroTable
    {
        public static string[] fieldarr = new string[] { "name", "canuse", "coreid", "modelid", "initquality" };
        public static string[] intfieldarr = new string[] { "hpbase", "atkbase", "defbase", "hpdelta", "atkdelta", "defdelta", "spd", "crt", "crd", "bfh", "bfa", "att", "hit", "avd" };
        public Dictionary<string, int> intParamDic = new Dictionary<string, int>();

        public int id;
        public string name;

        public bool canuse;
        public int coreid;
        public int modelid;
        public int initquality;

    }

    public class EquipStarTable
    {
        public static string[] fieldarr = new string[] { "equipstar", "mainattr", "viceattr", "awakenattr" };
        public int id;
        public int equipstar;
        public Dictionary<string, int[]> mainattrDic;
        public string mainattr
        {
            set
            {
                mainattrDic = getPropArrDic(value);
            }
        }
        public Dictionary<string, int[]> viceattrDic;
        public string viceattr
        {
            set
            {
                viceattrDic = getPropArrDic(value);

            }
        }
        public Dictionary<string, int[]> awakenattrDic;
        public string awakenattr
        {
            set
            {
                awakenattrDic = getPropArrDic(value);

            }
        }

        private Dictionary<string, int[]> getPropArrDic(string tstr)
        {
            Dictionary<string, int[]> tempDic = new Dictionary<string, int[]>();
            string[] tstr_arr = tstr.Split('|');
            for (int i = 0; i < tstr_arr.Length; i++)
            {
                string[] tnstr_arr = tstr_arr[i].Split(',');

                int[] num_arr = new int[] { Convert.ToInt32(tnstr_arr[1]), Convert.ToInt32(tnstr_arr[2]) };
                string[] tkstr_arr = tnstr_arr[0].Split('%');
                for (int k = 0; k < tkstr_arr.Length; k++)
                {
                    tempDic.Add(tkstr_arr[k], num_arr);
                }

            }
            return tempDic;
        }
    }
    public class EquipTypeTable
    {
        public static string[] fieldarr = new string[] { "equiptype", "mainattr", "viceattr", "awakenattr" };
        public int id;
        public int equiptype;
        public RandPropStruct[] mainattr_arr;
        public string mainattr
        {
            set
            {
                mainattr_arr = getRandArr(value);
            }
        }
        public string[] viceattr_arr;
        public string viceattr
        {
            set
            {
                viceattr_arr = value.Split('|');
            }
        }

        public RandPropStruct[] awakenattr_arr;
        public string awakenattr
        {
            set
            {
                awakenattr_arr = getRandArr(value);
            }
        }
        private RandPropStruct[] getRandArr(string tstr)
        {
            string[] tstr_arr = tstr.Split('|');
            RandPropStruct[] temp_arr = new RandPropStruct[tstr_arr.Length];
            for (int i = 0; i < tstr_arr.Length; i++)
            {
                string[] tnstr_arr = tstr_arr[i].Split(',');
                for (int j = 0; j < tnstr_arr.Length; j++)
                {
                    int tRandNum = i > 0 ? Convert.ToInt32(tnstr_arr[1]) + temp_arr[i - 1].RandNum : Convert.ToInt32(tnstr_arr[1]);
                    temp_arr[i] = new RandPropStruct() { RandNum = tRandNum, PropName = tnstr_arr[0] };
                }
            }
            return temp_arr;
        }
    }
    public class RandPropStruct
    {
        public int RandNum;
        public string PropName;
    }

    public class ItemTable
    {
        public static string[] fieldarr = new string[] { "name", "type", "quality", "param", "info" };
        public int id;
        public string name;
        public int type;
        public int quality;
        public int param;
        public string info;
    }

    public class BuffInfoTable
    {
        public static Dictionary<string, int> buffNameToID = new Dictionary<string, int>();
        public static string[] fieldarr = new string[] { "type", "maxcount", "info", "name" };
        private int _id;
        public int id { get { return _id; } set { _id = value; buffkind = _id / 100; } }
        public int type;
        public int maxcount;
        public string info;
        public string name;

        public int buffkind;

        public void initData()
        {
            if (!string.IsNullOrEmpty(name))
            {
                buffNameToID.Add(name, id);
            }
        }
    }
    public class BattleBuffTable
    {
        public static string[] fieldarr = new string[] { "name", "round", "addcount" };
        private int _id;
        public int id { get { return _id; } set { _id = value; buffid = _id / 100; buffkind = _id / 10000; } }
        public int buffid;
        public int buffkind;
        //name需要拆解
        public string name
        {
            get
            {
                return mainKey;
            }
            set
            {
                string[] modestr_arr = value.Split('|');
                modeList = new List<BuffModeTable>();
                for (int i = 0; i < modestr_arr.Length; i++)
                {
                    string[] paramstr_arr = modestr_arr[i].Split(',');
                    modeList.Add(new BuffModeTable(paramstr_arr));
                }
            }
        }
        public int type { get { return TableData.Instance.BuffInfoTablehash[buffid].type; } }
        public int maxcount { get { return TableData.Instance.BuffInfoTablehash[buffid].maxcount; } }
        public string info { get { return TableData.Instance.BuffInfoTablehash[buffid].info; } }

        public int round;
        public int addcount;
        public string mainKey { get { return modeList[0].mainKey; } }
        public bool isProp { get { return modeList[0].isProp; } }
        public List<BuffModeTable> modeList;
    }
    public class BuffModeTable
    {
        public string mainKey;
        public bool isProp;
        public List<int> paramList;

        public BuffModeTable(string[] paramstr_arr)
        {
            if (paramstr_arr[0] == "Prop")
            {
                isProp = true;
                mainKey = paramstr_arr[1];
                paramList = new List<int>();
                paramList.Add(Convert.ToInt32(paramstr_arr[2]));
            }
            else
            {
                mainKey = paramstr_arr[0];
                if (paramstr_arr.Length >= 2)
                {
                    paramList = new List<int>();
                    for (int i = 1; i < paramstr_arr.Length; i++)
                    {
                        paramList.Add(Convert.ToInt32(paramstr_arr[i]));
                    }
                }
            }
        }
    }
    public class BattleElementTable
    {
        public static string[] fieldarr = new string[] { "type", "info", "skill", "formula" };
        public int id;
        public string info;
        public int type;
        public int skill;
        public string formula;//设置公式模式

    }
    public class HeroSkillTable
    {
        public static string[] fieldarr = new string[] { "skillid", "name", "costmp", "slotid", "rocktype", "costmp", "elementtype" };
        public int id;
        public string name;
        public int costmp;
        public int elementtype;
        public List<int> skillidList;
        public List<int> rocktypeList;
        public int[] slotid_arr;
        public string skillid
        {
            set
            {
                string[] param_arr = value.Split(',');
                skillidList = new List<int>();
                for (int i = 0; i < param_arr.Length; i++)
                {
                    skillidList.Add(Convert.ToInt32(param_arr[i]));
                }
            }
        }

        public string rocktype
        {
            set
            {
                string[] param_arr = value.Split(',');
                rocktypeList = new List<int>();
                for (int i = 0; i < param_arr.Length; i++)
                {
                    rocktypeList.Add(Convert.ToInt32(param_arr[i]));
                }
            }
        }
        public string slotid
        {
            set
            {
                string[] param_arr = value.Split(',');
                slotid_arr = new int[param_arr.Length];
                for (int i = 0; i < param_arr.Length; i++)
                {
                    slotid_arr[i] = Convert.ToInt32(param_arr[i]);
                }
            }
        }
    }
    public class SkillDisplayModeTable
    {
        public static string[] fieldarr = new string[] { "info" };
        public static string[] intfieldarr = new string[] { "passivemainmode", "passivechildmode" };
        public Dictionary<string, int> intParamDic = new Dictionary<string, int>();
        public int id;//"skillmodeid";
        public string info;
    }
    public class SkillLVBuffLVTable
    {
        public static string[] fieldarr = new string[] { "rate", "bufflv" };
        public int id;
        public int rate;
        public int bufflv;
    }
    public class RockInfoTable
    {
        public static string[] fieldarr = new string[] { "name", "elementtype" };
        public int id;
        public string name;
        public int elementtype;
    }
    public class SkillRockTable
    {
        public static string[] fieldarr = new string[] { "quality", "info", "skillid", "rocktype", "icon", "maxlv" };
        public int id;
        public int quality;
        public string info;
        public int skillid;
        public int rocktype;
        public string icon;
        public int maxlv;

        public string name { get { return TableData.Instance.RockInfoTablehash[rocktype].name; } }
    }

    public class SkillParamTable
    {
        public static string[] fieldarr = new string[] { "skillmodeid", "info" };
        public static string[] stringfieldarr = new string[] { "target", "defaulttarget", "conditionattr", "attr", "changetype", "triggertimes", "hpratecondition" };
        //public static string[] boolfieldarr = new string[] {  };
        public static string[] intfieldarr = new string[] { "buffid", "clearbufftype", "clearbuffkind", "skillratebase", "skillratedelta", "healskillratebase", "healskillratedelta", "ratebase", "ratedelta", "bufflvbase", "bufflvdelta", "clearcountbase", "clearcountdelta", "changempbase", "changempdelta", "spdtimeratebase", "spdtimeratedelta", "conditionbuffid", "attrvaluebase", "attrvaluedelta", "armorskillratebase", "armorskillratedelta", "dearmorskillratebase", "dearmorskillratedelta", "getaorb", "rate", "findTargetNum", "stateMode", "changemp", "rangeicon", "hpskillratebase", "hpskillratedelta", "hpskillratelimitbase", "hpskillratelimitdelta", "hpratebelow", "atkskillratebase", "atkskillratedelta", "hprateconditionvalue", "conditionateam", "conditionbteam", "triggermaxcount" };
        public Dictionary<string, string> stringParamDic = new Dictionary<string, string>();
        public Dictionary<string, bool> boolParamDic = new Dictionary<string, bool>();
        public Dictionary<string, int> intParamDic = new Dictionary<string, int>();

        public int id;
        public int skillmodeid;
        public string info;

    }
   
}
