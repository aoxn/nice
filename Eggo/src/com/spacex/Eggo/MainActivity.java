package com.spacex.Eggo;

import android.app.Activity;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.ServiceConnection;
import android.media.MediaPlayer;
import android.media.MediaPlayer.*;
import android.os.Bundle;
import android.os.IBinder;
import android.util.Log;
import android.view.Window;
import com.iflytek.cloud.*;
import com.iflytek.cloud.util.ResourceUtil;
import com.litesuits.http.LiteHttp;
import com.litesuits.http.exception.HttpException;
import com.litesuits.http.listener.HttpListener;
import com.litesuits.http.request.StringRequest;
import com.litesuits.http.response.Response;
import com.spacex.core.ApiResult;
import com.spacex.core.JsonUtil;

import java.io.IOException;
import java.net.URLEncoder;

public class MainActivity extends Activity
{
    LiteHttp liteHttp      = null;
    SpeechRecognizer mIat  = null;
    SpeechSynthesizer mTts = null;
    VoiceWakeuper     mIvw = null;
    MediaPlayer         mp = new MediaPlayer();

//    final String serverUrl = "http://localhost:8080";
    final String serverUrl = "http://221.228.86.118:8080";
    PinService.LocalPinBinder pin= null;
    /** Called when the activity is first created. */
    @Override
    public void onCreate(Bundle savedInstanceState)
    {
        super.onCreate(savedInstanceState);
        this.requestWindowFeature(Window.FEATURE_NO_TITLE);
        setContentView(R.layout.main);

        mp.setOnCompletionListener(compListen);
        mp.setOnErrorListener(errListen);

        liteHttp = LiteHttp.newApacheHttpClient(null);
        SpeechUtility.createUtility(this, "appid=" + getString(R.string.app_id));
        //=================================================================================
        //1.创建 SpeechSynthesizer 对象, 第二个参数:本地合成时传 InitListener
        mTts= SpeechSynthesizer.createSynthesizer(this, null);
        //2.合成参数设置,详见《科大讯飞MSC API手册(Android)》SpeechSynthesizer 类
        // 设置发音人(更多在线发音人,用户可参见 附录12.2
        mTts.setParameter(SpeechConstant.VOICE_NAME, "xiaoqi");
        mTts.setParameter(SpeechConstant.SPEED, "50");//设置语速
        mTts.setParameter(SpeechConstant.VOLUME, "80");//设置音量,范围 0~100
        mTts.setParameter(SpeechConstant.ENGINE_TYPE, SpeechConstant.TYPE_CLOUD); //设置云端
        // 设置合成音频保存位置(可自定义保存位置),保存在“./sdcard/iflytek.pcm”
        //保存在 SD 卡需要在 AndroidManifest.xml 添加写 SD 卡权限
        // 仅支持保存为 pcm 格式,如果不需要保存合成音频,注释该行代码
        //mTts.setParameter(SpeechConstant.TTS_AUDIO_PATH, "./sdcard/iflytek.pcm");
        // 3.开始合成
        //mTts.startSpeaking("科大讯飞,让世界聆听我们的声音", mSynListener);
        //合成监听器

        //==================================================================================
        //1.加载唤醒词资源,resPath为唤醒资源路径
        StringBuffer param =new StringBuffer();
        String resPath = ResourceUtil.generateResourcePath(this, ResourceUtil.RESOURCE_TYPE.assets, "ivw/eggo.jet");
        param.append(ResourceUtil.IVW_RES_PATH+"="+resPath);
        param.append(","+ResourceUtil.ENGINE_START+"="+SpeechConstant.ENG_IVW);
        SpeechUtility.getUtility().setParameter(ResourceUtil.ENGINE_START,param.toString());
        //2.创建VoiceWakeuper对象
        mIvw = VoiceWakeuper.createWakeuper(this, null);
        //3.设置唤醒参数,详见《科大讯飞MSC API手册(Android)》SpeechConstant类
        //唤醒门限值,根据资源携带的唤醒词个数按照“id:门限;id:门限”的格式传入
        mIvw.setParameter(SpeechConstant.IVW_THRESHOLD,"0:"+"-20");
        //设置当前业务类型为唤醒
        mIvw.setParameter(SpeechConstant.IVW_SST,"wakeup");
        //设置唤醒一直保持,直到调用stopListening,传入0则完成一次唤醒后,会话立即结束(默认0)
        mIvw.setParameter(SpeechConstant.KEEP_ALIVE,"1");
        Log.d("SpeechResult:","Prepare to start wake");
        //4.开始唤醒
        //mIvw.startListening(mWakeuperListener);
        Log.d("SpeechResult:","start wake");

        //==================================================================================
        //1.创建SpeechRecognizer对象,第二个参数:本地听写时传InitListener
        mIat= SpeechRecognizer.createRecognizer(this, null);
        // 2.设置听写参数,详见《科大讯飞MSC API手册(Android)》SpeechConstant类
        mIat.setParameter(SpeechConstant.DOMAIN, "iat");
        mIat.setParameter(SpeechConstant.LANGUAGE, "zh_cn");
        mIat.setParameter(SpeechConstant.ACCENT, "mandarin ");
        //3.开始听写
        mIat.startListening(mRecoListener);

        //开启远程定位服务
        startService(new Intent(MainActivity.this,PinService.class));
        //绑定远程定位服务
        bindService(new Intent(MainActivity.this,PinService.class),sec, Context.BIND_AUTO_CREATE);
    }
    ServiceConnection sec = new ServiceConnection() {
        @Override
        public void onServiceConnected(ComponentName name, IBinder service) {
            Log.d("System.err","NAME");
            pin = (PinService.LocalPinBinder)service;
        }

        @Override
        public void onServiceDisconnected(ComponentName name) {
            Log.d("System.err","disconnectddddddddddd");
        }
    };

    public boolean shutdown = false;
    //听写监听器
    private WakeuperListener mWakeuperListener = new WakeuperListener() {

        public void onResult(WakeuperResult result) {
            Log.d("SpeechResult:","程序被唤醒，可以开始交流了。。。");
            mTts.startSpeaking("什么事?",mSynListener);
            String text = result.getResultString();
            Log.d("SpeechResult:",text);
            mIvw.stopListening();
            shutdown = false;
            mIat.startListening(mRecoListener);
        }
        public void onError(SpeechError error) {
            Log.d("SpeechResult:",error.toString());
        }
        public void onBeginOfSpeech() {}
        public void onEvent(int eventType, int arg1, int arg2, Bundle obj) {
            Log.d("SpeechResult:","OnEvent"+eventType);
            if (SpeechEvent.EVENT_IVW_RESULT == eventType) {
                //当使用唤醒+识别功能时获取识别结果 //arg1:是否最后一个结果,1:是,0:否。
                RecognizerResult reslut =((RecognizerResult)obj.get(SpeechEvent.KEY_EVENT_IVW_RESULT));
            }
        }

    };

    OnCompletionListener compListen = new OnCompletionListener() {
        @Override
        public void onCompletion(MediaPlayer mp) {
            if (mp!=null)
                mp.reset();
            listenAgain();
        }
    };
    OnErrorListener errListen = new OnErrorListener() {
        @Override
        public boolean onError(MediaPlayer mp, int what, int extra) {
            Log.d("EggoSpeak", "OnError fallback to call xunfei voice to speak. Reason=> " + what);
            //OnError fallback to call xunfei voice to speak
            if (mp != null)
                mp.reset();
            listenAgain();
            return false;
        }
    };
    HttpListener<String> turingRobot = new HttpListener<String>() {
        @Override
        public void onSuccess(String s, Response<String> response) {
            ApiResult api = JsonUtil.parseApiResult(s);
            Log.d("SpeechResult:","Eggo："+api.getText());
            //response.printInfo()
            try {
                mp.setDataSource(serverUrl+"/speak?words="+URLEncoder.encode(api.getText(),"utf-8"));
                Log.d("SpeechResult:","Eggo："+serverUrl+"/speak?words="+
                        URLEncoder.encode(api.getText(),"utf-8"));
                mp.prepare();
                mp.start();
            } catch (IOException e) {
                e.printStackTrace();
                Log.d("EggoSpeak: ", "Eggo remind Prepare or start Error=> "+ e.getLocalizedMessage());
                if (mp!=null) mp.reset();
                try {
                    mTts.startSpeaking(api.getText(), mSynListener);
                } catch (Exception ex) {
                    Log.e("SpeechResult:", "Error ()" + ex.getMessage());
                    listenAgain();
                }
            }
        }

        @Override
        public void onFailure(HttpException e, Response<String> response) {
            Log.e("SpeechResult:","http error []"+e.getMessage());
        }
    };

    //听写监听器
    private RecognizerListener mRecoListener = new RecognizerListener(){
         //听写结果回调接口(返回Json格式结果,用户可参见附录12.1);
        // 一般情况下会通过onResults接口多次返回结果,完整的识别内容是多次结果的累加;
        // 关于解析Json的代码可参见MscDemo中JsonParser类;
        //isLast等于true时会话结束。
        int silenceCnt = 0;
        String words = "";

        public void onResult(RecognizerResult results, boolean isLast) {
            words += JsonUtil.parseIatResult(results.getResultString());
            if (!isLast){
                return;
            }
            if(words.contains("关闭定位")){
                words = "";
                if (trackEnable(false))
                    mTts.startSpeaking("好的,已关闭！", mSynListener);
                else
                    mTts.startSpeaking("Sorry,关闭出错了！",mSynListener);
                return;
            }
            if(words.contains("打开定位")||words.contains("开启定位")){
                words = "";
                if (trackEnable(true))
                    mTts.startSpeaking("好的,已打开定位！", mSynListener);
                else
                    mTts.startSpeaking("Sorry,打开出错了！",mSynListener);
                return;
            }
            if (words.contains("闭嘴")){
                words = "";
                stopListener();
                mTts.startSpeaking("好的，我不说话就是了",mSynListener);
                return;
            }

            String url = "http://www.tuling123.com/openapi/api";
            final StringRequest request = new StringRequest(url)
                    .addUrlParam("key", "b10b30383938ebd5b880707b98661631")
                    .addUrlParam("info", words)
                    .addUrlParam("userid","1000")
                    .setHttpListener(turingRobot);
            // 1.1 execute async
            liteHttp.executeAsync(request);
            Log.d("SpeechResult:", "我："+words+"<>"+mIat.isListening());
            //Log.d("SpeechResult:", JsonParser.parseIatResult(results.getResultString()));
            words = "";

        }

        //会话发生错误回调接口
        public void onError(SpeechError error) {
            String desc = error.getPlainDescription(true);//获取错误码描述
            if (!mIat.isListening()&!shutdown) {
                if (silenceCnt++>10){
                    mIat.stopListening();
                    mIvw.startListening(mWakeuperListener);
                    silenceCnt = 0;
                    return;
                }
                Log.d("SpeechResult:","Erro 请开始说话。。。"+desc);
                mIat.startListening(mRecoListener);
            }
        }
        //开始录音
        public void onBeginOfSpeech () {
        }
        //音量值0~30
        public void onVolumeChanged(int volume) {
        }

        //结束录音
        public void onEndOfSpeech() {

        }

        //扩展用接口
        public void onEvent(int eventType, int arg1, int arg2, Bundle obj) {
        }
    };

    private void stopListener(){
        shutdown = true;
        if(mIat.isListening()) {
            mIat.stopListening();
        }
        if (!mIvw.isListening()){
            mIvw.startListening(mWakeuperListener);
        }
    }

    private Boolean trackEnable(Boolean enable){
        try{
            pin.getService().setEnable(enable);
        }catch (Exception e){
            e.printStackTrace();
            Log.e("SpeakResult","call pinservice error "+e.getMessage());
            return false;
        }
        return true;
    }


    private SynthesizerListener mSynListener = new SynthesizerListener() {
        //会话结束回调接口,没有错误时,error为null
        public void onCompleted(SpeechError error) {
            listenAgain();
        }

        @Override
        public void onEvent(int i, int i1, int i2, Bundle bundle) {

        }
        //缓冲进度回调
        // percent为缓冲进度0~100,beginPos为缓冲音频在文本中开始位置,
        // endPos表示缓冲音频在文本中结束位置,info为附加信息。

        public void onBufferProgress(int percent, int beginPos, int endPos, String info) {
        }

        //开始播放
        public void onSpeakBegin() {
        }

        //暂停播放
        public void onSpeakPaused() {
        }
        //播放进度回调
        // percent为播放进度0~100,beginPos为播放音频在文本中开始位置,endPos表示播放音频在文本中结束位置.

        public void onSpeakProgress(int percent, int beginPos, int endPos) {
        }
        //恢复播放回调接口

        public void onSpeakResumed() {
        }
        //会话事件回调接口
    };

    void listenAgain(){
        if (shutdown){
            return;
        }
        if (!mIat.isListening()) {

            Log.d("SpeechResult:","请开始说话。。。");
            mIat.startListening(mRecoListener);
        }else {
            Log.d("SpeechResult:","stop first");
            mIat.stopListening();
            mIat.startListening(mRecoListener);
        }
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        unbindService(sec);
    }
}
