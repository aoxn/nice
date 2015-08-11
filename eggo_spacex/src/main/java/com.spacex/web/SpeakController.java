package com.spacex.web;

import com.spacex.core.WareHourse;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

@Controller
public class SpeakController {

    @Autowired
    WareHourse wh = null;

	private final Log log = LogFactory.getLog(getClass());
    @RequestMapping(value = "/speak",produces = "audio/mpeg")
	public ResponseEntity<byte[]> speak(@RequestParam(value="words", defaultValue="World") String words) throws IOException {
	    HttpHeaders headers = new HttpHeaders();
		log.info("SPEAK REVOKED!!!!! => "+words);
		headers.add("Content-Disposition","inline; filename=\"speak.mp3\"");
        byte[] sound = wh.getSpeakData(words);
	    return new ResponseEntity<byte[]>(sound,
	                                      headers, HttpStatus.OK);
	}
}