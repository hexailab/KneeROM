import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar } from '@ionic/react';
import './Tab3.css';
import React from "react";

interface FullSpecifications {
  heightInches: number;
  weightPounds: number;
  truthLeftKnee: number;
  truthRightKnee: number;
  appLeftKnee: number;
  appRightKnee: number;
}

const Tab3 = (props: {specs: FullSpecifications}) => {
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Finish</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">Finish</IonTitle>
          </IonToolbar>
        </IonHeader>
        <p>Height in Inches: {props.specs.heightInches}</p>
        <p>Weight in Pounds: {props.specs.weightPounds}</p>
        <p>Left Knee Truth Value: {props.specs.truthLeftKnee}</p>
        <p>Right Knee Truth Value: {props.specs.truthRightKnee}</p>
        <p>Left Knee App Value: {props.specs.appLeftKnee}</p>
        <p>Right Knee App Value: {props.specs.appRightKnee}</p>

        <p>This isn't quite finished yet.</p>
      </IonContent>
    </IonPage>
  );
};

export default Tab3;
