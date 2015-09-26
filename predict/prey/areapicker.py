#!/usr/bin/env python2
# -*- coding: utf-8 -*-

__author__ = 'spacex'

import random,json,time,sys
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

class AreaPicker(Ssq):

    AREA_LENGTH = 11

    def __init__(self):
        Ssq.__init__(self)

    def area_count(self, idx):
        i,j,k = 0,0,0
        red = self.red_ball_row(idx)
        for item in red:
            t = (int(item)-1)/AREA_LENGTH
            if t==0:
                i += 1
            else if t==1:
                j += 1
            else if t==2:
                k += 1
        return (i,j,k)

    def statistic(self, start, end):
        ret = {}
        for i in range(start,end):
            cnt = self.area_count(i)
            ret[str(cnt)] = ret.get(str(cnt),0) + 1
        return ret

    def weight(self,k,v):

        return 0

    def average_emerge(self, ret, total):
        ret = []
        for k,v in ret.items():
            ret.append([k,v,total/v,total/v - total%/v])
        return sorted(ret, key=lambda item: item[3], reverse=False)

    def print_area(self,res):
        for i in res:
            self.log.debug(str(i))

    def run(self, start, end):

        res = self.average_emerge(self.statistic(start,end),start-end)
        self.print_area(res)



if __name__ == "__main__":
    rnd = AreaPicker()
    rnd.run(0,1800)
