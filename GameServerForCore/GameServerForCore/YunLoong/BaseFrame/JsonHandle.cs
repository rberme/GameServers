using System.Collections;
using System.Collections.Generic;
using System;
using System.Text;

namespace YunLoong
{
    public class JsonHandle
    {

        public JsonHandle()
        {

        }
        /*
        public static void parse<T>(ref T obj,string json){

        }
        */
        public static Hashtable parsejson(string json)
        {
            //Debug.Log(json);
            int parenttype;
            if (json.Substring(0, 1) == "{")
            {
                parenttype = 0;
            }
            else if (json.Substring(0, 1) == "[")
            {
                parenttype = 1;
            }
            else
            {
                return null;
            }
            json = json.Substring(1, json.Length - 2);


            List<int> leftdkh = new List<int>();
            List<int> rightdkh = new List<int>();
            List<int> leftzkh = new List<int>();
            List<int> rightzkh = new List<int>();
            List<int> dot = new List<int>();
            List<int> maohao = new List<int>();

            getjsondata(json, ref leftdkh, ref rightdkh, ref leftzkh, ref rightzkh, ref dot, ref maohao);

            Hashtable jsonhash = new Hashtable();
            if (json != string.Empty)
            {
                getHashtablebox(leftdkh, rightdkh, leftzkh, rightzkh, dot, maohao, json, 0, parenttype, 0, jsonhash);
            }
            return jsonhash;
        }
        private static void getjsondata(string json, ref List<int> leftdkh, ref List<int> rightdkh, ref List<int> leftzkh, ref List<int> rightzkh, ref List<int> dot, ref List<int> maohao)
        {
            char[] jsonchar = json.ToCharArray();
            List<int> yinhao = new List<int>();
            leftdkh.Clear();
            rightdkh.Clear();
            leftzkh.Clear();
            rightzkh.Clear();
            dot.Clear();
            maohao = new List<int>();
            for (int i = 0; i < jsonchar.Length; i++)
            {
                if (jsonchar[i] == '"')
                {
                    if (i > 0)
                    {
                        if (jsonchar[i - 1] != '\\')
                        {
                            yinhao.Add(i);
                        }
                    }
                    else
                    {
                        yinhao.Add(i);
                    }
                }
            }
            int yhcount = 0;

            int lastnum;
            if (yhcount < yinhao.Count)
            {
                lastnum = yinhao[yhcount];
            }
            else
            {
                lastnum = jsonchar.Length;
            }
            for (int i = 0; i < jsonchar.Length; i++)
            {
                if (i < lastnum)
                {
                    switch (jsonchar[i])
                    {
                        case '{':
                            leftdkh.Add(i);
                            break;
                        case '}':
                            rightdkh.Add(i);
                            break;
                        case '[':
                            leftzkh.Add(i);
                            break;
                        case ']':
                            rightzkh.Add(i);
                            break;
                        case ',':
                            dot.Add(i);
                            break;
                        case ':':
                            maohao.Add(i);
                            break;
                    }
                }
                else
                {
                    i = yinhao[yhcount + 1];
                    if (yhcount + 2 < yinhao.Count)
                    {
                        yhcount += 2;
                        lastnum = yinhao[yhcount];
                    }
                    else
                    {
                        lastnum = jsonchar.Length;
                    }
                }
            }
        }
        private static void getHashtablebox(List<int> leftdkh, List<int> rightdkh, List<int> leftzkh, List<int> rightzkh, List<int> dot, List<int> maohao, string json, int pyz, int parenttype, int depth, Hashtable parentclip, int id = 0)
        {

            //Debug.Log (json);
            int count = 0;
            int lcount = 1;
            int rcount = 0;
            List<int> leftkh = new List<int>();
            List<int> rightkh = new List<int>();
            int type;
            int dmin = json.Length + 100;
            int dkhmin = json.Length + 100;
            int zkhmin = json.Length + 100;
            if (leftdkh.Count > 0)
            {
                dkhmin = leftdkh[0];
            }
            if (leftzkh.Count > 0)
            {
                zkhmin = leftzkh[0];
            }
            if (dot.Count > 0)
            {
                dmin = dot[0];
            }
            int minnum = Math.Min(dkhmin, Math.Min(zkhmin, Math.Min(dmin, json.Length)));
            if (minnum == dkhmin)
            {
                leftkh = leftdkh;
                rightkh = rightdkh;
                type = 0;
            }
            else if (minnum == zkhmin)
            {
                leftkh = leftzkh;
                rightkh = rightzkh;
                type = 1;
            }
            else if (minnum == dmin)
            {
                type = 2;
                string nstr = json.Substring(pyz, dot[0] - pyz);
                //string lstr = json.Substring (dot [0] + 1);
                if (parenttype == 0)
                {
                    //对象元素
                    string tkey = nstr.Substring(0, maohao[0] - pyz);
                    string tvalue = nstr.Substring(maohao[0] - pyz + 1);
                    addhashvalue(parentclip, tkey.Trim('\"'), tvalue);

                    pyz = dot[0] + 1;
                    dot.RemoveAt(0);
                    maohao.RemoveAt(0);

                    getHashtablebox(leftdkh, rightdkh, leftzkh, rightzkh, dot, maohao, json, pyz, parenttype, depth, parentclip);
                }
                else if (parenttype == 1)
                {
                    //数组元素
                    addhashvalue(parentclip, id, nstr);
                    id++;
                    pyz = dot[0] + 1;
                    dot.RemoveAt(0);
                    getHashtablebox(leftdkh, rightdkh, leftzkh, rightzkh, dot, maohao, json, pyz, parenttype, depth, parentclip, id);
                }
            }
            else
            {
                type = 3;
                string nstr = json.Substring(pyz);
                if (parenttype == 0)
                {
                    //对象元素

                    string tkey = nstr.Substring(0, maohao[0] - pyz);
                    string tvalue = nstr.Substring(maohao[0] - pyz + 1);
                    addhashvalue(parentclip, tkey.Trim('\"'), tvalue);
                }
                else if (parenttype == 1)
                {
                    //数组元素
                    addhashvalue(parentclip, id, nstr);
                }

            }
            if (type < 2)
            {
                int nleft;
                if (lcount < leftkh.Count)
                {
                    nleft = leftkh[lcount];
                }
                else
                {
                    nleft = json.Length;
                }
                for (int i = 0; i <= leftkh.Count + rightkh.Count; i++)
                {
                    if (nleft < rightkh[rcount])
                    {
                        lcount++;
                        count++;
                        if (lcount < leftkh.Count)
                        {
                            nleft = leftkh[lcount];
                        }
                        else
                        {
                            nleft = json.Length;
                        }
                    }
                    else
                    {
                        if (count == 0)
                        {
                            break;
                        }
                        else
                        {
                            rcount++;
                            count--;
                        }
                    }
                }
                if (parenttype == 0)
                {
                    //对象元素
                    string nstr = json.Substring(pyz, rightkh[rcount] - pyz + 1);
                    string tkey = nstr.Substring(0, maohao[0] - pyz);
                    string tvalue = nstr.Substring(maohao[0] - pyz + 1);
                    Hashtable thash = parsejson(tvalue);
                    addhashvalue(parentclip, tkey.Trim('\"'), thash);
                    if (rightkh[rcount] + 1 < json.Length)
                    {
                        string lstr = json.Substring(rightkh[rcount] + 2);
                        getjsondata(lstr, ref leftdkh, ref rightdkh, ref leftzkh, ref rightzkh, ref dot, ref maohao);
                        getHashtablebox(leftdkh, rightdkh, leftzkh, rightzkh, dot, maohao, lstr, 0, parenttype, depth, parentclip);
                    }
                }
                else
                {
                    //数组元素
                    string nstr = json.Substring(pyz, rightkh[rcount] - pyz + 1);
                    Hashtable thash = parsejson(nstr);
                    addhashvalue(parentclip, id, thash);
                    id++;
                    if (rightkh[rcount] + 1 < json.Length)
                    {
                        string lstr = json.Substring(rightkh[rcount] + 2);
                        getjsondata(lstr, ref leftdkh, ref rightdkh, ref leftzkh, ref rightzkh, ref dot, ref maohao);
                        getHashtablebox(leftdkh, rightdkh, leftzkh, rightzkh, dot, maohao, lstr, 0, parenttype, depth, parentclip, id);
                    }
                }
            }
        }

        private static void addhashvalue(Hashtable thash, object tkey, object tvalue)
        {

            try
            {
                if (tvalue.ToString() != "System.Collections.Hashtable")
                {
                    if (((string)tvalue).IndexOf("\"") == 0)
                    {
                        tvalue = ((string)tvalue).Trim('\"');
                    }
                    /*else{
                        tvalue = System.Convert.ToInt32(tvalue);
                    }*/
                }
                if (thash.ContainsKey(tkey))
                {
                    thash[tkey] = tvalue;
                }
                else
                {
                    thash.Add(tkey, tvalue);
                }
            }
            catch
            {
                Console.WriteLine(tkey);
                Console.WriteLine(tvalue);
            }
        }
        public static Hashtable getstrlist(Hashtable tmptable, List<string> strlistid)
        {
            for (int i = 0; i < strlistid.Count; i++)
            {
                tmptable = (Hashtable)tmptable[strlistid[i]];
            }
            return tmptable;
        }
        public static void debughashtable(Hashtable strlist, string keybefore = "")
        {
            foreach (DictionaryEntry de in strlist)
            {
                if (de.Value.ToString() == "System.Collections.Hashtable")
                {
                    debughashtable((Hashtable)strlist[de.Key], keybefore + de.Key.ToString() + ".");
                }
                else
                {
                    Console.WriteLine(keybefore + de.Key.ToString() + "    " + de.Value.ToString());
                }
            }
        }
        public static string Hashtabletojson(Hashtable thash, int type = 0)
        {
            string str;
            if (type == 0)
            {
                str = "{";
                foreach (DictionaryEntry de in thash)
                {
                    str += "\"" + de.Key.ToString() + "\":";
                    if (de.Value.ToString() == "System.Collections.Hashtable")
                    {
                        if (((Hashtable)de.Value).ContainsKey(0))
                        {
                            str += Hashtabletojson((Hashtable)de.Value, 1) + ",";
                        }
                        else
                        {
                            str += Hashtabletojson((Hashtable)de.Value) + ",";
                        }
                    }
                    else
                    {
                        if (typeof(int) == de.Value.GetType())
                        {
                            str += de.Value.ToString() + ",";
                        }
                        else
                        {
                            str += "\"" + de.Value.ToString() + "\",";
                        }
                    }
                }
                str = str.TrimEnd(',');
                str += "}";
            }
            else
            {
                str = "[";
                for (int i = 0; i < thash.Count; i++)
                {
                    //Debug.Log (System.Convert.ToString(i));
                    //Debug.Log (((Hashtable)thash)[i]);
                    if ((((Hashtable)thash)[i]).ToString() == "System.Collections.Hashtable")
                    {
                        if (((Hashtable)((Hashtable)thash)[i]).ContainsKey(0))
                        {
                            str += Hashtabletojson((Hashtable)((Hashtable)thash)[i], 1) + ",";
                        }
                        else
                        {
                            str += Hashtabletojson((Hashtable)((Hashtable)thash)[i]) + ",";
                        }
                    }
                    else
                    {
                        if (typeof(int) == ((Hashtable)thash)[i].GetType())
                        {
                            str += ((Hashtable)thash)[i] + ",";
                        }
                        else
                        {
                            str += "\"" + ((Hashtable)thash)[i] + "\",";
                        }

                    }
                }
                str = str.TrimEnd(',');
                str += "]";
            }
            return str;
        }

        public static string ClearExtraJson(string tjsonstr)
        {
            tjsonstr = tjsonstr.Replace("\t", "");
            tjsonstr = tjsonstr.Replace("\r\n", "");
            tjsonstr = tjsonstr.Replace(" ", "");
            return tjsonstr;
        }

        public static string dictoKVstring<T, U>(Dictionary<T, U> dic)
        {
            string kvstr = string.Empty;
            var keys = dic.Keys;
            int tmax = keys.Count;
            int tcount = 0;
            kvstr = "[";
            foreach (var tkey in keys)
            {
                string tstr = "{" + tkey + "," + dic[tkey] + "}";
                kvstr += tstr;
                tcount++;
                if (tcount < tmax)
                {
                    kvstr += ",";
                }
            }
            kvstr += "]";
            return kvstr;
        }

        public static string listToArrayString<T>(List<T> list)
        {
            string listStr = string.Empty;
            int tmax = list.Count;
            int tcount = 0;
            listStr = "[";
            foreach (var ele in list)
            {
                string tstr = ele.ToString();
                listStr += tstr;
                tcount++;
                if (tcount < tmax)
                {
                    listStr += ",";
                }
            }
            listStr += "]";
            return listStr;
        }
    }
}