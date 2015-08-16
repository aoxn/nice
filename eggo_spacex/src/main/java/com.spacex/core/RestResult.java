package com.spacex.core;

/**
 * Created by spacex on 15/6/20.
 */
public class RestResult {
    String result;
    Object message;

    public RestResult() {
    }

    public RestResult(String result, Object message) {
        this.result = result;
        this.message = message;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public Object getMessage() {
        return message;
    }

    public void setMessage(Object message) {
        this.message = message;
    }
}
