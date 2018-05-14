using System;
using System.Collections;
using System.Collections.Generic;
using System.Xml;
using System.IO;
using System.Text;
using System.Threading;
using System.Reflection;
using YunLoong;

namespace GameServer.CsScript
{
    public class GameResources
    {

        public static GameResources Instance;
        public Action loadoverhandle;
        public Dictionary<string, XmlDocument> xmlhash;
        public Dictionary<string, Hashtable> jsonhash;
        private string[] xmlname_arr = new string[] { "battlebuff", "battleelement", "buffinfo", "equipstar", "equiptype", "hero", "heroskill", "item", "rockinfo", "skilldisplaymode", "skilllvbufflv", "skillparam", "skillrock" };
        private string[] jsonname_arr = new string[] { "gameDataConfig" };
        private int nowloadstep;
        private int totalcount = 0;

        static GameResources()
        {
            Instance = new GameResources();
        }
        public void startLoadRes(Action loadover = null)
        {
            loadoverhandle = loadover;
            xmlhash = new Dictionary<string, XmlDocument>();
            jsonhash = new Dictionary<string, Hashtable>();
            nowloadstep = 0;
            loadandCreateData(xmlname_arr[nowloadstep], ".xml", "xml", "xml/");
        }

        public void loadandCreateData(string dataName, string dataExt, string dataType, string dataPath = "", int dataID = 0)
        {

            string path = Environment.CurrentDirectory + "/Model/" + dataPath + dataName + dataExt;//ConfigUtils.GetSetting("ScriptRelativePath")+"/Model/" + dataPath + dataName + dataExt;
            loaddata(path, dataLoadover, dataName, dataType, dataID);
        }

        public void loaddata(string filepath, Action<string, string, string, int> loadcallback, string resName, string restype, int resID)
        {
            FileInfo fi = new FileInfo(filepath);
            //Console.WriteLine(fi.DirectoryName);
            if (fi.Exists)
            {
                FileStream fs = new FileStream(filepath, FileMode.Open);
                StreamReader reader = new StreamReader(fs);
                string tdata = reader.ReadToEnd();
                if (loadcallback != null)
                {
                    loadcallback(tdata, resName, restype, resID);
                }
                fs.Dispose();
                reader.Dispose();
            }
        }
        public void dataLoadover(string tdata, string resName, string restype, int resID)
        {
            totalcount++;
            if (restype == "xml")
            {
                XmlDocument xmltable = new XmlDocument();
                xmltable.LoadXml(tdata);
                xmlhash.Add(resName, xmltable);

                if (nowloadstep < xmlname_arr.Length - 1)
                {
                    nowloadstep++;
                    loadandCreateData(xmlname_arr[nowloadstep], ".xml", "xml", "xml/");
                }
                else
                {
                    if (jsonname_arr.Length > 0)
                    {
                        nowloadstep = 0;
                        loadandCreateData(jsonname_arr[nowloadstep], ".json", "json", "json/");
                    }
                    else
                    {
                        Console.WriteLine("数据表加载完成");
                    }
                }

            }
            if (restype == "json")
            {
                Hashtable thash = JsonHandle.parsejson(JsonHandle.ClearExtraJson(tdata));
                jsonhash.Add(resName, thash);
                if (nowloadstep < jsonname_arr.Length - 1)
                {
                    nowloadstep++;
                    loadandCreateData(jsonname_arr[nowloadstep], ".json", "json", "json/");
                }
                else
                {
                    Console.WriteLine("数据表加载完成");
                    if (loadoverhandle != null)
                    {
                        loadoverhandle();
                    }
                }
            }
        }
    }
}