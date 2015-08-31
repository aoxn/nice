package com.spacex.nice;

import com.spacex.core.PredictAPI;
import org.apache.commons.logging.Log;

import java.io.*;
import java.nio.file.Paths;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.List;
import java.util.function.Consumer;
import java.util.stream.Stream;
import com.google.gson.Gson;
import org.apache.commons.logging.LogFactory;
import org.springframework.stereotype.Service;

/**
 * Created by space on 2015/8/16.
 */
@Service
public class NicePicker {
    final Log log = LogFactory.getLog(getClass());
    private static int PRE_DATE = -4;
    private static final String PREDICT_FILE=
            Paths.get(".").toAbsolutePath().toString()+ File.separator+"result.txt";

    PredictAPI toPredicAPI(String content){
        PredictAPI res=null;
        try {
            res = new Gson().fromJson(content, PredictAPI.class);
        }catch (Exception e){
            log.debug("ParseAPI ERROR:" + e.getMessage());
        }
        return res;
    }

    List<PredictAPI> readResult(){
        List<PredictAPI> res = new ArrayList<PredictAPI>();
        try{
            Stream<String> stream =new BufferedReader(
                    new FileReader(PREDICT_FILE)).lines();
            log.info("PREDICT FILE: "+PREDICT_FILE);
            stream.forEach(new Consumer<String>() {
                @Override
                public void accept(String s) {
                    PredictAPI pre = toPredicAPI(s);
                    if (pre == null)
                        return;
                    if ("RANDOM_ROUND".equals(pre.getType().toUpperCase()))
                        return;
                    if (!dateFilter(pre.getStart()))
                        return;
                    res.add(pre);
                }
            });
        }catch (Exception e){
            e.printStackTrace();
        }
        return res;
    }


    boolean dateFilter(String date){
        try {
            Date result =new SimpleDateFormat("yyyy-MM-dd HH:mm:ss").parse(date);
            Calendar limit = Calendar.getInstance();
            limit.add(Calendar.DAY_OF_MONTH, PRE_DATE);
            Date now = limit.getTime();
            System.out.println(now.before(result));
            return now.before(result);
        } catch (ParseException e) {
            e.printStackTrace();
            return false;
        }
    }

    public List<PredictAPI> getLuckyNumber(int date){
        if (date< 30&&date>0){
            PRE_DATE = -1 * date;
        }
        List<PredictAPI> pre = readResult();
        log.info(pre);
        return pre;
        //return new Gson().toJson(pre);
    }

    public String doPrey(){
        String cmd = "python prey.py",tmp="";
        StringBuilder res = new StringBuilder("");
        BufferedReader b = null,r = null;
        try {
            Process p = Runtime.getRuntime().exec(cmd);
            int status =p.waitFor();
            b =new BufferedReader(new InputStreamReader(p.getErrorStream()));
            log.info("Exit status : "+status);
            if (status!=0){
                String line,msg="";
                while ((line = b.readLine()) != null)
                    msg+=line+"\n";
                log.error("Fail to CALL CMD: "+cmd +" "+msg);
                return "ERROR:"+msg;
            }
            log.debug("CMD execute finish...");
            r =new BufferedReader(new InputStreamReader(p.getErrorStream()));
            while ((tmp = r.readLine()) != null) {
                res.append(tmp);
            }
            return res.toString();
        }catch (InterruptedException|IOException ex) {
            ex.printStackTrace();
        }finally {
            close(b);
            close(r);
        }
        return "ERROR";
    }
    void close(Reader b ){
        if (b!=null){
            try {
                b.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
    public void main(String[] args){
        List<PredictAPI> res = new NicePicker().readResult();
        System.out.println(res.size());
    }

}
