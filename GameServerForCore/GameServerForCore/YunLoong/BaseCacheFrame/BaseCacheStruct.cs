using System;
using System.Collections.Generic;
using System.Text;

namespace YunLoong
{
    public abstract class BaseCacheStruct
    {
        //主实体对象的ID与Cachelock和CacheStruct的ID一致,且主键唯一,(作用,判定CacheStruct对象是否存在)
        public bool isinit;
        public object cacheID;
        public abstract GSqlManager gSqlManager { get; set; }
        public BaseCacheStruct()
        {
            
            //Console.WriteLine("数据对象构造器*****************");
        }
        public virtual void initCacheData()
        {
            isinit = true;
        }
        public virtual void createCacheData()
        {
            isinit = true;
        }


    }
}