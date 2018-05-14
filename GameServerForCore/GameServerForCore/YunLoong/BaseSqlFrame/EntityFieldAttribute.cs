using System;
using System.Collections.Generic;
using System.Text;

namespace YunLoong
{
    [AttributeUsage(AttributeTargets.Field | AttributeTargets.Property, AllowMultiple = false)]
    public class EntityFieldAttribute : Attribute
    {
        /// <summary>
        /// 字段名
        /// </summary>
        public string FieldName { get; set; }
        /// <summary>
        /// 是否主键
        /// </summary>
        public bool IsKey { get; set; }
        /// <summary>
        /// 禁用或排除数据库取值
        /// </summary>
        public bool Disable { get; set; }
        /// <summary>
        /// 列允许为空
        /// </summary>
        public bool Isnullable { get; set; }
        /// <summary>
        /// key的sortID
        /// </summary>
        public int Keyid { get; set; }
        public EntityFieldAttribute()
        {
            FieldName = string.Empty;
            Isnullable = true;
        }
        public EntityFieldAttribute(string fieldName)
            : this()
        {
            FieldName = fieldName;
        }
        public EntityFieldAttribute(bool isKey, int keyid = 0)
           : this()
        {
            IsKey = isKey;
            Keyid = keyid;
        }
    }
}
