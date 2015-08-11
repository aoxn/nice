package com.spacex.fate;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Created by spacex on 15/6/25.
 */
public class Ball {

    List<Integer> red;
    int blue;
    int qc ;
    String date;
    int relative;
    public Ball() {
    }


    public Ball(List<Integer> red, int blue, int qc, String date,int relative) {
        Collections.sort(red);
        this.red = red;
        this.blue = blue;
        this.qc = qc;
        this.date = date;
        this.relative=relative;
    }

    public int getRelative() {
        return relative;
    }

    public void setRelative(int relative) {
        this.relative = relative;
    }

    public List<Integer> getRed() {
        return red;
    }

    public void setRed(List<Integer> red) {
        Collections.sort(red);
        this.red = red;
    }

    public int getQc() {
        return qc;
    }

    public void setQc(int qc) {
        this.qc = qc;
    }

    public String getDate() {
        return date;
    }

    public void setDate(String date) {
        this.date = date;
    }

    public int getBlue() {
        return blue;
    }

    public void setBlue(int blue) {
        this.blue = blue;
    }

    @Override
    public String toString() {
        return red.toString()+"|"+blue;
    }

    public List<Integer> intersection(Ball b){
        List<Integer> r_a = this.getRed(),r_b= b.getRed();
        int a_len = r_a.size() ,b_len = r_b.size(),i=0,j=0;
        List<Integer> res = new ArrayList<Integer>();
        int k=0;
        Integer f_a,f_b;
        while(i < a_len && j < b_len){
            k++;
            f_a = r_a.get(i);
            f_b = r_b.get(j);
            if (f_a.equals(f_b)){
                res.add(f_a);
                i++;
                j++;
            }else if (f_a > f_b){
                j++;
            }else {
                i++;
            }
        }
        //System.out.print(k);
        return res;
    }
}