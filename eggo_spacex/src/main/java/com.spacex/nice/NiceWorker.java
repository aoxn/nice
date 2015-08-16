package com.spacex.nice;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import java.util.function.Consumer;
import java.util.stream.Stream;

/**
 * Created by space on 2015/8/16.
 */
public class NiceWorker{
    final Log log = LogFactory.getLog(getClass());
    long DELAY  = 5*1000;
    long PERIOD = 4*60*60*1000;
    Timer timer = new Timer();


    public void start(){

        timer.schedule(new TimerTask(){

            @Override
            public void run() {
                String cmd = "python randpicker.py";
                try {
                    log.debug("SSQ Thread: " + cmd + " FilePath:" + Paths.get(".").toAbsolutePath().toString());
                    Runtime.getRuntime().exec(cmd).waitFor();
                } catch (IOException e) {
                    e.printStackTrace();
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                log.debug("CMD execute finish...");
            }
        }, DELAY, PERIOD);

    }

}
