package com.spacex.nice;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.io.IOException;
import java.util.Timer;
import java.util.TimerTask;

/**
 * Created by space on 2015/8/16.
 */
public class NiceWorker{
    final Log log = LogFactory.getLog(getClass());
    long DELAY  = 60*1000;
    long PERIOD = 4*60*60*1000;

    Timer timer = new Timer();


    public void start(){

        timer.schedule(new TimerTask(){

            @Override
            public void run() {
                String cmd = "python randpicker.py";
                try {
                    log.debug("SSQ Thread: " + cmd);
                    Runtime.getRuntime().exec(cmd);
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }, DELAY, PERIOD);

    }
}
