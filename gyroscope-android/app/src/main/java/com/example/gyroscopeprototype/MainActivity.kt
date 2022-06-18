package com.example.gyroscopeprototype

import android.content.Context
import android.hardware.Sensor
import android.hardware.SensorEvent
import android.hardware.SensorEventListener
import android.hardware.SensorManager
import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import kotlin.math.abs
import kotlin.math.floor

class MainActivity : AppCompatActivity(), SensorEventListener {
    private lateinit var sensorManager: SensorManager
    private var mRotation: Sensor? = null
    private var mMagnetic: Sensor? = null

    private var mGravity = FloatArray(3)
    private var mGeomagnetic = FloatArray(3)

    private var pitch = 0.0
    private var roll = 0.0

    private var minPitch = 0.0
    private var maxPitch = 0.0
    private var pitchRange = 0.0

    private var state = 0

    /*
     *
     * Eventual Needs:
     *   - Create internal GitHub server for easy collaboration. (Done)
     *   - Create AWS server for access. (Done)
     *   - Create REST API layer with MongoDB backend. (Done)
     *   - Individual sign-in for patient / doctor. TODO
     *   - Doctor/Patient pairing request. TODO
     *   - Interpret data, give recommendations & informatics TODO
     *          (i.e. line graph and projected recovery time)
     *
     * MORE RECENT TODO:
     *   - Create a slideshow which explains the backend/frontend connection.
     *   - Create screen-based flow diagram in Adobe XD, also for presenting.
     */

    fun toggleState() {
        if (state == 0) {
            // Reset default states
            pitch = 0.0
            minPitch = 180.0
            maxPitch = 0.0
            state = 1
            findViewById<Button>(R.id.refreshButton).text = "Stop"
        } else {
            state = 0
            findViewById<Button>(R.id.refreshButton).text = "Start"
        }
    }

    public override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        sensorManager = getSystemService(Context.SENSOR_SERVICE) as SensorManager
        mRotation = sensorManager.getDefaultSensor(Sensor.TYPE_ACCELEROMETER)
        mMagnetic = sensorManager.getDefaultSensor(Sensor.TYPE_MAGNETIC_FIELD)

        findViewById<Button>(R.id.refreshButton).text = "Start"
        findViewById<Button>(R.id.refreshButton).setOnClickListener {
            toggleState()
        }
    }

    override fun onAccuracyChanged(sensor: Sensor, accuracy: Int) {
        // Stub
    }

    override fun onSensorChanged(event: SensorEvent) {
        if (event.sensor.type === Sensor.TYPE_ACCELEROMETER) mGravity = event.values
        if (event.sensor.type === Sensor.TYPE_MAGNETIC_FIELD) mGeomagnetic = event.values

        if (mGravity != null && mGeomagnetic != null && state == 1) {
            val R = FloatArray(9)
            val I = FloatArray(9)
            val success = SensorManager.getRotationMatrix(R, I, mGravity, mGeomagnetic)

            if (success) {
                val orientation = FloatArray(3)
                SensorManager.getOrientation(R, orientation)
                pitch = Math.toDegrees(orientation[1].toDouble())
                roll = Math.toDegrees(orientation[2].toDouble())

                if (pitch < 0) pitch = 0.0
                if (abs(roll) > 90) pitch = 180.0 - pitch

                if (pitch < minPitch) { minPitch = pitch }
                if (pitch > maxPitch) { maxPitch = pitch }

                pitchRange = maxPitch - minPitch;
            }
        }
        findViewById<TextView>(R.id.pitch).text = floor(pitchRange).toString()
    }

    override fun onResume() {
        super.onResume()
        mRotation?.also { light ->
            sensorManager.registerListener(this, light, SensorManager.SENSOR_DELAY_UI)
        }
        mMagnetic?.also { obj ->
            sensorManager.registerListener(this, obj, SensorManager.SENSOR_DELAY_UI)
        }
    }

    override fun onPause() {
        super.onPause()
        sensorManager.unregisterListener(this)
    }
}