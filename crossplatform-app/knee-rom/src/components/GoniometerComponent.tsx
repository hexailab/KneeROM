import React from "react";
import "./GoniometerComponent.css";
import {PluginListenerHandle} from "@capacitor/core";
import {Motion, RotationRate} from "@capacitor/motion";
import {IonAlert, IonButton} from "@ionic/react";

export interface GoniometerFinishEvent {
  leftKneeRange: number,
  rightKneeRange: number,
}

export interface IGoniometerProps {
  type: number;
  onFinish: (ev: GoniometerFinishEvent) => void;
}

const GoniometerComponent = (props: IGoniometerProps) => {
  let minAngle: number = 360;
  let maxAngle: number = 0;

  const MINIMUM_BUFFER = 5;
  const GONIOMETER_START = "Start Goniometer";
  const GONIOMETER_END = "Stop Goniometer";

  let bufferWait = 0;
  const [debugCurrentEventObj, setDebugCurrentEventObj] = React.useState<RotationRate>({alpha: 0, beta: 0, gamma: 0})
  const [pitchRange, setPitchRange] = React.useState(0)
  const [isRunning, setRunning] = React.useState(false);
  const [motionHandler, setMotionHandler] = React.useState<PluginListenerHandle | undefined>(undefined);
  const [alertActive, setAlertActive] = React.useState(false);
  const [confirmAlert, setConfirmAlert] = React.useState(false);
  const [kneeNumber, setKneeNumber] = React.useState(0);
  const [leftKneePitch, setLeftKneePitch] = React.useState(0);


  const resetState = () => {
    setPitchRange(0);
    bufferWait = 0;
    minAngle = 360;
    maxAngle = 0;
  }

  const getKnee = () => {
    return kneeNumber == 0 ? "RIGHT" : "LEFT";
  }

  const onOrientationChange = (event: RotationRate) => {
    if (bufferWait++ > MINIMUM_BUFFER) {
      let angle;

      switch (props.type) {
        case 0: angle = (event.gamma > 90) ? 180 - event.beta : event.beta; break;
        case 1: angle = (event.gamma > 0) ? event.beta - 180 : event.beta; break;
        default: angle = -1;
      }

      maxAngle = Math.max(maxAngle, angle);
      minAngle = Math.min(minAngle, angle);

      setDebugCurrentEventObj(event);
      setPitchRange(maxAngle - minAngle);
    }
  }

  const registerAccelerationHandler = async () => {
    if (!isRunning) {
      resetState();
      if (motionHandler == undefined) setMotionHandler(await Motion.addListener('orientation', onOrientationChange));
    } else {
      if (motionHandler != undefined) await motionHandler.remove();
      setMotionHandler(undefined);
      setAlertActive(true);
    }

    setRunning(!isRunning);
  }

  const redoButton = async () => {
    resetState();
  }

  const confirmButton = async () => {
    if (kneeNumber == 0) {
      setLeftKneePitch(pitchRange);
      setKneeNumber(1);
      resetState()
      setConfirmAlert(true);
    } else {
      props.onFinish({
        leftKneeRange: Math.floor(leftKneePitch),
        rightKneeRange: Math.floor(pitchRange),
      })
    }
  }

  return (
    <div className={"container"}>
      <span style={{"fontSize": "14px", "color": "gray", "padding": "32px", "marginBottom": "16px", "textAlign": 'center'}}>Please secure your phone onto the front of your <b style={{"color": "black"}}>{getKnee()}</b> shin and press the "Start" button. Complete three full ranges of motion, and then stop the goniometer.</span>
      <div role="progressbar" style={({"--value":Math.floor(pitchRange)} as any) as React.CSSProperties}></div>
      <IonButton onClick={registerAccelerationHandler}>{isRunning ? GONIOMETER_END : GONIOMETER_START}</IonButton>
      <span style={{"fontSize": "10px", "color": "lightgray"}}>DEBUG: α{debugCurrentEventObj.alpha.toFixed(0)}, β{debugCurrentEventObj.beta.toFixed(0)}, λ{debugCurrentEventObj.gamma.toFixed(0)}</span>
      <IonAlert
        isOpen={alertActive}
        onDidDismiss={() => setAlertActive(false)}
        header={'Confirm Recording?'}
        message={`Affirm your given recording for the <b>${getKnee()}</b> knee: <b style="color: black;">${pitchRange.toFixed(0)}</b> degrees?`}
        buttons={[
          {
            text: 'Redo',
            role: 'redo',
            cssClass: 'secondary',
            handler: redoButton,
          },
          {
            text: 'Confirm',
            handler: confirmButton,
          }
        ]}
      />
      <IonAlert
        isOpen={confirmAlert}
        onDidDismiss={() => setConfirmAlert(false)}
        header={'Great!'}
        message={`Next, please switch your phone's position over to the <b>LEFT</b> shin and repeat the same steps.`}
        buttons={[
          {
            text: 'OK',
            handler: () => {},
          }
        ]}
      />
    </div>
  );
}

export default GoniometerComponent;