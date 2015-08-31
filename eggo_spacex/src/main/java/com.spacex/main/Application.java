package com.spacex.main;

import com.spacex.core.VoiceWorker;
import com.spacex.nice.NiceWorker;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@ComponentScan(basePackages = {"com.spacex.core","com.spacex.web","com.spacex.nice"})
public class Application {

    public static void main(String[] args) {
        ConfigurableApplicationContext at = SpringApplication.run(Application.class, args);
        VoiceWorker worker = (VoiceWorker)at.getBean(VoiceWorker.class);

        // launch VoiceWorker thread
        new Thread(worker).start();

        // launch SSQ thread, scheduled every 4 hours
        new NiceWorker("python randpicker.py cross "+NiceWorker.crossTimes).start();
        new NiceWorker("python randpicker.py random "+ NiceWorker.randomTimes).start();
    }
}
