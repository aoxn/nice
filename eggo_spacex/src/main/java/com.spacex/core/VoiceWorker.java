package com.spacex.core;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;
import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

@Service
public class VoiceWorker implements Runnable{

	private final Log log = LogFactory.getLog(getClass());
	//Logger log = Logger.getLogger(VoiceWorker.class);
	DateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss.S");
	String   BASE = "2015-06-14";
	@Autowired
	WareHourse wi = null;

	@Override
	public void run() {
		while(true){
			log.info(String.format("worker thread[%s] consumer()...",wi.BASE_PATH));
			consumer();
			try {
				Thread.sleep(10000);
			} catch (InterruptedException e) {
				log.error("VoiceWorker Thread sleep error , "+e.getLocalizedMessage());
				e.printStackTrace();
			}
		}
	}
	
	void consumer(){
		File[] files = new File(wi.WORKER_PATH).listFiles();
		if (files==null){return ;}
		for(File f:files){
			if (!f.getName().endsWith(".txt")){
				log.debug("Not a txt file continue!");
				continue;
			}
			doIndex(f);
		}
	}
	private String getFileExtension(File file) {
		String fileName = file.getName();
		if (fileName.lastIndexOf(".") != -1 &&
				fileName.lastIndexOf(".") != 0) {
			return fileName.substring(fileName.lastIndexOf(".") + 1);
		}
		return "";
	}

	String getAudioFormatPara(String line){
		String[] pas = line.split(",");
		for (String item:pas){
			String[] pa = item.split(":");
			if (pa.length != 2){
				return "";
			}
			if ("FORMAT" == pa[0].toUpperCase() ){
				return pa[1];
			}
		}
		return "";
	}

	void doIndex(File file){
		BufferedReader br = null;
		try{
			br = new BufferedReader(new FileReader(file));
			String tmp;
			int    cnt = 0;
			String from = "";
			while ((tmp = br.readLine())!=null){
				if (cnt++ == 0){
					//First line is metadata [Format:mp3,Meta:something]
					String format = getAudioFormatPara(tmp);
					if (format.equals("")){
						throw new Exception("Unkonw audio Format , please specify the audio" +
								" Format in the first line of input.txt like Format:mp3,Meta:something");
					}
					from = file.getName().replace(".txt","")+"."+format;
				}
				String rec[] = tmp.split(",");
				if (rec.length!=3){
					System.out.println("Length not 3");
					continue;
				}
				Date base  = df.parse(BASE+" "+"00:00:00.000");
				Date start = df.parse(BASE+" "+"0"+rec[0]+"0");
				Date end   = df.parse(BASE+" "+"0"+rec[1]+"0");
				double duration  = ((double)end.getTime() - (double)start.getTime())/1000;
				double first = ((double)start.getTime() - (double)base.getTime())/1000;

				//System.out.println(base.toString()+"  "+start.toString()+"  "+end.toString()+" "+duration+"  "+first);
				wi.index(new SoundETY(rec[2],first,duration,from));
			}
			wi.commit();
			File target = new File(wi.BACKUP_PATH+File.separator+file.getName()+".done");
			doRename(file,target);
			doRename(new File(wi.WORKER_PATH+File.separator+from),
					new File(wi.RESOURCE_PATH+File.separator+from));
		}catch (ParseException e) {
			e.printStackTrace();
		}catch(IOException ex){
			ex.printStackTrace();
		} catch (Exception e) {
			e.printStackTrace();
		} finally{
			close(br);
		}
		
	}
	void doRename(File file,File target){
		Boolean succ = true;
		if (target.exists()){
			succ = target.delete();
		}
		if (succ) {
			succ = file.renameTo(target);
			if (!succ)log.error(String.format("Rename File Error[src:%s,dst:%s]",file.getName(),target.getName()));
			if (file.exists())file.delete();
		}else{
			log.error(String.format("Remove File Error[src:%s,dst:%s]",file.getName(),target.getName()));
		}
	}

	void close(BufferedReader br){
		if (br==null){
			return;
		}
		try {
			br.close();
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

}
