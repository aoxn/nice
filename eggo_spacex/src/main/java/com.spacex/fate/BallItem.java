package com.spacex.fate;


import java.util.Comparator;

public class BallItem{
    Ball ball;
    Integer  currentCnt;
    Integer  historyCnt;
    String location;
    String preLoc;

    public BallItem() {
    }

    public BallItem(Ball ball, int currentCnt, int historyCnt, String location, String preLoc) {
        this.ball = ball;
        this.currentCnt = currentCnt;
        this.historyCnt = historyCnt;
        this.location = location;
        this.preLoc = preLoc;
    }

    public Ball getBall() {
        return ball;
    }

    public String getPreLoc() {
        return preLoc;
    }

    public void setPreLoc(String preLoc) {
        this.preLoc = preLoc;
    }

    public void setBall(Ball ball) {
        this.ball = ball;
    }

    public int getCurrentCnt() {
        return currentCnt;
    }

    public void setCurrentCnt(int currentCnt) {
        this.currentCnt = currentCnt;
    }

    public int getHistoryCnt() {
        return historyCnt;
    }

    public void setHistoryCnt(int historyCnt) {
        this.historyCnt = historyCnt;
    }

    public String getLocation() {
        return location;
    }

    public void setLocation(String location) {
        this.location = location;
    }

    @Override
    public String toString() {
        return "期号:"+String.format("%-10s",ball.getQc())+"号码:"+String.format("%-30s",ball.toString())+" 当前累积:"+String.format("%-5d",currentCnt)
                +" 历史累积:"+String.format("%-5d",historyCnt)+" 历史位置:"+String.format("%-40s",location)
                +" 预测累积:"+String.format("%-30s",preLoc);
    }

    static class RankComparator implements Comparator<BallItem> {

        @Override
        public int compare(BallItem b1, BallItem b2) {
            return b1.currentCnt.compareTo(b2.currentCnt);
//            if (){
//                return -1;
//            }
//            return 1;
        }
    }

}