#!/usr/bin/env python2
# -*- coding: utf-8 -*-

__author__ = 'spacex'

import random,json,time,numpy as np
import itertools as itto
from ssq import Ssq

class WriteResult(object):
    def __init__(self,path,mode="a+"):
        self.path = path
        self.mode = mode

    def write(self, word,mode):
        with open(self.path,mode) as fd:
            fd.writelines(word)
        return

    def __call__(self, oper):
        def wrapper(wrap, *args, **kwargs):
            self.write(str(oper(wrap, *args, **kwargs))+"\r\n",self.mode)

        return wrapper

class RandomNum(Ssq):
    CROSS_LENGTH = 6
    def __init__(self):
        Ssq.__init__(self)

    def pick_up_ball(self, ln=6):
        ret = {}
        while len(ret.keys()) < ln:
            tmp = random.randint(1, 33)
            if ret.get(tmp, 0) != 0: continue
            ret[tmp] = tmp

        return [k for k,v in ret.items()]

    def random_round(self, times, ln =6):
        ret = []
        for i in range(0,times):
            rnd = self.pick_up_ball(ln)
            ret.append(rnd)
        return ret

    def cross_set(self,first, second):
        ret =[]
        for i in first:
            for j in second:
                if len(self.intersection(i,j))>=self.CROSS_LENGTH:
                    self.log.debug("Found Cross: %s,%s"%(i,j))
                    ret.append(str(i))
                    break
        return ret

    @WriteResult("result.txt","a+")
    def round_one(self,red,times,cnt=6):
        #random pick 7 ball N times, to see if whether the ball we wanted in it
        ret = []
        fi = self.random_round(times,cnt)
        red = self.red_ball_row(red)
        for i in fi:
            if len(self.intersection(i,red))>=self.CROSS_LENGTH:
                self.log.debug("Found: %s"%i)
                ret.append(str(i))
        return ret


    def round_two(self,times,cnt=6):
        # random pick 7 ball N times, we believe that our BINGO in it.
        # random pick 7 ball N times for twice, the cross set is we meant.

        return self.cross_set(self.random_round(times,cnt),
                                 self.random_round(times,cnt))

    @WriteResult("result.txt","a+")
    def run_nice(self,times,cnt=6):
        # API data struct
        #   seq : the sequence of Number
        # start : the start time of running this search
        #   end : the finish time of this search
        # result: the result of this search
        return json.dumps({
            "seq":len(self.data),
            "start":time.strftime("%y-%m-%d %H:%M:%S",time.gmtime()),
            "result":self.round_two(times,cnt),
            "end":time.strftime("%y-%m-%d %H:%M:%S",time.gmtime())
        })

if __name__ == "__main__":
    import sys
    rnd = RandomNum()
    times = 15000
    lines = 6
    arg = ""
    if len(sys.argv)<2:
        rnd.log.debug("Not enough argument %s"%sys.argv)
    cmd = sys.argv[1]
    if len(sys.argv) == 3:
        arg = sys.argv[2]
    if str(cmd).lower() == "start":
        if arg:
            try:
                times = int(arg)
            except:
                pass
        rnd.run_nice(times)
    else:
        rnd.log.debug("CMD error : %s %s"%(cmd,arg))
