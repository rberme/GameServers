using System.Collections;
using System.Text;
using System.IO;
using System;
namespace YunLoong
{
    public class Gbytes
    {
        //写入数据
        public BinaryWriter bw;
        private MemoryStream rms;
        private MemoryStream wms;

        //读取数据
        public BinaryReader br;
        public int bufferSize
        {
            get
            {
                if (buffer != null)
                {
                    return buffer.Length;
                }
                else
                {
                    return 0;
                }
            }
        }

        public byte[] buffer;

        public Gbytes()
        {
        }

        internal void init(int bufferSize)
        {
            buffer = new byte[bufferSize];
        }

        public BinaryReader reader()
        {
            rms = new MemoryStream(buffer);
            br = new BinaryReader(rms);
            return br;
        }
        public byte[] rbytes
        {
            get
            {
                return buffer;
            }
        }

        public byte[] wbytes
        {
            get
            {
                if (wms != null)
                {
                    return wms.ToArray();
                }
                else
                {
                    return new byte[0];
                }
            }
        }

        public BinaryWriter writer()
        {
            wms = new MemoryStream();
            bw = new BinaryWriter(wms);
            return bw;
        }

        public BinaryReader giveBinaryReader(byte[] rtbytes)
        {
            rms = new MemoryStream(rtbytes);
            br = new BinaryReader(rms);
            return br;
        }

        public void Dispose()
        {
            buffer = null;
        }
    }
}