import {IonContent, IonHeader, IonPage, IonTitle, IonToolbar} from '@ionic/react';
import './Tab2.css';
import GoniometerComponent, {IGoniometerProps} from "../components/GoniometerComponent";
import React from "react";

const Tab2 = (props: IGoniometerProps) => {
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Goniometer</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">Goniometer</IonTitle>
          </IonToolbar>
        </IonHeader>
        <GoniometerComponent type={0} onFinish={props.onFinish}/>
      </IonContent>
    </IonPage>
  );
};

export default Tab2;
