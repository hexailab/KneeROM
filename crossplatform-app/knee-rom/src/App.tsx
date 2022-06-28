import { Redirect, Route, useLocation } from 'react-router-dom';
import {
  IonApp,
  IonIcon,
  IonLabel,
  IonRouterOutlet,
  IonTabBar,
  IonTabButton,
  IonTabs,
  setupIonicReact
} from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import { compass, pencil, checkmarkCircle } from 'ionicons/icons';
import Tab1 from './pages/Tab1';
import Tab2 from './pages/Tab2';
import Tab3 from './pages/Tab3';

/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Optional CSS utils that can be commented out */
import '@ionic/react/css/padding.css';
import '@ionic/react/css/float-elements.css';
import '@ionic/react/css/text-alignment.css';
import '@ionic/react/css/text-transformation.css';
import '@ionic/react/css/flex-utils.css';
import '@ionic/react/css/display.css';

/* Theme variables */
import './theme/variables.css';
import React from "react";
import ProgressStatusComponent from "./components/ProgressStatusComponent";
import {ApolloClient, ApolloProvider, InMemoryCache} from "@apollo/client";

setupIonicReact();

const App: React.FC = () => {
  const client = new ApolloClient({
    uri: 'http://54.146.48.48:8080/graphql',
    cache: new InMemoryCache()
  });

  const [heightInches, setHeightInches] = React.useState(0);
  const [weightPounds, setWeightPounds] = React.useState(0);
  const [truthLeftKnee, setTruthLeftKnee] = React.useState(0);
  const [truthRightKnee, setTruthRightKnee] = React.useState(0);
  const [appLeftKnee, setAppLeftKnee] = React.useState(0);
  const [appRightKnee, setAppRightKnee] = React.useState(0);
  const [currentTab, setCurrentTab] = React.useState(0);

  const resetApp = () => {
    setHeightInches(0);
    setWeightPounds(0);
    setTruthRightKnee(0);
    setTruthLeftKnee(0);
    setAppLeftKnee(0);
    setAppRightKnee(0);
    setCurrentTab(0);
  }

  return <IonApp>
    <ApolloProvider client={client}>
      {currentTab == 0 ?
        <Tab1 onFinish={ev => {
          setHeightInches(ev.heightInches);
          setWeightPounds(ev.weightPounds);
          setTruthLeftKnee(ev.truthValueLeft);
          setTruthRightKnee(ev.truthValueRight);
          setCurrentTab(1);
        }}/>
        : null
      }
      {currentTab == 1 ?
        <Tab2 type={0} onFinish={ev => {
          setAppLeftKnee(ev.leftKneeRange);
          setAppRightKnee(ev.rightKneeRange);
          setCurrentTab(2);
        }}/>
        : null
      }
      {currentTab == 2 ?
        <Tab3 specs={{heightInches, weightPounds, truthLeftKnee, truthRightKnee, appLeftKnee, appRightKnee}} onFinish={resetApp}/>
        : null
      }
      <ProgressStatusComponent id={currentTab}/>
    </ApolloProvider>
  </IonApp>;
};

export default App;
