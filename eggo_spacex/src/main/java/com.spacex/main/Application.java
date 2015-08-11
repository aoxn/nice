package com.spacex.main;

import com.spacex.core.Worker;
import org.apache.catalina.core.ApplicationContext;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@ComponentScan(basePackages = {"com.spacex.core","com.spacex.web"})
public class Application {

    public static void main(String[] args) {
        ConfigurableApplicationContext at = SpringApplication.run(Application.class, args);
        Worker worker = (Worker)at.getBean(Worker.class);
        new Thread(worker).start();
    }
}