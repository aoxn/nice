package com.spacex.core;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class SoundETY {

    private String words;       //sentence
    private double   start;       //start time (seconds)
    private double   duration;    //duration time (seconds)
    private String from;        //come from which movie

    public SoundETY() {
    }

    public SoundETY(String words, double start, double duration, String from) {
        this.words = words;
        this.start = start;
        this.duration = duration;
        this.from = from;
    }

    @Override
    public String toString() {
        try {
            return new ObjectMapper().writeValueAsString(this);
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }
        return "Error Convert to Json";
    }
    @JsonIgnore
    public String outFileName(){
        return this.from.replace(".mp3","")+"-"+replaceSpecial(this.words)+"-"+this.start+"-"+this.duration+".mp3";
    }
    public String replaceSpecial(String words){
        return words.replaceAll("[\\pP‘’“”]", "").replaceAll("\\s*", "");
    }

    public String getWords() {
        return words;
    }

    public void setWords(String words) {
        this.words = words;
    }

    public double getStart() {
        return start;
    }

    public void setStart(double start) {
        this.start = start;
    }

    public double getDuration() {
        return duration;
    }

    public void setDuration(double duration) {
        this.duration = duration;
    }

    public String getFrom() {
        return from;
    }

    public void setFrom(String from) {
        this.from = from;
    }
}

