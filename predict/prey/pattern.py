#!/usr/bin/env python2
# -*- coding: utf-8 -*-

__author__ = 'spacex'

import random,json,time,sys
from ssq import Ssq
import numpy as np

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

class Matrix(Ssq):

    def __init__(self, n=0):
        Ssq.__init__(self)
        self.matrix = []
        if n == 0:
            n = len(self.data)
        for i in range(0,n,1):
            self.matrix.append(self.transform(self.red_ball_row(i)))
        self.matrix.append(np.ones(33))
        self.matrix = np.array(self.matrix)

    def transform(self,red):
        n = np.zeros(33)
        for i in red:
            n[i-1]=1
        return n

    def init(self):
        self.matrix = []
        for i in range(0,len(self.data),1):
            self.matrix.append(self.transform(self.red_ball_row(i)))
    def similarity(self, x, y):
        total,XSTEP,YSTEP = len(self.matrix),6 if x<30 else 3,6
        base  = self.matrix[total-YSTEP:total,x:x+XSTEP]
        varxy = self.matrix[y-YSTEP:y,x:x+XSTEP]
        ca = base - varxy
        sim = sum(sum(abs(ca[0:5])))
        return (y,sim,(x,y),varxy[5:6][0].tolist())

    def calcu(self):
        res = []
        le = len(self.matrix)
        for y in range(le - 2, 6, -1):
            ret = []
            for x in range(0, 33, 6):
                ret.append(self.similarity(x, y))
            res.append(ret)

        return res

    def do_sort(self,re ,n):
        return sorted(re, key=lambda item: item[n][1], reverse=False)

    def cut(self, re,n):
        ms = self.do_sort(re,n)
        tmp = ms[0:4]

        return [x[n][3] for x in tmp]
    def run(self):
        re = self.calcu()

        a0 = self.cut(re, 0)
        a1 = self.cut(re, 1)
        a2 = self.cut(re, 2)
        a3 = self.cut(re, 3)
        a4 = self.cut(re, 4)
        a5 = self.cut(re, 5)
        for x in range(len(a0)):
            print a0[x],a1[x],a2[x],a3[x],a4[x],a5[x]
        man = []
        for a in a0:
            for b in a1:
                for c in a2:
                    for d in a3:
                        for e in a4:
                            for f in a5:
                                a.append(b)
                                a.append(c)
                                a.append(d)
                                a.append(e)
                                a.append(f)
                                if sum(a)==6:
                                    man.append(a)
        print "LEN: ",len(man)
        print man



if __name__ == "__main__":
    rnd = Matrix(39)
    rnd.run()

