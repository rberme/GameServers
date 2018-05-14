using System;
using System.Collections.Generic;
using System.Text;

namespace YunLoong
{
    public abstract class BaseCacheLock<T> : UsingLock<T> where T : BaseCacheStruct, new()
    {
        public BaseCacheLock() : base(new T()) { }

        //设置HashSet
        public void changeDataSet(BaseEntity obj)
        {
            //data.changeDataHashSet.Add(obj);
        }
        //取值并清空HashSet
    }
}
