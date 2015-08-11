package com.spacex.core;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;
import org.apache.lucene.analysis.Analyzer;
import org.apache.lucene.analysis.standard.StandardAnalyzer;
import org.apache.lucene.document.Document;
import org.apache.lucene.document.Field;
import org.apache.lucene.document.TextField;
import org.apache.lucene.index.DirectoryReader;
import org.apache.lucene.index.IndexReader;
import org.apache.lucene.index.IndexWriter;
import org.apache.lucene.index.IndexWriterConfig;
import org.apache.lucene.queryparser.classic.ParseException;
import org.apache.lucene.queryparser.classic.QueryParser;
import org.apache.lucene.search.IndexSearcher;
import org.apache.lucene.search.Query;
import org.apache.lucene.search.TopDocs;
import org.apache.lucene.store.Directory;
import org.apache.lucene.store.FSDirectory;
import org.springframework.stereotype.Service;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;

@Service
public class WareHourse {
    Logger log = Logger.getLogger(WareHourse.class);
    public final String BASE_PATH     = Paths.get(".").toAbsolutePath().toString();
    public final String INDEX_PATH    = BASE_PATH + File.separator+"data"+File.separator+"index"+File.separator;
    public final String RESOURCE_PATH = BASE_PATH + File.separator+"data"+File.separator+"resource"+File.separator;
    public final String WORKER_PATH   = BASE_PATH + File.separator+"data"+File.separator+"worker"+File.separator;
    public final String BACKUP_PATH   = BASE_PATH + File.separator+"data"+File.separator+"backup"+File.separator;

    public final String LOCATION_PATH = BASE_PATH + File.separator+"data"+File.separator+"location"+File.separator;


    private IndexWriter   indexer  = null;
    private IndexSearcher searcher = null;
    private QueryParser   parser   = null;
    private Analyzer      analyzer = null;
    private ObjectMapper  mapper   = new ObjectMapper();


    public WareHourse() {
        this.init(new StandardAnalyzer());
    }

    public WareHourse(Analyzer anlyzer) {
        init(anlyzer);
    }

    private void init(Analyzer analyzer) {
        log.debug("BASE_PATH: "+this.BASE_PATH);
        ensurePath();
        this.analyzer  = analyzer;
        this.initWriter();
        this.initSearcher(analyzer);
    }

    private void initWriter(){
        Directory         directory;
        IndexWriterConfig config;
        try {
            directory = FSDirectory.open(Paths.get(INDEX_PATH));
            config    = new IndexWriterConfig(analyzer);
            indexer   = new IndexWriter(directory, config);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    private void initSearcher(Analyzer analyzer){
        Directory         directory;
        IndexReader       reader;
        try {
            ensurePath();
            this.analyzer  = analyzer;
            directory = FSDirectory.open(Paths.get(INDEX_PATH));
            reader    = DirectoryReader.open(directory);
            searcher  = new IndexSearcher(reader);
            parser    = new QueryParser("words", this.analyzer);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    void ensurePath(){
        File index = new File(INDEX_PATH);
        if (!index.exists()){
            index.mkdirs();
        }
        File resource = new File(RESOURCE_PATH);
        if (!resource.exists()){
            resource.mkdirs();
        }
        File worker = new File(WORKER_PATH);
        if (!worker.exists()){
            worker.mkdirs();
        }
        File backup = new File(BACKUP_PATH);
        if (!backup.exists()){
            backup.mkdirs();
        }
        File locate = new File(LOCATION_PATH);
        if (!locate.exists()){
            locate.mkdirs();
        }
    }


    public void index(SoundETY sound) {
        this.doIndex(sound.getWords(),sound);
    }

    private void doIndex(String key, SoundETY sound) {
        Document doc = new Document();
        //这儿text可能为空，需要改进的地方
        doc.add(new TextField("words", key, Field.Store.YES));
        doc.add(new TextField("contents", sound.toString(), Field.Store.YES));
        try {
            indexer.addDocument(doc);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


    public byte[] getSpeakData(String words){
        SoundETY sound = doSearch(words);
        log.info("after do search");
        if(sound!=null){
            String ou_file = RESOURCE_PATH+File.separator+sound.outFileName();
            String in_file = RESOURCE_PATH+File.separator+sound.getFrom();

            File target = new File(ou_file);
            if (!target.exists()){
                try {
                    Runtime.getRuntime().exec(String.format("avconv -i %s -c copy -t %s -ss %s %s",
                            in_file,sound.getDuration(),sound.getStart(),ou_file)).waitFor();
                    log.info(String.format("avconv -i %s -c copy -t %s -ss %s %s",
                            in_file,sound.getDuration(),sound.getStart(),ou_file));

                } catch (IOException e) {
                    e.printStackTrace();
                    return new byte[0];
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            return readData(ou_file);
            // cut an
        }
        return new byte[0];
    }

    public byte[] readData(String file){
        try {
            return Files.readAllBytes(Paths.get(file));
        } catch (IOException e) {
            e.printStackTrace();
        }
        return new byte[0];
    }

    public SoundETY doSearch(String line) {
        SoundETY sound = null;
        try {
            if (searcher==null)
                initSearcher(this.analyzer);
            Query q     = parser.parse(line);
            TopDocs top = searcher.search(q, 20);
            if (top.totalHits<=0){
                return sound;
            }
            String r    = searcher.doc(top.scoreDocs[0].doc).getField("contents").stringValue();
            log.info("SEARCH RESULT => "+r);
            sound = mapper.readValue(r, SoundETY.class);
        } catch (IOException e) {
            e.printStackTrace();
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return sound;
    }

    public String getLastLocation(String cnt){
        if ("".equals(cnt)){
            cnt = "10";
        }
        String ret ="",tmp="";
        int c = Integer.parseInt(cnt);
        BufferedReader br = null;

        try{
            br =new BufferedReader(new FileReader(LOCATION_PATH+"locate.txt"));
            while((tmp=br.readLine())!=null)
                ret += tmp+",";
        }catch (Exception e){
            e.printStackTrace();
        }
        return ret;
    }

    public Boolean createLocation(String item){
        FileWriter wt = null;

        try{
            wt = new FileWriter(LOCATION_PATH+"locate.txt",true);
            wt.write(item+"\r\n");
            return true;
        }catch (IOException e) {
            e.printStackTrace();
        }finally {
            if (wt ==null){
                return true;
            }
            try {
                wt.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
        return false;
    }

    /*
     * 提交索引
     */
    public void commitAndClose() {
        if (indexer == null) {
            return;
        }
        try {
            indexer.commit();
            indexer.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    public void commit() {
        if (indexer == null) {
            return;
        }
        try {
            indexer.commit();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
