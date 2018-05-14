using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace YunLoong
{
    public static class GeneralShortCut
    {
        public static List<T> clone<T>(this List<T> target)
        {
            List<T> cloneList = new List<T>();
            for (int i = 0; i < target.Count; i++)
            {
                cloneList.Add(target[i]);
            }
            return cloneList;
        }
    }
}

