package com.spacex.core;

import java.util.List;

/**
 * Created by space on 2015/8/16.
 */
public class PredictAPI {

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
