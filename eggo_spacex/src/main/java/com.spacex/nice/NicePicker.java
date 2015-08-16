package com.spacex.nice;

import org.apache.commons.logging.Log;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
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

            stream.forEach(new Consumer<String>() {
                @Override
                public void accept(String s) {
                    PredictAPI pre = toPredicAPI(s);
                    if (pre == null)
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
            Date result =new SimpleDateFormat("yyyy-mm-dd HH:MM:SS").parse(date);
            Calendar limit = Calendar.getInstance();
            limit.add(Calendar.DAY_OF_MONTH, PRE_DATE);
            Date now = limit.getTime();
            return now.before(result);
        } catch (ParseException e) {
            e.printStackTrace();
            return false;
        }
    }

    public String getLuckyNumber(int date){
        if (date< 30&&date>0){
            PRE_DATE = -1 * date;
        }
        List<PredictAPI> pre = readResult();
        return new Gson().toJson(pre);
    }

    class PredictAPI{
        Integer seq;
        String start;
        String end;
        List<String> result;

        public PredictAPI(Integer seq, String start, String end, List<String> result) {
            this.seq = seq;
            this.start = start;
            this.end = end;
            this.result = result;
        }

        public List<String> getResult() {
            return result;
        }

        public void setResult(List<String> result) {
            this.result = result;
        }

        public Integer getSeq() {
            return seq;
        }

        public void setSeq(Integer seq) {
            this.seq = seq;
        }

        public String getStart() {
            return start;
        }

        public void setStart(String start) {
            this.start = start;
        }

        public String getEnd() {
            return end;
        }

        public void setEnd(String end) {
            this.end = end;
        }
    }

    public void main(String[] args){
        List<PredictAPI> res = new NicePicker().readResult();
        System.out.println(res.size());
    }

}
