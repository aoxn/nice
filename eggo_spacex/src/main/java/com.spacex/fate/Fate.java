package com.spacex.fate;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;


public class Fate {
    Ssq ssq = new Ssq();
    List<BallItem> items = new ArrayList<>();
    BallItem.RankComparator cmp = new BallItem.RankComparator();

    public static void main(String[] args){
        new Fate().inspectPoint(1700);
    }
    private int cnt=0;
    public void inspectPoint(int num){

        if (num>=ssq.data.size()){
            System.out.println("Num Exceede");
            return;
        }
        for(int i = 0;i<num;i++){
            shift(ssq.data.get(i));
        }
        printBallItemArray();
        System.out.println("小于100次数:"+cnt);
    }

    private void shift(Ball ball){
        List<Integer> inter;
        String loc = "";
        String preLoc = "";
        for(int i =0;i<items.size();i++){
            BallItem it = items.get(i);
            inter = ball.intersection(it.ball);
            if(inter.size()>=4){
                it.currentCnt = 0;
//                loc += it.getBall().getQc()+",";
                loc += it.getBall().getRelative()+",";
                preLoc += i+",";
                if(i<=50)cnt++;
                break;
            }else {
                it.currentCnt ++;
                it.historyCnt ++;
            }
        }
        items.add(new BallItem(ball,0,0,loc,preLoc));
        //Collections.sort(items,cmp);
    }

    private void printBallItemArray(){
        for(BallItem it:items){
            System.out.println(it.toString());
        }
    }
}
