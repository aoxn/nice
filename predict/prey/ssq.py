#!/usr/bin/env python2
# -*- coding: utf-8 -*-

__author__ = 'spacex'

import urllib,os,math,time
import numpy as np
import itertools as itto
import logging,logging.config


class Ssq:
    """数据结构
            self.data:
            red_ball : [1,2,3,4]
            freq     : [(k,v),(k1,v1)]

    """
    KEY_INDEX   = 0
    VALUE_INDEX = 1

    def __init__(self):
        logging.config.fileConfig("logger.conf")
        self.log      = logging.getLogger("main")
        self.ssq_file = "ssq.txt"
        self.ssq_url  = "http://www.17500.cn/getData/ssq.TXT"
        self.data     = None
        self.init_data()

    def init_data(self):
        self.retrieve_ssq_file()
        if not self.data:
            self.data =  self.load_data()

    def retrieve_ssq_file(self):
        if not os.path.isfile(os.path.abspath(self.ssq_file)):
            urllib.urlretrieve(self.ssq_url, self.ssq_file)
        return

    def load_data(self):
        data = []
        with open(self.ssq_file, "r") as fd:
            for line in fd:
                tmp = line.split()
                data.append({"num": tmp[0], "date": tmp[1], "r1": tmp[2], "r2": tmp[3], "r3": tmp[4],
                             "r4": tmp[5], "r5": tmp[6], "r6": tmp[7], "b1": tmp[8]})
        return data

    def red_ball_row(self, idx):
        ret = []
        row = self.data[idx]
        for i in range(1, 7):
            ret.append(int(row["r%s" % i]))
        return ret

    def blue_ball_row(self, idx):
        return [self.data[idx]['b1']]

    def sort_list_result(self, fq, target=0, iv=False):
        """
            [(k,v),[k1,v1]]
        :param fq:
        :param target:
        :param iv:
        :return:
        """
        return sorted(fq, key=lambda item: item[target], reverse=iv)

    def area_freq_red(self, start, end):
        ret = {}
        for num in range(1, 34): ret[num] = 0

        for i in range(start, end+1):
            reds = self.red_ball_row(i)
            for red in reds: ret[int(red)] += 1

        return [(k, v) for k, v in ret.items()]

    def intersection(self, main, second):
        """
            两个集合的交集
        :param main:
        :param second:
        :return:
        """
        return list(set(main).intersection(set(second)))

    def union(self, main, second):
        """
            两个集合的并集
        :param main:
        :param second:
        :return:
        """
        return list(set(main).union(set(second)))
    def difference(self, main, second):
        """
            两个集合的交集
        :param main:
        :param second:
        :return:
        """
        return list(set(main).difference(set(second)))
    def possiblity(self, num, rank=1):
        """
            打印前N期的预测值，在后面几期中出现的位置
        :param num:
        :param rank: 代表几等奖
        :return:
        """

        return

    def print_freq(self, start, end, sort_key=VALUE_INDEX,iv=True):
        print self.sort_list_result(self.area_freq_red(start,end), sort_key, iv=iv)

    def search_result(self, start, cnt, win_list):
        ret = []
        for i in range(start, start + cnt):
            red = self.red_ball_row(i)
            ret.append((red, win_list, self.intersection(win_list, red)))
        return ret

    def write2file(self, word,file='result.txt'):
        with open(file,"a+") as fd:
            fd.writelines(word)
        return

    def remove_file(self,file):
        if os.path.isfile(os.path.abspath(file)):os.remove(file)
        return


class PatternGraph(Ssq):

    """
        每期的中奖号码出现的模式
    """

    def __init__(self):
        Ssq.__init__(self)
        self.pattern_set = None

    def search_pattern(self, start, cnt):
        ret = {}
        for i in range(start, cnt):
            ball_a_red  = self.red_ball_row(i)
            ball_a_blue = self.blue_ball_row(i)
            for j in range(i+1, len(self.data)):
                ball_b_red  = self.red_ball_row(j)
                ball_b_blue = self.blue_ball_row(j)
                inter_r  = self.intersection(ball_a_red,ball_b_red)
                inter_b  = self.intersection(ball_a_blue,ball_b_blue)

                if len(inter_r) >= 5:
                    print "[%s:%s:]%s|%s:%s|%s:%s|%s"%(i,j,ball_a_red,ball_a_blue[0], ball_b_red, ball_b_blue[0], inter_r,inter_b)
            #break
        return ret

    def search_pattern_reverse(self, start, end):
        ret = {}
        for i in range(start, end, -1):
            ball_a_red  = self.red_ball_row(i)
            ball_a_blue = self.blue_ball_row(i)
            balls_union, balls_all = [], []
            for j in range(i -1, 0,-1):
                ball_b_red  = self.red_ball_row(j)
                ball_b_blue = self.blue_ball_row(j)
                inter_r  = self.intersection(ball_a_red,ball_b_red)
                #inter_b  = self.intersection(ball_a_blue,ball_b_blue)
                if len(inter_r) >= 4:
                    #self.write2file("[%s:%s:]%s|%s:%s|%s:%s|%s\n"%(i,j,ball_a_red,ball_a_blue[0], ball_b_red, ball_b_blue[0], inter_r,inter_b),"pattern.txt")
                    #print "[%s:%s:]%s|%s:%s|%s:%s|%s"%(i,j,ball_a_red,ball_a_blue[0], ball_b_red, ball_b_blue[0], inter_r,inter_b)
                    balls_union = self.union(set(balls_union), set(inter_r))
                    balls_all = self.union(set(balls_all),ball_b_red)
                    ret[i] = ret.get(i,[])+[j]
            #print ball_a_red,balls_all,balls_union
        return ret

    def find_predict(self,idx,target):
        if not self.pattern_set:self.pattern_set = self.search_pattern_reverse(1770,0)
        #print self.pattern_set
        pre = self.intersection(self.pattern_set[idx],target)
        return idx,target,pre,len(target),len(pre)

    #从第N期开始搜索交集大于M个的期数（C）集合
    def search_union(self,n,c):
        ret,cnt ,ok = [], 0,[]
        comb = itto.combinations(range(n-1,0,-1),c)
        for item in comb:
            cnt += 1
            print item,len(ok)
            if not self.check_intersection(item):continue
            ok = self.do_search_match(item)
            if len(ok) <= 0: continue
            self.write2file(str(ok)+"\n","RET.txt")
            ret += ok
        return ret

    def search_union_combination(self,target,c):
        ret,cnt ,ok = [], 0,[]
        comb = itto.combinations(target,c)
        for item in comb:
            cnt += 1
            print item,len(ok)
            if not self.check_intersection(item):continue
            ok = self.do_search_match(item)
            if len(ok) <= 0: continue
            self.write2file(str(ok)+"\n","RET.txt")
            ret += ok
        return ret

    #从第N期开始搜索交集大于M个的期数（C）集合
    def search_union_recursive(self, n, m, c, path):
        path.append(n)
        if c==1:
            ret = self.do_search_match(path)
            if len(ret) > 0: self.write2file(str(ret)+"\n")
            print "search for:%s[%s]"%(path,ret)
            path.remove(n)
            return
        for i in range(n-1, 0, -1):
            if len(self.intersection(self.red_ball_row(n),self.red_ball_row(i))) < 2: continue
            self.search_union_recursive(i, m, c-1, path)
        path.remove(n)


    def do_search_match(self,path):
        candidate ,u, ret = [], [], []
        for item in path:
            red = self.red_ball_row(item)
            candidate.append(red)
            u = self.union(u,red)
        print u
        poss = itto.combinations(u,6)
        for p in poss:
            flag = True
            for can in candidate:
                #print p,can
                if len(self.intersection(p,can)) < 4:
                    flag = False
                    break
            if flag: ret.append((path,p))
        return ret

    def check_intersection(self,path):
        for i in path:
            first = self.red_ball_row(i)
            for j in path:
                if i >= j:continue
                if len(self.intersection(first,self.red_ball_row(j))) < 2:
                    return False
        return True




class PatternCycle(PatternGraph):

    THRESH = 4

    def __init__(self):
        PatternGraph.__init__(self)
        self.co_sim_tree = {}
        self.cycle_weight = [(0,100,1.0)]

    def get_cycle_weight(self, idx):
        for mi, ma, w in self.cycle_weight:
            if mi <= idx < ma:
                return w
        return 0

    def get_cycle_number(self, idx):
        for index, (mi, ma, w) in enumerate(self.cycle_weight):
            if mi <= idx < ma:
                return index
        return -1

    def rebalance_cycle_weight(self):
        for i, (mi, ma, w) in enumerate(self.cycle_weight):
            self.cycle_weight[i] = mi,ma,float("%.2f"%w)/2
        l = len(self.cycle_weight)
        self.cycle_weight.append((l*100,(l+1)*100,0.5))
        return l

    def shift_cycle_weight(self, children):
        #return
        if not children:return
        tmp = {}
        for child in children:
            c_num = self.get_cycle_number(child)
            if c_num == -1:
                c_num = self.rebalance_cycle_weight()
            tmp[c_num] = tmp.get(c_num, 0) + 1
        cnt = len(children)
        for i,(mi,ma,w) in enumerate(self.cycle_weight):
            self.cycle_weight[i] = mi, ma, (w + 0.1*float("%.2f"%tmp.get(i, 0))/cnt)/1.1
        return

    def create_co_sim_tree(self, idx):
        for i in range(0, idx):
            self.recreate_co_sim_tree(i)
        return

    def recreate_co_sim_tree(self, idx):
        return self.shift_tree(idx)

    def shift_tree(self, idx):
        children = self.find_similarity(idx)
        self.adjust_tree_root(idx, children)
        self.shift_cycle_weight(children)
        return self.co_sim_tree

    def find_similarity(self, idx):
        ret = []
        target = self.red_ball_row(idx)
        for i in range(idx-1, -1, -1):
            red = self.red_ball_row(i)
            if len(self.intersection(target, red)) >= self.THRESH: ret.append(i)
        return ret

    def adjust_tree_root(self, idx, children):
        self.co_sim_tree[idx] = (0, children)
        for child in children:
            cnt, cd = self.co_sim_tree.get(child, (0, []))
            self.co_sim_tree[child] = cnt + 1, cd
        return

    def get_predict_num(self, path):
        if len(path) == 1: return [path[0], path[0]],path,self.red_ball_row(path[0])

        max_step, min_step, inter = 0, float('inf'), self.red_ball_row(path[0])
        for idx in range(0,len(path)):
            if idx == 0: continue
            inter = self.intersection(inter, self.red_ball_row(path[idx]))
            tmp = int(path[idx-1]) - int(path[idx])
            max_step = tmp if tmp > max_step else max_step
            min_step = tmp if tmp < min_step else min_step
        #print [([path[0] + min_step , path[0] + max_step ], path, inter)]
        return [path[0] + min_step , path[0] + max_step ], path[:], inter

    def get_predict_num_2(self, idx,path):
        if len(path) == 1: return [path[0], path[0]],path,self.red_ball_row(path[0]),path[0]

        step = 1
        max_step, min_step, inter = 0, float('inf'), self.red_ball_row(path[0])
        for idx in range(0,len(path)):
            if idx == 0: continue
            inter = self.intersection(inter, self.red_ball_row(path[idx]))
            #tmp = int(path[idx-1]) - int(path[idx])

        avg = (int(path[0])-int(path[-1]))/len(path)
        #print path[-1],len(path)
        #print [([path[0] + min_step , path[0] + max_step ], path, inter)]
        return [path[0]+avg - step , path[0]+avg + step ], path[:], inter,path[0]

    def get_predict_num_3(self,idp , path):
        if len(path) == 1: return [path[0], path[0]],path,self.red_ball_row(path[0]),path[0]

        step = 0
        max_step, min_step, inter = 0, float('inf'), self.red_ball_row(path[0])
        for idx in range(0,len(path)):
            if idx == 0: continue
            inter = self.intersection(inter, self.red_ball_row(path[idx]))
            #tmp = int(path[idx-1]) - int(path[idx])
        avg = (int(path[0])-int(path[-1]))/len(path)
        i=1
        while True:
            max = path[0]+ i*avg + step
            min = path[0]+ i*avg - step
            if min <= idp <= max:
                return [min , max ], path[:], inter,path[0]
            if idp < min:break
            i += 1
        #print path[-1],len(path)
        #print [([path[0] + min_step , path[0] + max_step ], path, inter)]
        return [path[0]+avg - step , path[0]+avg + step ], path[:], inter,path[0]

    def get_co_sim_tree_children(self, idx):
        return self.co_sim_tree[idx][1]

    def do_search(self, idx, path):
        ret = []
        children = self.get_co_sim_tree_children(path[-1])
        if not children or len(path) >= 3:
            pre = self.get_predict_num_2(idx,path)
            #print pre
            #print int(pre[0][0]),idx,int(pre[0][1]) ,int(pre[0][0]) <= idx <= int(pre[0][1])
            return [pre] if int(pre[0][0]) <= idx <= int(pre[0][1]) else []
        for child in children:
            path.append(child)
            tmp = self.do_search(idx, path)
            path.remove(child)
            if len(tmp) > 0: ret += tmp
        return ret

    def search_idx_all(self,idx):
        ret = []
        for key,item in self.co_sim_tree.items():
            tmp = self.do_search(idx, [key])
            #print tmp
            if tmp: ret += tmp

        return ret

    def test_cycle(self, start):
        #if not self.pattern_set:
        #    self.pattern_set = self.search_pattern_reverse(len(self.data)-1, 0)

        self.create_co_sim_tree(start)
        #total =len(self.data)
        total = 1770
        for q in range(start,total):
            result = self.search_idx_all(q)
            ret = {}
            for p in result:
                pd, num = p[2],p[3]
                weight = self.get_cycle_weight(num)
                for item in pd:
                    ret[item] = ret.get(item, 0.0) + weight
            txt = self.sort_list_result(ret.items(),1,True)
            ok = []
            l = len(txt)
            if l>=28:
                l=28
            for i in range(0,l):
                ok.append(txt[i][0])
            #print self.cycle_weight
            red = self.red_ball_row(q)
            a = self.intersection(self.red_ball_row(q-1),self.red_ball_row(q-2))
            #b = self.intersection(self.red_ball_row(q-3),self.red_ball_row(q-2))
            print a

            ok =[i for i in ok if i not in a]
            #ok =[i for i in ok if i not in b]
            print red,len(ok),self.intersection(ok,red)
            self.write2file(str(txt)+"\n","TEST.txt")
            self.recreate_co_sim_tree(q)
        return

    def test_x(self):
        k,m = 0 , 0

        for i in range(17,1771):
            f = []
            for p in range(i-4+1,i):
                a = self.intersection(self.red_ball_row(p),self.red_ball_row(p-1))
                f = self.union(f,a)

            b = self.intersection(self.red_ball_row(i),f)
            print self.red_ball_row(i),b,f
            if not b:
                #print "xx",b
                k += 1
            else:
                m+=1
        print k,m



if __name__ == "__main__":
    num = 1700
    print time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
    np.set_printoptions(threshold='nan', linewidth=300)
    pc = PatternCycle()
    #pc.test_cycle(1768)
    pc.test_x()
    print time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))





