package com.spacex.nice;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.io.*;
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
    String cmd  = "python randpicker.py random";

    public NiceWorker(String cmd){
        this.cmd = cmd;
    }

    public void start(){

        timer.schedule(new TimerTask(){
            void close(Reader b ){
                if (b!=null){
                    try {
                        b.close();
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                }
            }
            @Override
            public void run() {
                log.info("SSQ Thread: " + this.cmd + " FilePath:" + Paths.get(".").toAbsolutePath().toString());
                BufferedReader b = null;
                try {
                    Process p = Runtime.getRuntime().exec(this.cmd);
                    int status =p.waitFor();
                    b =new BufferedReader(new InputStreamReader(p.getErrorStream()));
                    log.info("Exit status : "+status);
                    if (status!=0){
                        String line,msg="";
                        while ((line = b.readLine()) != null)
                            msg+=line+"\n";
                        log.error("Fail to CALL CMD: "+this.cmd +" "+msg);
                    }
                }catch (InterruptedException|IOException ex) {
                    ex.printStackTrace();
                }finally {
                    close(b);
                }
                log.debug("CMD execute finish...");
            }
        }, DELAY, PERIOD);

    }

    public void main(String[] args){
        new NiceWorker().start();
    }
}
