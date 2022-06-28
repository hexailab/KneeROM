import {
  IonContent,
  IonHeader,
  IonItem,
  IonLabel,
  IonPage,
  IonTitle,
  IonToolbar,
  IonInput,
  IonButton,
} from '@ionic/react';
import React from "react";
import './Tab1.css';

interface InformationEventReturn {
  heightInches: number;
  weightPounds: number;
  truthValueLeft: number;
  truthValueRight: number;
}

interface InformationTabProps {
  onFinish: (ev: InformationEventReturn) => void;
}

const InformationTab = (props: InformationTabProps) => {
  const [heightInches, setHeightInches] = React.useState(0);
  const [weightPounds, setWeightPounds] = React.useState(0);
  const [truthValueLeft, setTruthValueLeft] = React.useState(0);
  const [truthValueRight, setTruthValueRight] = React.useState(0);


  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Information</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">Information</IonTitle>
          </IonToolbar>
        </IonHeader>
        <p style={{"padding": "24px"}}>Thank you for taking part in this evaluation. First, we would like some information:</p>
        <IonItem style={{"padding": "16px"}}>
          <IonLabel position="floating">Your Height in Inches</IonLabel>
          <IonInput type="number" onIonInput={(ev: any) => {setHeightInches(parseFloat(ev.target.value))}}></IonInput>
        </IonItem>
        <IonItem style={{"padding": "16px"}}>
          <IonLabel position="floating">Your Weight in Pounds</IonLabel>
          <IonInput type="number" onIonInput={(ev: any) => {setWeightPounds(parseFloat(ev.target.value))}}></IonInput>
        </IonItem>
        <IonItem style={{"padding": "16px"}}>
          <IonLabel position="floating">Actual Goniometer Value for Left Knee</IonLabel>
          <IonInput type="number" onIonInput={(ev: any) => {setTruthValueLeft(parseFloat(ev.target.value))}}></IonInput>
        </IonItem>
        <IonItem style={{"padding": "16px"}}>
          <IonLabel position="floating">Actual Goniometer Value for Right Knee</IonLabel>
          <IonInput type="number" onIonInput={(ev: any) => {setTruthValueRight(parseFloat(ev.target.value))}}></IonInput>
        </IonItem>
        <div style={{"padding": "16px", "float": "right"}}>
          <IonButton onClick={() => {props.onFinish({heightInches, weightPounds, truthValueLeft, truthValueRight})}}>Next</IonButton>
        </div>
      </IonContent>
    </IonPage>
  );
}

export default InformationTab;