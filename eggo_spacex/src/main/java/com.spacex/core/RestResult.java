package com.spacex.core;

/**
 * Created by spacex on 15/6/20.
 */
public class RestResult {
    String result;
    String message;

    public RestResult() {
    }

    public RestResult(String result, String message) {
        this.result = result;
        this.message = message;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
