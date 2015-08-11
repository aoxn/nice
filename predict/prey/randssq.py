#!/usr/bin/env python2
# -*- coding: utf-8 -*-

__author__ = 'spacex'

import random,time,numpy as np
import itertools as itto
from ssq import Ssq


class RandomNum(Ssq):

    def __init__(self):
        Ssq.__init__(self)

    def area_freq_red(self, start, end):
        ret = {}
        for num in range(1, 34): ret[num] = 0

        for i in range(start, end+1):
            reds = self.red_ball_row(i)
            for red in reds: ret[int(red)] += 1

        return ret

    def pick_up_ball(self, ln=6):
        ret = {}
        while len(ret.keys()) < ln:
            tmp = random.randint(1, 33)
            if ret.get(tmp, 0) != 0: continue
            ret[tmp] = tmp

        return [k for k,v in ret.items()]

    def filter_over_continuous_in(self,pdt_set, idx, cnt=4):
        if idx < cnt: return
        ret = []
        for p in range(idx-cnt+1, idx):
            tmp = self.intersection(self.red_ball_row(p), self.red_ball_row(p-1))
            ret = self.union(ret, tmp)
        self.log.debug("Filter_A:"+str(ret)+" 过滤掉前几期重号的")
        return self.difference(pdt_set, ret)

    def over_continuous_in(self,pdt_set, idx, cnt=6):
        if idx < cnt: return
        ret = []
        for p in range(idx-cnt+1, idx):
            tmp = self.intersection(self.red_ball_row(p), self.red_ball_row(p-1))
            ret = self.union(ret, tmp)
        self.log.debug("Filter_A:"+str(ret)+" 过滤掉前几期重号的")
        return ret

    def filter_over_continuous_out(self,pdt_set, idx,cnt=19):
        #排除前3期都没出现的，或者前6、7、8、9期都没出现的。像小概率事件
        if idx < cnt: return
        ret = {}
        for p in range(idx-1, idx-cnt, -1):
            for m in self.red_ball_row(p):
                tmp = ret.get(m, 0)
                if tmp: continue
                ret[m] = idx - p
        res = [k for k,v in ret.items() if v == 4 ]
        self.log.debug("Filter_B:"+str(res)+" 排除前3期都没出现的。像小概率事件")
        return self.difference(pdt_set, res)

    def find_cycle(self, idx):
        res = {}
        for i in range(idx-1,0, -1):
            tmp = self.red_ball_row(i)
            for j in tmp:
                res[j] = res.get(j,0) + 1
            if len(res) >=31: return idx - i,res
        return idx,res

    def filter_over_connect_num(self,pdt_set, idx):
        #上一期连号出现的几个，最多只留一个
        return
    def over_connect_repeat(self,pdt_set, idx):
        #与上一期最多三个重号的，平均一个，从上一期中淘汰三个。出现次数少的数，重号的可能性不大
        pre = self.red_ball_row(idx - 1)
        itr = self.intersection(pdt_set, pre)
        #print "xyy:",itr
        if len(itr) <= 3:
            return pdt_set
        #frq = self.area_freq_red(idx-60, idx - 2)
        cyl,frq = self.find_cycle(idx-2)
        #print frq
        tmp = [(k, v) for k,v in frq.items() if k in itr]
        res = self.sort_list_result(tmp, Ssq.VALUE_INDEX)
        res = [k for k, v in res[0:len(itr) - 3]]
        #res = [k for k, v in res[0:len(itr) - 3]]
        self.log.debug( "Filter_C:"+str(res)+" 与上一期最多三个重号，平均一个，从上一期中淘汰三个。出现次数少的数，重号的可能性不大")
        return res
    def filter_over_connect_repeat(self,pdt_set, idx):
        #与上一期最多三个重号的，平均一个，从上一期中淘汰三个。出现次数少的数，重号的可能性不大
        pre = self.red_ball_row(idx - 1)
        itr = self.intersection(pdt_set, pre)
        #print "xyy:",itr
        if len(itr) <= 3:
            return pdt_set
        #frq = self.area_freq_red(idx-60, idx - 2)
        cyl,frq = self.find_cycle(idx-2)
        #print frq
        tmp = [(k, v) for k,v in frq.items() if k in itr]
        res = self.sort_list_result(tmp, Ssq.VALUE_INDEX)
        res = [k for k, v in res[0:len(itr) - 3]]
        #res = [k for k, v in res[0:len(itr) - 3]]
        self.log.debug( "Filter_C:"+str(res)+" 与上一期最多三个重号，平均一个，从上一期中淘汰三个。出现次数少的数，重号的可能性不大")
        return self.difference(pdt_set, res)

    def filter_over_last_cycle(self,pdt_set, idx):
        #上一个全量周期中，出现次数最少的，次数为1的
        res, px = {},0
        for i in range(idx-1,0, -1):
            tmp = self.red_ball_row(i)
            for j in tmp:
                res[j] = res.get(j,0) + 1
            if len(res) >=33:
                px = idx - i
                break
        tmp = [k for k,v in res.items() if v==1]
        f =  tmp[0:4] if len(tmp)>=4 else tmp
        self.log.debug( "Filter_D:",str(f))
        return self.difference(pdt_set,f)

    def filter_over_last_cycle_group(self,pdt_set, idx):
        #上一个全量周期中，出现次数最少的，次数为1的
        res, px = {},0
        for i in range(idx-1,0, -1):
            tmp = self.red_ball_row(i)
            for j in tmp:
                res[j] = res.get(j,0) + 1
            if len(res) >=33:
                px = idx - i
                break
        tmp = [k for k,v in res.items() if v==1]
        f =  tmp[0:4] if len(tmp)>=4 else tmp
        self.log.debug("Filter_D:"+str(f))
        return self.difference(pdt_set,f)

    def filter_over_direct(self,pdt_set,f_set):

        return self.difference(pdt_set,f_set)

    def test_cycle(self,idx):
        res, px = {},0
        for i in range(idx-1,0, -1):

            tmp = self.red_ball_row(i)
            #print tmp
            for j in tmp:
                res[j] = res.get(j,0) + 1
            if len(res) >=33:
                break
        return res

    def test(self,):
        cnt = 0
        for idx in range(50,1770):
            tmp = self.test_cycle(idx)
            print "------------------------------------------------------------------------"
            m=self.sort_list_result(tmp.items(),Ssq.VALUE_INDEX,True)
            print m
            print self.red_ball_row(idx+1)
            #p = [k for k,v in tmp.items() if v ==1]
            p=[m[32][0]]
            print "红外",p
            x= self.intersection(self.red_ball_row(idx+1),p)
            print x
            if len(x)==1:cnt+=1

        print cnt,float(cnt)/(1770-50)

    def predict_num(self, idx):
        res = self.pick_up_ball(33)
        res = self.filter_over_continuous_in(res, idx)
        #res = self.filter_over_continuous_out(res,idx)
        #res = self.filter_over_connect_repeat(res,idx)
        #res = self.filter_over_connect_repeat(res,idx-6)
        #res = self.filter_over_connect_repeat(res,idx-7)
        #res = self.filter_over_connect_repeat(res,idx-8)
        #res = self.filter_over_connect_repeat(res,idx-9)
        #res = self.filter_over_connect_repeat(res,idx-10)
        #res= self.filter_over_last_cycle(res,idx)
        #res =  self.filter_over_direct(res,[1,2,5,8,11,24,29,32,6])
        return res

    def test_random_round(self, start, end, ln =30):

        for i in range(start,end):
            rnd = self.pick_up_ball(ln)
            r_i = self.red_ball_row(i)
            #print r_i,rnd, self.intersection(r_i,rnd)
            print self.intersection(r_i,rnd)
        return


    def search(self,l):
        for i in range(0,1777):
            r =self.red_ball_row(i)
            if len(self.intersection(r,l))>=4:
                print i,r

    def reward_cnt(self,prd,cur,cnt):
        r_6,r_5,r_4 = 0,0,0
        for item in prd:
            #print list(item),cur
            l = len(self.intersection(cur,list(item)))
            if l==6:
                r_6+=1
            elif l==5:
                r_5+=1
            elif l==4:
                r_4+=1

        return (r_6,r_5,r_4)

    def combinate_3_road(self, pre_select, rule, item):
        ret = []
        first = [k for k in pre_select if 1 <= k <= 11]
        middle= [k for k in pre_select if 12 <= k <= 24]
        end   = [k for k in pre_select if 25 <= k <= 33]
        f = itto.combinations(first,rule[0])

        for fi in f:
            m = itto.combinations(middle,rule[1])
            for mi in m:
                e = itto.combinations(end,rule[2])
                for ei in e:
                    tmp = fi+mi+ei
                    if not self.filter(item,tmp): ret.append(tmp)

        #print ret
        return ret

    def combinate_2_road(self, pre_select, rule, item):
        ret = []
        first = [k for k in pre_select if 1 <= k <= 16]
        middle= [k for k in pre_select if 17 <= k <= 33]
        f = itto.combinations(first,rule[0])

        for fi in f:
            m = itto.combinations(middle,rule[1])
            for mi in m:
                tmp = fi+mi
                if not self.filter(item,tmp): ret.append(tmp)
        return ret

    def filter(self,s,tmp):
        for t in s:
            if len(self.intersection(t,tmp))==len(t):
                return True

        return False

    def filter_less_similarity(self,target,cnt=4):
        t_len = len(self.data)
        for it in target:
            t_cnt = 0
            for tmp in range(t_len - 600, t_len):
                if len(self.intersection(it,self.red_ball_row(tmp)))>=4:
                    t_cnt += 1
                if t_cnt >= cnt: break
            if t_cnt < cnt: target.remove(it)
        return target


    def main(self,num):
        self.log.debug("总期数："+str(len(self.data))+" 待预测期："+str(num+1))
        prd = self.predict_num(num)
        #prd = [1, 2, 3,  7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 25, 27, 28, 30, 31, 33]
        #prd = self.difference(prd,[02, 18, 19, 27, 9, 8, 15])

        self.log.debug("预测集合 "+str(len(prd))+" "+str(prd))
        ret = self.combinate_2_road(prd,[2,4],[])
        self.log.debug("总的组合数 "+str(len(ret)))
        ret = self.filter_less_similarity(ret)
        self.log.debug("过滤后总的组合数 "+str(len(ret)))

        if num <len(self.data):
            cur = self.red_ball_row(num)
            cnt = self.reward_cnt(ret,cur,5)
            self.log.debug("第 "+str(num+1)+" 期，真实开奖"+str(cur)+" 中奖号码："+str(self.intersection(prd,cur))+" 中奖个数："+str(cnt))
        return






if __name__ == "__main__":
    rnd = RandomNum()
    num = len(rnd.data)-1
    np.set_printoptions(threshold='nan', linewidth=300)
    rnd.test_random_round(0,100,8)
    #rnd.test()
    #print rnd.sort_list_result(rnd.test_cycle(1779).items(),Ssq.VALUE_INDEX,True)
    #print "dddd"
    #print rnd.red_ball_row(1770)
    #ret = rnd.combinate([1,3,7,13,14,15,30,33],[3,2,1],[[13,14]])
    #print len(rnd.data)
    #rnd.search([2 ,4,6,9,10,19,28,30,31,33])
    #rnd.main(num)
    #print prd,len(prd)
    #print rnd.red_ball_row(num),len(prd),rnd.intersection(rnd.red_ball_row(num),prd), prd




