package com.spacex.Eggo;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;
import com.spacex.Eggo.PinService;

/**
 * Created by spacex on 15/6/20.
 */
public class BootOnReceiver extends BroadcastReceiver {

    @Override
    public void onReceive(Context context, Intent intent) {
        Log.d("System.err","onReceive..");
        //if (intent.getAction().equals("android.intent.action.BOOT_COMPLETED")){
            context.startService(new Intent(context, PinService.class));
        //}
        //
    }
}
