using System;
using System.Collections;
using System.Collections.Generic;
using System.Reflection;
using System.Security.Cryptography;

namespace YunLoong
{
    public class GameData
    {

        //++++++++++++++++功能性方法++++++++++++++++++++//
        //“反射机制”给指定对象的指定变量赋值（变量名为字符串）
        public static void valuetofield(object cobj, string field, object valueobj)
        {
            if (cobj.GetType().GetField(field) != null)
            {
                try
                {
                    switch (cobj.GetType().GetField(field).FieldType.ToString())
                    {
                        case "System.String":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToString(valueobj));
                            break;
                        case "System.Int16":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToInt16(valueobj));
                            break;
                        case "System.Int32":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToInt32(valueobj));
                            break;
                        case "System.Int64":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToInt64(valueobj));
                            break;
                        case "System.Single":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToSingle(valueobj));
                            break;
                        case "System.Boolean":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToBoolean(valueobj));
                            break;
                        case "System.DateTime":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToDateTime(valueobj));
                            break;
                        case "System.Byte":
                            cobj.GetType().GetField(field).SetValue(cobj, System.Convert.ToByte(valueobj));
                            break;
                        case "System.Byte[]":
                            cobj.GetType().GetField(field).SetValue(cobj, (System.Byte[])(valueobj));
                            break;
                        default:
                            cobj.GetType().GetField(field).SetValue(cobj, null);// valueobj);
                            Console.WriteLine("未实现的类型:" + cobj.GetType().GetField(field).FieldType.ToString());
                            break;
                    }
                }
                catch (System.Exception ex)
                {
                    Console.WriteLine(cobj + "," + field + "," + valueobj);
                    Console.WriteLine("变量赋值类型错误:" + cobj.GetType().GetField(field).FieldType.ToString() + "对象" + cobj.ToString() + "变量" + field + "值" + valueobj + "\n错误:" + ex.ToString());
                }
            }
            else if (cobj.GetType().GetProperty(field) != null)
            {
                try
                {
                    switch (cobj.GetType().GetProperty(field).PropertyType.ToString())
                    {
                        case "System.String":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToString(valueobj), null);
                            break;
                        case "System.Int16":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToInt16(valueobj), null);
                            break;
                        case "System.Int32":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToInt32(valueobj), null);
                            break;
                        case "System.Int64":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToInt64(valueobj), null);
                            break;
                        case "System.Single":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToSingle(valueobj), null);
                            break;
                        case "System.Boolean":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToBoolean(valueobj), null);
                            break;
                        case "System.DateTime":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToDateTime(valueobj));
                            break;
                        case "System.Byte":
                            cobj.GetType().GetProperty(field).SetValue(cobj, System.Convert.ToByte(valueobj));
                            break;
                        case "System.Byte[]":
                            cobj.GetType().GetProperty(field).SetValue(cobj, (System.Byte[])(valueobj));
                            break;
                        default:
                            cobj.GetType().GetProperty(field).SetValue(cobj, null, null);// valueobj, null);
                            Console.WriteLine("未实现的类型:" + cobj.GetType().GetProperty(field).PropertyType.ToString());
                            break;
                    }
                }
                catch (System.Exception ex)
                {
                    Console.WriteLine(cobj + "," + field + "," + valueobj);
                    Console.WriteLine("变量赋值类型错误:" + cobj.GetType().GetProperty(field).PropertyType.ToString() + "对象" + cobj.ToString() + "变量" + field + "值" + valueobj + "\n错误:" + ex.ToString());
                }
            }

        }
        //获得指定对象的指定变量的值
        public static object givevalue(object cobj, string field)
        {
            if (cobj.GetType().GetField(field) != null)
            {
                return cobj.GetType().GetField(field).GetValue(cobj);
            }
            else if (cobj.GetType().GetProperty(field) != null)
            {
                return cobj.GetType().GetProperty(field).GetValue(cobj, null);
            }
            else
            {
                return null;
            }

        }
        public static T CloneModel<T>(System.Object Model, T b = null) where T : class, new()
        {
            T tobj;
            if (b != null)
            {
                tobj = b;
            }
            else
            {
                tobj = new T();
            }
            MemberInfo[] member_arr = tobj.GetType().GetProperties();
            MemberInfo[] Fields_arr = tobj.GetType().GetFields();
            if (Fields_arr.Length > 0)
            {
                int memberlen = member_arr.Length;
                Array.Resize(ref member_arr, member_arr.Length + Fields_arr.Length);
                Fields_arr.CopyTo(member_arr, memberlen);
            }
            for (int i = 0; i < member_arr.Length; i++)
            {
                object value = GameData.givevalue(Model, member_arr[i].Name);
                if (value != null)
                {
                    GameData.valuetofield(tobj, member_arr[i].Name, value);
                }
            }
            return tobj;
        }
        private static int _servertimer = gettimer();
        public static int servertimer
        {
            get
            {
                return _servertimer;
            }
            set
            {
                _servertimer = value;
            }
        }

        public static int gettimer()
        {
            long ntime = (DateTime.Now.ToUniversalTime().Ticks - 621355968000000000) / 10000000;
            return Convert.ToInt32(ntime);
        }
        public static DateTime timerToDate(int ntime)
        {
            long tick = Convert.ToInt64(ntime) * 10000000 + 621355968000000000;
            DateTime tdate = new DateTime(tick);
            return tdate;
        }
        public static int dateToTimer(DateTime tdate)
        {
            long ntime = (tdate.Ticks - 621355968000000000) / 10000000;
            return Convert.ToInt32(ntime);
        }
        //返回时间字符串(参数为时间戳)
        public static string getTimeStr(int timestamp, int type)
        {
            string str = "";
            DateTime time = GetTime(Convert.ToString(timestamp));
            if (type == 0)
            {
                str = time.Year + "/" + time.Month + "/" + time.Day + "  " + time.Hour + ":" + time.Minute + ":" + time.Second;
            }
            else
            {
                str = time.Year + "/" + time.Month + "/" + time.Day + "  " + time.Hour + ":" + time.Minute;
            }
            return str;
        }

        /// 时间戳转为C#格式时间
        private static DateTime GetTime(string timeStamp)
        {
            DateTime dtStart = TimeZone.CurrentTimeZone.ToLocalTime(new DateTime(1970, 1, 1));
            long lTime = long.Parse(timeStamp + "0000000");
            TimeSpan toNow = new TimeSpan(lTime);
            return dtStart.Add(toNow);
        }

        /// DateTime时间格式转换为Unix时间戳格式
        private static int ConvertDateTimeInt(System.DateTime time)
        {
            System.DateTime startTime = TimeZone.CurrentTimeZone.ToLocalTime(new System.DateTime(1970, 1, 1));
            return (int)(time - startTime).TotalSeconds;
        }

        public static string getclock(int num)
        {
            int s = num % 60;
            num = num / 60;
            int m = num % 60;
            num = num / 60;
            int h = num;
            string str = clockneedzero(h) + ":" + clockneedzero(m) + ":" + clockneedzero(s);
            return str;
        }
        public static string getclockStr(int num)
        {
            int s = num % 60;
            num = num / 60;
            int m = num % 60;
            num = num / 60;
            int h = num;
            string str = clockneedzero(h) + "小时" + clockneedzero(m) + "分钟" + clockneedzero(s) + "秒";
            return str;
        }

        public static string clockneedzero(int num)
        {
            if (num < 10)
            {
                return "0" + num;
            }
            else
            {
                return Convert.ToString(num);
            }
        }

        public static string MD5Encrypt(string strText)
        {
            MD5 md5 = new MD5CryptoServiceProvider();
            byte[] bytes = md5.ComputeHash(System.Text.Encoding.UTF8.GetBytes(strText));
            //return System.Text.Encoding.Default.GetString(result);         
            md5.Clear();
            string ret = "";
            for (int i = 0; i < bytes.Length; i++)
            {
                ret += Convert.ToString(bytes[i], 16).PadLeft(2, '0');
            }
            return ret.PadLeft(32, '0');
        }


        public static int Rand(ref int randznum, int n = 100000, int minznum = 2001, int maxznum = 20000)
        {
            randznum += 2;
            while (randznum > minznum + maxznum)
            {
                randznum -= maxznum;
            }
            int rnum = zrand(randznum, n);
            rnum = Convert.ToInt32(Math.Floor(rnum / 100f));
            return rnum;
        }


        public static int zrand(int zn, int fw = 100000)
        {
            long znn = zn;
            if (zn % 5 != 0)
            {
                return (int)(znn * znn * znn % fw);
            }
            else
            {
                return (int)((znn * (znn + 2) * (znn + 4)) % fw);
            }
        }
    }
}