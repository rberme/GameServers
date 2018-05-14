using GameServer.CsScript.Model;
using System;
using System.Collections.Generic;
using System.Text;

namespace GameServer.CsScript.Model
{
    public static class DefaultModelData
    {
        public static PlayerData defaultPlayerData(int userID, int serverID)
        {
            return new PlayerData()
            {
                pid = userID * 100000L + serverID,
                userid = userID,
                serverid = serverID,
                pname = "",
                picon = "",
                piconborder = "",
                pexp = 0,
                pdiamond = 0,
                pgold = 0,
                pspirit = 0,
                pmedal = 0,
                phonor = 0,
                pcreatetime = DateTime.Now,
                plastlogintime = DateTime.Now,
                prandseed = 0,
                pguildid = 0,
                plastspiritsddtime = DateTime.Now
            };
        }

        public static PlayerHero defaultPlayerHero(long playerid, int heroID, int quality = 1)
        {
            return new PlayerHero()
            {
                pid = playerid,
                heroid = heroID,
                hexp = 0,
                hquality = quality,
                hawaken = 0,
                hequip1 = 0,
                hequip2 = 0,
                hequip3 = 0,
                hequip4 = 0,
                hequip5 = 0,
                hequip6 = 0,
                hskillrock1 = 0,
                hskillrock2 = 0,
                hskillrock3 = 0,
                hskillrock4 = 0,
                hskillrock5 = 0,
                hskillrock6 = 0,
                hskillrock7 = 0,
                hskillrock8 = 0,
                hskin = 0
            };
        }

        public static PlayerEquip defaultPlayerEquip(long playerid, int SID, int equipID, int rand, int awakenRand)
        {
            return new PlayerEquip()
            {
                pid = playerid,
                sid = SID,
                equipid = equipID,
                eexp = 0,
                erand = rand,
                eawakenrand = awakenRand,
                eawakencount = 0,
                ecreateticks = DateTime.Now,
                eowner = 0,
                eislock = false,
                estate = 1
            };
        }

        public static PlayerRock defaultPlayerRock(long playerid,int rockID,int Num) {
            return new PlayerRock()
            {
                pid = playerid,
                rockid = rockID,
                num = Num,
            };
        }
        public static PlayerItem defaultPlayerItem(long playerid, int itemID, int Num)
        {
            return new PlayerItem()
            {
                pid = playerid,
                itemid = itemID,
                num = Num,
            };
        }
    }
}


