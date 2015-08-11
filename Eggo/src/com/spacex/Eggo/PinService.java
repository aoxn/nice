package com.spacex.Eggo;

import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.location.Criteria;
import android.location.Location;
import android.location.LocationManager;
import android.os.Binder;
import android.os.IBinder;
import android.util.Log;
import com.litesuits.http.LiteHttp;
import com.litesuits.http.exception.HttpException;
import com.litesuits.http.listener.HttpListener;
import com.litesuits.http.request.StringRequest;
import com.litesuits.http.response.Response;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;
import java.util.Timer;
import java.util.TimerTask;


public class PinService extends Service {
    SimpleDateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.CHINA);
    LocationManager  lm = null;
    Boolean ENABLE = true;
    String url = "http://221.228.86.118:8080/location";
//    String url = "http://192.168.1.103:8080/location";//
//    String url = "http://172.19.147.0:8080/location";
    Thread wth = null;




    Runnable worker = new Runnable() {
        LiteHttp liteHttp = LiteHttp.newApacheHttpClient(null);
        @Override
        public void run() {
            while(ENABLE){
                try {
                    do_work();
                    Thread.sleep(4*60*1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }catch (Exception e){
                    e.printStackTrace();
                    Log.e("PinService","Unexpected "+e.getMessage());
                }
            }
        }
        public void do_work() {
            String time = "",provider="";
            Location loc= null;
            try {
                time = df.format(new Date());
                provider = lm.getBestProvider(new Criteria(), true);
                loc = lm.getLastKnownLocation(provider);
            }catch (Exception e){
                time = "getLocation ERROR: "+e.getLocalizedMessage();
                loc  = new Location("Descrip");
            }
            final StringRequest request = new StringRequest(url)
                    .addUrlParam("time", time)
                    .addUrlParam("longitude", loc.getLongitude()+"")
                    .addUrlParam("latitude", loc.getLatitude()+"")
                    .addUrlParam("altitude", loc.getAltitude()+"")
                    .addUrlParam("accuracy",loc.getAccuracy()+"")
                    .setHttpListener(
                            new HttpListener<String>() {
                                @Override
                                public void onSuccess(String s, Response<String> response) {
                                    Log.d("EggoSpeak",String.format("send ok [%s]",""));
                                }

                                @Override
                                public void onFailure(HttpException e, Response<String> response) {
                                    Log.d("EggoSpeak","ops"+e.getMessage());
                                }
                            }
                    );

            // 1.1 execute async
            liteHttp.executeAsync(request);

        }
    };


    public void setEnable(Boolean enable){
        this.ENABLE = enable;

        if (ENABLE){
            wth = new Thread(worker);
            wth.start();
        }
    }

    @Override
    public IBinder onBind(Intent intent) {
        return pinBinder;
    }
    LocalPinBinder pinBinder = new LocalPinBinder();
    @Override
    public void onCreate() {
        super.onCreate();
        lm = (LocationManager)getSystemService(Context.LOCATION_SERVICE);
        wth = new Thread(worker);
        wth.start();
    }

    @Override
    public boolean onUnbind(Intent intent) {
        return super.onUnbind(intent);
    }

    @Override
    public void onDestroy() {
        ENABLE = false;
        super.onDestroy();
    }

    public class LocalPinBinder extends Binder{

        PinService getService(){
            return PinService.this;
        }
    }
}
