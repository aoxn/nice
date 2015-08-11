package com.spacex.fate;

import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.HttpClients;
import org.springframework.http.HttpMethod;
import org.springframework.util.FileCopyUtils;

import java.io.*;
import java.util.*;

/**
 * Created by spacex on 15/6/24.
 */
public class Ssq {
    public static void main(String[] args){
//        Ssq s = new Ssq();
//        List<Integer> red_a = new ArrayList<Integer>(6);
//        red_a.add(12);
//        red_a.add(3);
//        red_a.add(8);
//        red_a.add(1);
//        red_a.add(31);
//        red_a.add(22);
//        Ball a = new Ssq.Ball(red_a,6);
//        List<Integer> red_b = new ArrayList<Integer>(6);
//        red_b.add(1);
//        red_b.add(20);
//        red_b.add(18);
//        red_b.add(13);
//        red_b.add(6);
//        red_b.add(12);
//        Ball b = new Ssq.Ball(red_b,6);
//        System.out.println(a.getRed());
//        System.out.println(b.getRed());
//        List<Integer> res = s.intersection(a,b);
//        System.out.println(res);
        new Ssq().init();
    }
    List<Ball> data = new ArrayList<>();

    public Ssq() {
        init();
    }

    private void loadFile(){
        File ssq = new File("mp.txt");
        //System.out.println(ssq.getAbsoluteFile());
        if (ssq.exists()){
            return;
        }
        //Get File and Parse
        HttpClient httpclient = HttpClients.createDefault();
        HttpResponse response;
        OutputStream out      = null;
        try {
            response = httpclient.execute(new HttpGet("http://www.17500.cn/getData/ssq.TXT"));
            HttpEntity entity = response.getEntity();
            if (entity != null) {
                out = new FileOutputStream(ssq);
                FileCopyUtils.copy(entity.getContent(),out);
            }
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            if (out!=null)
                try {
                    out.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
        }
    }

    private void pasreFile(){
        File ssq = new File("mp.txt");
        BufferedReader br = null;
        String tmp = "";
        String[] p = null;
        try{
            br = new BufferedReader(new FileReader(ssq));
            int i = 0;
            while((tmp=br.readLine())!=null) {
                p = tmp.split("\\s+");
                //System.out.println(p[2]);
                List<Integer> list = new ArrayList<>();
                list.add(Integer.parseInt(p[2]));
                list.add(Integer.parseInt(p[3]));
                list.add(Integer.parseInt(p[4]));
                list.add(Integer.parseInt(p[5]));
                list.add(Integer.parseInt(p[6]));
                list.add(Integer.parseInt(p[7]));
                Ball b = new Ball(list,Integer.parseInt(p[8]),Integer.parseInt(p[0]),p[1],i++);
                data.add(b);
            }
        }catch (Exception e){
            e.printStackTrace();
        }
    }

    private void init(){
        loadFile();
        pasreFile();

    }


    public List<Integer> intersection(Ball a,Ball b){
        List<Integer> r_a = a.getRed(),r_b= b.getRed();
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

    public int[] difference(Ball a,Ball b){

        return new int[0];
    }

    public int[] union(Ball a,Ball b){

        return new int[0];
    }





}
