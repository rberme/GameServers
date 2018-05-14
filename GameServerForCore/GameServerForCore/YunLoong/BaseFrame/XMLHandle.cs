using System.Collections;
using System.Collections.Generic;
using System;
using System.Xml;
namespace YunLoong
{
    public class XMLHandle
    {
        public static XmlNodeList givexmldata(string xmlstr)
        {
            XmlDocument xmltable = new XmlDocument();
            xmltable.LoadXml(xmlstr);
            XmlElement xmlroot = xmltable.DocumentElement;
            XmlNodeList xmltablelist = xmlroot.GetElementsByTagName("info");
            return xmltablelist;
        }
        public static Dictionary<int, XmlElement> givexmlfortableID(string xmlstr)
        {
            XmlNodeList xmltablelist = givexmldata(xmlstr);
            Dictionary<int, XmlElement> hashtmp = new Dictionary<int, XmlElement>();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                XmlElement tmp = ((XmlElement)xmltablelist[i]);
                hashtmp.Add(Convert.ToInt32(tmp.GetAttribute("tableID")), tmp);
            }
            return hashtmp;
        }
        public static Dictionary<string, string> givexmlforkeyvalue(string xmlstr, string keyname, string valuename)
        {
            XmlNodeList xmltablelist = givexmldata(xmlstr);
            Dictionary<string, string> hashtmp = new Dictionary<string, string>();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                XmlElement tmp = ((XmlElement)xmltablelist[i]);
                hashtmp.Add(tmp.GetAttribute(keyname), tmp.GetAttribute(valuename));
            }
            return hashtmp;
        }
        public static Dictionary<string, string[]> givexmlforkeyvalue(string xmlstr, string keyname, string[] valuename)
        {
            XmlNodeList xmltablelist = givexmldata(xmlstr);
            Dictionary<string, string[]> hashtmp = new Dictionary<string, string[]>();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                XmlElement tmp = ((XmlElement)xmltablelist[i]);

                string[] strarr = new string[valuename.Length];
                for (int j = 0; j < valuename.Length; j++)
                {
                    strarr[j] = tmp.GetAttribute(valuename[j]);
                }
                hashtmp.Add(tmp.GetAttribute(keyname), strarr);
            }
            return hashtmp;
        }
        public static Hashtable givexmlforkeyhash(XmlDocument xmldoc, string keyname)
        {
            XmlElement xmlroot = xmldoc.DocumentElement;
            XmlNodeList xmltablelist = xmlroot.GetElementsByTagName("info");
            Hashtable hashtmp = new Hashtable();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                XmlElement tmp = ((XmlElement)xmltablelist[i]);
                Hashtable hashtmpsingle = new Hashtable();
                for (int j = 0; j < tmp.Attributes.Count; j++)
                {
                    hashtmpsingle.Add(tmp.Attributes[j].Name, tmp.Attributes[j].Value);
                }
                hashtmp.Add(tmp.GetAttribute(keyname), hashtmpsingle);
            }
            return hashtmp;
        }

        //解析复杂结构的uiprop
        public static Hashtable givexmlforUIprop(XmlDocument xmldoc, string keyname, string smallKeyname)
        {
            XmlElement xmlroot = xmldoc.DocumentElement;
            XmlNodeList xmltablelist = xmlroot.GetElementsByTagName("weapon");
            Hashtable hashtmp = new Hashtable();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                XmlElement tmp = ((XmlElement)xmltablelist[i]);
                Hashtable hashtmpsingle = new Hashtable();
                for (int j = 0; j < tmp.Attributes.Count; j++)
                {
                    hashtmpsingle.Add(tmp.Attributes[j].Name, tmp.Attributes[j].Value);
                }

                Hashtable propsHash = new Hashtable();
                XmlNodeList titlelist = tmp.GetElementsByTagName("title");
                for (int k = 0; k < titlelist.Count; k++)
                {
                    XmlElement titleTmp = ((XmlElement)titlelist[k]);
                    Hashtable titleTmpItem = new Hashtable();
                    for (int m = 0; m < titleTmp.Attributes.Count; m++)
                    {
                        titleTmpItem.Add(titleTmp.Attributes[m].Name, titleTmp.Attributes[m].Value);
                    }

                    propsHash.Add(titleTmp.GetAttribute(smallKeyname), titleTmpItem);
                }
                hashtmpsingle.Add("props", propsHash);

                hashtmp.Add(tmp.GetAttribute(keyname), hashtmpsingle);
            }
            return hashtmp;
        }

        public static Hashtable givexmlforhash(string xmlstr)
        {
            XmlDocument xmltable = new XmlDocument();
            xmltable.LoadXml(xmlstr);
            XmlElement xmlroot = xmltable.DocumentElement;
            XmlNodeList xmltablelist = xmlroot.ChildNodes;
            Hashtable hashtmp = new Hashtable();
            for (int i = 0; i < xmltablelist.Count; i++)
            {
                hashtmp[xmltablelist[i].Name] = xmltablelist[i].ChildNodes[0].Value;
            }
            return hashtmp;
        }
    }
}