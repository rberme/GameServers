using System;
using System.Collections.Generic;
using System.Text;

namespace YunLoong
{

    public class GSqlManager
    {
        public MySqlDataProvider mySqlDataProvider;
        public GSqlManager(GMySqlInfo tGMySqlInfo)
        {
            mySqlDataProvider = new MySqlDataProvider(tGMySqlInfo);
        }

        //查询指定表是否存在指定ID(单主键,映射表对象类型)
        public bool TryFindResultSet(object personalId, Type type)
        {
            Type[] smonotype_Arr = { type };
            return Convert.ToBoolean(mySqlDataProvider.GetType().GetMethod("findData").MakeGenericMethod(smonotype_Arr).Invoke(mySqlDataProvider, new object[] { new object[] { personalId } }));
        }

        public void TryGiveResultSet<T>(object personalId, out T mainEntity) where T : new()
        {
            mySqlDataProvider.giveSingleData(out mainEntity, personalId);
        }
    }
}
