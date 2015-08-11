package com.spacex.main;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * Created by spacex on 15/6/18.
 */
public class Test {
    public void main(String[] args){
        //readData("/Users/spacex/work/dp/android/eggo_spacex/./data/resource/哆啦A梦伴我同行-我已经叫你好几遍了-93-2.mp3");
        DateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss.S");
        try {
            Date d = df.parse("2014-05-06 00:00:30.240");
            Date e = df.parse("2014-05-06 00:00:31.300");
            double m =((double)e.getTime()-(double)d.getTime());
            double k = m/1000;
            System.out.println((double)e.getTime()+" "+e.getTime()+" "+m+" "+k+"  "+d.getTime());
        } catch (ParseException e) {
            e.printStackTrace();
        }
    }
    public static byte[] readData(String file){
        try {
            return Files.readAllBytes(Paths.get(file));
        } catch (IOException e) {
            e.printStackTrace();
        }
        return new byte[0];
    }
}
