import json
from ssq import Ssq

class Preyer(Ssq):
    """
        Get Result
    """

    RESULT_FILE = "result.txt"

    def __init__(self):
        Ssq.__init__(self)

    def load(self,):
        items = []
        d_len = len(self.data)
        with open(self.RESULT_FILE,"rw") as fi:
            for line in fi.readlines():
                try:
                    res = json.loads(line)
                except:
                    self.log.error("Json load Error.. %s"%res)
                curr = int(res.get("seq",-1))
                if curr >= d_len:
                    continue
                res['ball'] = self.red_ball_row(curr)
                res['cnt']  = self.count_number(res,res["ball"])
                items.append(res)
        return items

    def count_number(self,res_set,ball):
        t4,t5,t6 = 0, 0, 0
        for item in res_set.get("result",[]):
            t_len = len(self.intersection(eval(item),ball))
            t4 += 1 if t_len == 4 else 0
            t5 += 1 if t_len == 4 else 0
            t6 += 1 if t_len == 4 else 0
        return t4,t5,t6

if __name__=="__main__":
    rs = Preyer().load()
    print len(rs)
    for r in rs:
        print r["seq"],len(r["result"]),r["ball"],r["cnt"]
    print "DONE!"

