package com.spacex.web;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

import com.spacex.core.PredictAPI;
import com.spacex.core.RestResult;
import com.spacex.core.VersionETY;
import com.spacex.core.WareHourse;
import com.spacex.nice.NicePicker;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class WordsController {

    private final Log log = LogFactory.getLog(getClass());
    private static final String template = "Hello, %s!";
    private final AtomicLong counter = new AtomicLong();

    @Autowired
    WareHourse wh = null;

    @Autowired
    NicePicker np = null;

    @RequestMapping("/version")
    public VersionETY version() {
        return new VersionETY("1.0.0","spacex","15521164411@163.com");
    }

    @RequestMapping("/location")
	public RestResult speak(@RequestParam(value="time", defaultValue="World") String time,
                                        @RequestParam(value="longitude", defaultValue="Longitude") String longitude,
                                        @RequestParam(value="latitude", defaultValue="Latitude") String latitude,
                                        @RequestParam(value="altitude", defaultValue="Altitude") String altitude,
                                        @RequestParam(value="accuracy", defaultValue="Accuracy") String accuracy) throws IOException {
	    String item = time+"|"+longitude+"|"+latitude+"|"+altitude+"|"+accuracy;
        log.info(item);
        if (!wh.createLocation(item)){
            return new RestResult("FALSE","ERROR OCCUR");
        }
	    return new RestResult("TRUE",item);
	}

    @RequestMapping("/result")
    public RestResult result(@RequestParam(value="last", defaultValue="10") String last) throws IOException {
        log.info(last);
        String ret = wh.getLastLocation(last);
        if (ret.equals("")){
            return new RestResult("FALSE","ERROR OCCUR");
        }
        return new RestResult("TRUE",ret);
    }

    @RequestMapping("/lucky")
    public RestResult ssqResult(@RequestParam(value="last", defaultValue="4") int last) throws IOException {
        log.info(last);
        List<PredictAPI> ret = np.getLuckyNumber(last);
        if (ret == null){
            return new RestResult("FALSE","ERROR OCCUR");
        }
        return new RestResult("TRUE",ret);
    }
}