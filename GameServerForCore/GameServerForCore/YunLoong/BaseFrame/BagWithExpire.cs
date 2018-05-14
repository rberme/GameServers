
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading;

namespace YunLoong
{
    public sealed class BagWithExpire<T> : IDisposable
    {

        private readonly ConcurrentDictionary<T, DateTime> pool;
        private Timer expireTimer;//过期计时器
        private int expireTime;//过期时间
        private int isInTimer;//是否正在执行计时操作
        private int minPoolSize;
        /// <summary>
        /// pool min size
        /// </summary>
        public int MinPoolSize
        {
            get { return minPoolSize; }
            set { minPoolSize = value; }
        }

        public BagWithExpire(bool enableExpire = false, int expireTime = 300, int minPoolSize = 5)
        {
            this.expireTime = expireTime;
            this.minPoolSize = minPoolSize;
            pool = new ConcurrentDictionary<T, DateTime>();
            if (enableExpire)//每一分钟检查一次过期
                expireTimer = new Timer(ExpireCheck, null, 60000, 60000);

        }
        private void ExpireCheck(object state)
        {
            if (Interlocked.CompareExchange(ref isInTimer, 1, 0) == 1)
                return;//正在执行上一次过期检查
            try
            {
                List<T> expireList = new List<T>();

                int leftnum = pool.Count;
                foreach (var item in pool)
                {
                    if (leftnum <= minPoolSize) break;
                    if (DateTime.Now.Subtract(item.Value).TotalSeconds >= expireTime)
                    {
                        expireList.Add(item.Key);
                        leftnum--;
                    }
                }
                foreach (var item in expireList)
                {
                    pool.TryRemove(item, out DateTime value);
                    try
                    {
                        IDisposable dispose = item as IDisposable;
                        if (dispose != null)
                            dispose.Dispose();
                        else
                        {

                        }
                    }
                    catch (Exception ex)
                    {

                    }
                }
            }
            finally
            {
                Interlocked.Exchange(ref isInTimer, 0);
            }
        }

        public int Count
        {
            get
            {
                return pool.Count;
            }
        }



        public bool TryTake(ref T item)
        {
            if (pool.Count == 0)
                return false;
            if (pool.TryRemove(item, out DateTime time))
                return true;
            return false;
        }


        public void CopyTo(T[] array, int index)
        {
            T[] Keyarray = pool.Keys.ToArray();
            Keyarray.CopyTo(array, index);
        }

        //public T Create()
        //{
        //    return factory();
        //}

        //public void Put()
        //{
        //    T item = factory();
        //}

        public void Add(T item)
        {

            pool.TryAdd(item, DateTime.Now);
        }


        //停止计时器，释放池里面所有的Item
        public void Dispose()
        {
            if (expireTimer != null)
                expireTimer.Dispose();

            while (pool.Count != 0)
            {
                var result = pool.First();
                pool.TryRemove(result.Key, out DateTime value);
                try
                {
                    IDisposable dispose = result.Key as IDisposable;
                    if (dispose != null)
                        dispose.Dispose();
                    else
                    {
                        //var wcfProxy = result.Item as ICommunicationObject;
                        //var socket = result.Item as Socket;
                        //if (socket != null)
                        //{
                        //    socket.Close();
                        //    socket.Dispose();
                        //}
                    }
                }
                catch (Exception ex)
                {

                }
            }
        }
    }
}