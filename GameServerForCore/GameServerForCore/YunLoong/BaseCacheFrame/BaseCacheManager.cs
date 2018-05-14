using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Text;
using System.Threading;

namespace YunLoong
{
    public abstract class BaseCacheManager<T, U> where T : BaseCacheLock<U>, new() where U : BaseCacheStruct, new()
    {

        public ConcurrentDictionary<object, T> CacheLockDic;
        public LazyConcurrentDictionary<object, T> CacheLockLazyDic;
        //数据修改操作的HashSet
        public HashSet<BaseEntity> changeDataHashSet = new HashSet<BaseEntity>();
        protected HashSet<object> keyHashSet;
        private ReaderWriterLockSlim keyHashSetLock;
        public GSqlManager gSqlManager;
        public BaseCacheManager(GMySqlInfo tGMySqlInfo)
        {
            gSqlManager = new GSqlManager(tGMySqlInfo);
            initCacheManager();
        }
        public BaseCacheManager()
        {
            initCacheManager();
        }
        public void initCacheManager()
        {
            CacheLockDic = new ConcurrentDictionary<object, T>();
            CacheLockLazyDic = new LazyConcurrentDictionary<object, T>();
            keyHashSet = new HashSet<object>();
            keyHashSetLock = new ReaderWriterLockSlim();
            initkeyHashSet();
        }
        public abstract void initkeyHashSet();
        //public bool ContainKey(int personalId)
        //{
        //    return CacheLockDic.ContainsKey(personalId);
        //}

        public bool TryFindKey(int personalId, out T data)
        {
            return getOrCreateKey(personalId, out data, ReadCacheMode.readOnly);
        }

        public bool TryFindOrCreateKey(object personalId, out T data)
        {
            return getOrCreateKey(personalId, out data, ReadCacheMode.readAndCreate);
        }

        private bool getOrCreateKey(object personalId, out T data, ReadCacheMode mode = ReadCacheMode.readOnly)
        {
            data = default(T);
            if (!CacheLockDic.TryGetValue(personalId, out data))
            {
                keyHashSetLock.EnterReadLock();
                bool hasDataForSql = keyHashSet.Contains(personalId);
                keyHashSetLock.ExitReadLock();
                if (!hasDataForSql)
                {
                    Type mainEntity = typeof(T).GetProperty("mainEntity").GetValue(null) as Type;
                    hasDataForSql = gSqlManager.TryFindResultSet(personalId, mainEntity);
                }
                if (hasDataForSql)
                {
                    try
                    {
                        data = CacheLockLazyDic.GetOrAdd(personalId, (object Tkey) =>
                        {
                            //创建添加用户Lock
                            T tdata = new T();
                            using (tdata.Write())
                            {
                                tdata.Data.cacheID = personalId;
                                if (!tdata.Data.isinit)
                                {
                                    //加载数据实体
                                    tdata.Data.initCacheData();
                                    CacheLockDic.TryAdd(personalId, tdata);
                                }
                            }
                            return tdata;
                        });
                    }
                    catch
                    {
                        return false;
                    }
                }
                else if (mode == ReadCacheMode.readAndCreate)
                {
                    try
                    {
                        data = CacheLockLazyDic.GetOrAdd(personalId, (object Tkey) =>
                        {
                            //创建添加用户Lock
                            T tdata = new T();
                            using (tdata.Write())
                            {
                                tdata.Data.cacheID = personalId;
                                if (!tdata.Data.isinit)
                                {
                                    //创建数据实体
                                    tdata.Data.createCacheData();
                                    CacheLockDic.TryAdd(personalId, tdata);
                                }
                                keyHashSetLock.EnterWriteLock();
                                keyHashSet.Add(personalId);
                                keyHashSetLock.ExitWriteLock();
                            }
                            
                            return tdata;
                        });
                    }
                    catch
                    {
                        return false;
                    }
                }
            }
            return data != default(T);
        }
    }
    enum ReadCacheMode
    {
        readOnly,
        readAndCreate
    }
}


