package com.spacex.core;

import java.util.List;

/**
 * Created by space on 2015/8/16.
 */
public class PredictAPI {

    Integer seq;
    String start;
    String end;
    Integer length;
    String qc;
    String type;
    List<String> result;

    public PredictAPI(Integer seq, String start, String end, Integer length, String qc, String type, List<String> result) {
        this.seq = seq;
        this.start = start;
        this.end = end;
        this.length = length;
        this.qc = qc;
        this.type = type;
        this.result = result;
    }

    public Integer getLength() {
        return length;
    }

    public void setLength(Integer length) {
        this.length = length;
    }

    public String getQc() {
        return qc;
    }

    public void setQc(String qc) {
        this.qc = qc;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
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
