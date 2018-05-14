using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Text;

namespace YunLoong
{
    public static class GNetManager
    {
        public static String SendRequest(String url, Encoding encoding)
        {
            HttpWebRequest webRequest = (HttpWebRequest)WebRequest.Create(url);
            webRequest.Method = "GET";
            HttpWebResponse webResponse = (HttpWebResponse)webRequest.GetResponse();
            StreamReader sr = new StreamReader(webResponse.GetResponseStream(), encoding);
            return sr.ReadToEnd();
        }
        public static String SendRequestForGet(String url, string param)
        {
            return SendRequest(url+"?"+ param, Encoding.UTF8);
        }
        public static string SendRequestForPost(String url, string param)
        {

            HttpWebRequest webRequest = (HttpWebRequest)WebRequest.Create(url);
            webRequest.Method = "POST";
            webRequest.ContentType = "application/x-www-form-urlencoded; charset=UTF-8";

            Stream newStream = webRequest.GetRequestStream();
            byte[] data = Encoding.UTF8.GetBytes(param);
            newStream.Write(data, 0, data.Length);
            newStream.Close();

            HttpWebResponse webResponse = (HttpWebResponse)webRequest.GetResponse();
            StreamReader sr = new StreamReader(webResponse.GetResponseStream(), Encoding.UTF8);
            return sr.ReadToEnd();
        }
        public static string CreateRequestParam(params object[] obj_arr )
        {
            StringBuilder sb = new StringBuilder();
            if (obj_arr.Length % 2 == 0)
            {
                for (int i = 0; i < obj_arr.Length; i += 2)
                {
                    sb.AppendFormat("{0}={1}", obj_arr[i], obj_arr[i + 1]);
                }
                return sb.ToString();
            }
            else
            {
                Console.WriteLine("RequestParam参数数量错误");
                return "";
            }
        }
        public static void TestGNetManager()
        {
            Console.WriteLine(GNetManager.SendRequest("http://localhost/test/testgetpost.php", Encoding.UTF8));
            Console.WriteLine(GNetManager.SendRequestForGet("http://localhost/test/testgetpost.php", CreateRequestParam("param", 6582)));
            Console.WriteLine(GNetManager.SendRequestForPost("http://localhost/test/testgetpost.php", CreateRequestParam("param", 9655)));
        }
    }

}
