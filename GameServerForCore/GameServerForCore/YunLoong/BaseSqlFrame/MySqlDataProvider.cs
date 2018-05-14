using System;
using System.Collections.Generic;
using System.Data;
using System.Data.Common;
using System.Linq;
using System.Reflection;
using System.Text;
using MySql.Data.MySqlClient;
using ProtoBuf;

namespace YunLoong
{

    /// <summary>
    /// MSSQL数据库服务提供者
    /// </summary>
    public class MySqlDataProvider
    {
        public GMySqlInfo mySqlInfo;
        public string ConnectionString
        {
            get
            {
                return $"Data Source={mySqlInfo.mysqlIP};Port={mySqlInfo.mysqlPort};Database={mySqlInfo.mysqlDatabase};User Id={mySqlInfo.mysqlUserID};Password={mySqlInfo.mysqlPassword};pooling=true;MaximumPoolSize=100;CharSet=utf8;port={mySqlInfo.mysqlPort};SslMode=none";
            }
        }

        //测试 MySqlDataProvider tMySqlDataProvider = new MySqlDataProvider(GMySqlInfo.getGServerInfo());
        public MySqlDataProvider(GMySqlInfo tMysqlInfo)
        {
            mySqlInfo = tMysqlInfo;
        }
        public MySqlConnection getMySqlCon()
        {
            MySqlConnection mysql = new MySqlConnection(ConnectionString);
            return mysql;
        }
        public MySqlCommand getSqlCommand(MySqlConnection mysqlConn)
        {
            MySqlCommand mySqlCommand = new MySqlCommand();
            mySqlCommand.Connection = mysqlConn;
            return mySqlCommand;
        }

        //保存指定存储对象到表
        public void saveData<T>(params T[] obj_arr)
        {


            string saveStr = @"INSERT INTO {0}({1})VALUES{2} ON DUPLICATE KEY UPDATE {3}";

            if (obj_arr.Length <= 0) return;
            Type dataType = typeof(T);

            string tableName = dataType.Name;
            StringBuilder insertFields = new StringBuilder();
            StringBuilder updateSets = new StringBuilder();

            List<string> fNameList = new List<string>();
            List<string> mainKeyList = new List<string>();

            MemberInfo[] tMemberInfo = dataType.GetMembers();


            checkfieldName(tMemberInfo, ref fNameList, ref mainKeyList);

            if (fNameList.Count > 0)
            {
                for (int i = 0; i < fNameList.Count; i++)
                {
                    if (i != 0) { insertFields.Append(","); }
                    insertFields.Append(fNameList[i]);

                }
                int upcount = 0;
                for (int i = 0; i < fNameList.Count; i++)
                {
                    if (!mainKeyList.Contains(fNameList[i]))
                    {
                        if (upcount != 0) { updateSets.Append(","); }
                        updateSets.AppendFormat("{0}=VALUES({0})", fNameList[i]);
                        upcount++;
                    }
                }

                StringBuilder insertValues = new StringBuilder();
                for (int i = 0; i < obj_arr.Length; i++)
                {
                    StringBuilder valuestr = new StringBuilder();
                    for (int j = 0; j < fNameList.Count; j++)
                    {
                        if (j != 0) { valuestr.Append(","); }
                        valuestr.AppendFormat("@P{0}_{1}", i, j);
                    }
                    if (i != 0) { insertValues.Append(","); }
                    insertValues.AppendFormat("({0})", valuestr);
                }

                MySqlConnection sqlConnection = getMySqlCon();
                sqlConnection.Open();
                MySqlCommand tMySqlCommand = getSqlCommand(sqlConnection);



                for (int i = 0; i < obj_arr.Length; i++)
                {
                    for (int j = 0; j < fNameList.Count; j++)
                    {
                        object fvalue = GameData.givevalue(obj_arr[i], fNameList[j]);
                        tMySqlCommand.Parameters.AddWithValue(String.Format("@P{0}_{1}", i, j), fvalue);
                    }
                }

                StringBuilder command = new StringBuilder();
                command.AppendFormat(saveStr, tableName, insertFields, insertValues, updateSets);
                tMySqlCommand.CommandText = command.ToString();

                //Console.WriteLine(command);

                tMySqlCommand.ExecuteScalar();
                //int rows = Convert.ToInt32(tMySqlCommand.ExecuteScalar());
                //Console.WriteLine($"rows:{rows}");

                sqlConnection.Close();
            }
        }

        public void checkfieldName(MemberInfo[] tMemberInfo_arr, ref List<string> fNameList, ref List<string> mainKeyList)
        {
            for (int i = 0; i < tMemberInfo_arr.Length; i++)
            {
                EntityFieldAttribute tEntityField = tMemberInfo_arr[i].GetCustomAttribute<EntityFieldAttribute>();
                if (tEntityField != null)
                {
                    string tfieldName;
                    if (!tEntityField.FieldName.IsEmpty())
                    {
                        tfieldName = tEntityField.FieldName;
                    }
                    else
                    {
                        tfieldName = tMemberInfo_arr[i].Name;
                    }
                    if (tEntityField.IsKey)
                    {
                        mainKeyList.Add(tfieldName);
                    }
                    fNameList.Add(tfieldName);
                }
            }
        }
        private MySqlCommand createSeleteCommand<T>(string giveStr, out List<string> fNameList, out MySqlConnection sqlConnection, params object[] keys)
        {
            Type dataType = typeof(T);
            string tableName = dataType.Name;
            StringBuilder selectCondition = new StringBuilder();
            fNameList = new List<string>();
            Dictionary<int, string> mainKeyDic = new Dictionary<int, string>();

            MemberInfo[] tMemberInfo = dataType.GetMembers();

            checkFieldMainKey(tMemberInfo, ref fNameList, ref mainKeyDic);

            if (mainKeyDic.Count > 0)
            {

                sqlConnection = getMySqlCon();
                sqlConnection.Open();
                MySqlCommand tMySqlCommand = getSqlCommand(sqlConnection);

                int tLen = Math.Min(mainKeyDic.Count, keys.Length);
                int tcount = 0;
                for (int i = 0; i < tLen; i++)
                {
                    if (keys[i] != null)
                    {
                        if (tcount != 0) { selectCondition.Append(" and "); }
                        selectCondition.AppendFormat("{0}=@P_{1}", mainKeyDic[i], i);
                        tMySqlCommand.Parameters.AddWithValue(String.Format("@P_{0}", i), keys[i]);
                        tcount++;
                    }
                }

                StringBuilder command = new StringBuilder();
                command.AppendFormat(giveStr, tableName, selectCondition);
                tMySqlCommand.CommandText = command.ToString();

                //Console.WriteLine(command);
                return tMySqlCommand;
            }
            else
            {
                sqlConnection = null;
                return null;
            }
        }

        //从表中获取指定ID的指定对象
        public bool findData<T>(params object[] keys) where T : new()
        {
            //T[] t_arr;
            //giveData(out t_arr, keys);
            //return t_arr.Length > 0;
            MySqlConnection sqlConnection;
            List<string> fNameList;
            string giveStr = "SELECT count(0) FROM {0} WHERE {1};";
            MySqlCommand tMySqlCommand = createSeleteCommand<T>(giveStr, out fNameList, out sqlConnection, keys);
            if (tMySqlCommand != null)
            {
                int num = giveResultNum(tMySqlCommand);
                sqlConnection.Close();
                return num > 0;
            }
            else
            {
                return false;
            }
        }

        public void giveSingleData<T>(out T tSqlTestTableData, params object[] keys) where T : new()
        {
            tSqlTestTableData = default(T);

            MySqlConnection sqlConnection;
            List<string> fNameList;
            string giveStr = "SELECT * FROM {0} WHERE {1};";
            MySqlCommand tMySqlCommand = createSeleteCommand<T>(giveStr, out fNameList, out sqlConnection, keys);
            if (tMySqlCommand != null)
            {
                T[] tSqlTestTableData_arr;
                giveResultSet<T>(out tSqlTestTableData_arr, fNameList, tMySqlCommand);
                sqlConnection.Close();
                if (tSqlTestTableData_arr != null)
                {
                    if (tSqlTestTableData_arr.Length > 0)
                    {
                        tSqlTestTableData = tSqlTestTableData_arr[0];

                        if (tSqlTestTableData_arr.Length > 1)
                        {
                            StringBuilder errerBuider = new StringBuilder();
                            errerBuider.Append($"giveSingleData查询{typeof(T).ToString()}结果不唯一:");
                            for (int er = 0; er < keys.Length; er++)
                            {
                                errerBuider.Append($"{keys[er]} ");

                            }
                            Console.WriteLine(errerBuider);
                        }
                    }
                }
            }
        }
        public void giveData<T>(out T[] tSqlTestTableData_arr, params object[] keys) where T : new()
        {

            MySqlConnection sqlConnection;
            List<string> fNameList;
            string giveStr = "SELECT * FROM {0} WHERE {1};";
            MySqlCommand tMySqlCommand = createSeleteCommand<T>(giveStr, out fNameList, out sqlConnection, keys);
            if (tMySqlCommand != null)
            {
                giveResultSet<T>(out tSqlTestTableData_arr, fNameList, tMySqlCommand);
                sqlConnection.Close();
            }
            else
            {
                tSqlTestTableData_arr = null;
            }
        }

        public void giveKeyForField<T>(out HashSet<object> keyHashSet, string field) where T : new()
        {
            string giveStr = "SELECT {1} FROM {0};";
            Type dataType = typeof(T);
            string tableName = dataType.Name;
            StringBuilder selectCondition = new StringBuilder();

            MySqlConnection sqlConnection = getMySqlCon();
            sqlConnection.Open();
            MySqlCommand tMySqlCommand = getSqlCommand(sqlConnection);

            StringBuilder command = new StringBuilder();
            command.AppendFormat(giveStr, tableName, field);
            tMySqlCommand.CommandText = command.ToString();

            //Console.WriteLine(command);
            giveResultKeyHashSet(out keyHashSet, tMySqlCommand);


            sqlConnection.Close();

        }
        public static int giveResultNum(MySqlCommand mySqlCommand)
        {
            MySqlDataReader reader = mySqlCommand.ExecuteReader();
            try
            {
                while (reader.Read())
                {
                    if (reader.HasRows)
                    {
                        return reader.GetInt32(0);
                    }
                }
            }
            catch (Exception e)
            {

                Console.WriteLine($"查询失败了！{e}");
            }
            finally
            {
                reader.Close();
            }
            return 0;
        }
        public static void giveResultKeyHashSet(out HashSet<object> keyHashSet, MySqlCommand mySqlCommand)
        {
            keyHashSet = new HashSet<object>();
            MySqlDataReader reader = mySqlCommand.ExecuteReader();
            try
            {
                while (reader.Read())
                {
                    if (reader.HasRows)
                    {
                        try
                        {
                            keyHashSet.Add(reader.GetValue(0));
                        }
                        catch
                        {
                        }
                    }
                }
            }
            catch (Exception e)
            {
                Console.WriteLine($"查询失败了！{e}");
            }
            finally
            {
                reader.Close();
            }
        }

        public static void giveResultSet<T>(out T[] tSqlTestTableData_arr, List<string> fNameList, MySqlCommand mySqlCommand) where T : new()
        {
            MySqlDataReader reader = mySqlCommand.ExecuteReader();
            try
            {

                List<T> tSqlTestTableDataList = new List<T>();
                while (reader.Read())
                {
                    if (reader.HasRows)
                    {
                        T item = new T();
                        for (int i = 0; i < reader.FieldCount; i++)
                        {
                            try
                            {
                                if (!(reader[fNameList[i]] is System.DBNull))
                                {
                                    GameData.valuetofield(item, fNameList[i], reader[fNameList[i]]);
                                }
                                //Console.WriteLine($"{fNameList[i]}={reader[fNameList[i]]}");
                            }
                            catch (Exception e)
                            {

                                Console.WriteLine($"字段不存在！{fNameList[i]}");
                            }
                        }
                        tSqlTestTableDataList.Add(item);
                    }
                }
                tSqlTestTableData_arr = tSqlTestTableDataList.ToArray();
            }
            catch (Exception e)
            {

                Console.WriteLine($"查询失败了！{e}");
                tSqlTestTableData_arr = null;
            }
            finally
            {
                reader.Close();
            }
        }
        public void deleteData<T>(params object[] keys)
        {
            string deleteStr = "DELETE FROM {0} WHERE {1};";

            Type dataType = typeof(T);
            string tableName = dataType.Name;
            StringBuilder deleteCondition = new StringBuilder();
            List<string> fNameList = new List<string>();
            Dictionary<int, string> mainKeyDic = new Dictionary<int, string>();

            MemberInfo[] tMemberInfo = dataType.GetMembers();

            checkFieldMainKey(tMemberInfo, ref fNameList, ref mainKeyDic);

            if (mainKeyDic.Count > 0)
            {

                MySqlConnection sqlConnection = getMySqlCon();
                sqlConnection.Open();
                MySqlCommand tMySqlCommand = getSqlCommand(sqlConnection);

                int tLen = Math.Min(mainKeyDic.Count, keys.Length);
                int tcount = 0;
                for (int i = 0; i < tLen; i++)
                {
                    if (keys[i] != null)
                    {
                        if (tcount != 0) { deleteCondition.Append(" and "); }
                        deleteCondition.AppendFormat("{0}=@P_{1}", mainKeyDic[i], i);
                        tMySqlCommand.Parameters.AddWithValue(String.Format("@P_{0}", i), keys[i]);
                        tcount++;
                    }
                }

                StringBuilder command = new StringBuilder();
                command.AppendFormat(deleteStr, tableName, deleteCondition);
                tMySqlCommand.CommandText = command.ToString();

                //Console.WriteLine(command);

                tMySqlCommand.ExecuteScalar();
                //int rows = Convert.ToInt32(tMySqlCommand.ExecuteScalar());
                //Console.WriteLine($"rows:{rows}");

                sqlConnection.Close();
            }


        }
        public void checkFieldMainKey(MemberInfo[] tMemberInfo_arr, ref List<string> fNameList, ref Dictionary<int, string> mainKeyDic)
        {
            for (int i = 0; i < tMemberInfo_arr.Length; i++)
            {
                EntityFieldAttribute tEntityField = tMemberInfo_arr[i].GetCustomAttribute<EntityFieldAttribute>();
                if (tEntityField != null)
                {
                    string tfieldName;
                    if (!tEntityField.FieldName.IsEmpty())
                    {
                        tfieldName = tEntityField.FieldName;
                    }
                    else
                    {
                        tfieldName = tMemberInfo_arr[i].Name;
                    }
                    if (tEntityField.IsKey)
                    {
                        mainKeyDic.Add(tEntityField.Keyid, tfieldName);
                    }
                    fNameList.Add(tfieldName);
                }
            }
        }

    }

    public struct GMySqlInfo
    {
        public int sid;
        public string mysqlIP;
        public string mysqlPort;
        public string mysqlUserID;
        public string mysqlPassword;
        public string mysqlDatabase;

        public static GMySqlInfo getGServerInfo()
        {
            return new GMySqlInfo()
            {
                sid = 1,
                mysqlIP = "192.168.0.205",
                mysqlPort = "3306",
                mysqlUserID = "gzb",
                mysqlPassword = "0",
                mysqlDatabase = "barriers_1",
            };
        }
    }
}
