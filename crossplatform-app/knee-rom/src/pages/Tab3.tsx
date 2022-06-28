import {IonAlert, IonButton, IonContent, IonHeader, IonPage, IonTitle, IonToolbar} from '@ionic/react';
import {gql, useMutation} from '@apollo/client';
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

const CREATE_NEW_STUDY_ENTRY = gql`
    mutation CreateNewStudyEntry($rKF:Float!, $rKG:Float!, $lKF:Float!, $lKG:Float!, $pH:Float!, $pW: Float!) {
        createNewStudyEntry(
            rightKneeAppSideMeasure:0,
            rightKneeAppFrontMeasure:$rKF,
            rightKneeGoniometerTruth:$rKG,
            leftKneeAppSideMeasure:0,
            leftKneeAppFrontMeasure:$lKF,
            leftKneeGoniometerTruth:$lKG,
            patientHeight:$pH,
            patientWeight:$pW,
        ) {
            patientHeight,
            patientWeight,
            patientKneeLeft {
                appSideMeasure,
                appFrontMeasure,
                goniometerTruth,
            },
            patientKneeRight {
                appSideMeasure,
                appFrontMeasure,
                goniometerTruth
            }
        }
    }
`;

const Tab3 = (props: {specs: FullSpecifications, onFinish: () => void}) => {
  const [createEntry, entryData] = useMutation(CREATE_NEW_STUDY_ENTRY);
  const [finishAlert, setFinishAlert] = React.useState(false);
  const [errorAlert, setErrorAlert] = React.useState(false);

  const submit = async () => {
    await createEntry({
      variables: {
        rKF: props.specs.appRightKnee,
        lKF: props.specs.appLeftKnee,
        rKG: props.specs.truthRightKnee,
        lKG: props.specs.truthLeftKnee,
        pH: props.specs.heightInches,
        pW: props.specs.weightPounds,
      }
    });

    if (entryData.error != undefined) {
      setErrorAlert(true);
    } else {
      setFinishAlert(true);
    }
  }

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
        <h4 style={{marginLeft: "24px"}}>Your Results:</h4>
        <table>
          <tr>
            <td>Height in Inches</td>
            <td>{props.specs.heightInches}"</td>
          </tr>
          <tr>
            <td>Weight in Pounds</td>
            <td>{props.specs.weightPounds}lbs</td>
          </tr>
          <tr>
            <td>Left Knee Truth Value</td>
            <td>{props.specs.truthLeftKnee}°</td>
          </tr>
          <tr>
            <td>Right Knee Truth Value</td>
            <td>{props.specs.truthRightKnee}°</td>
          </tr>
          <tr>
            <td>Left Knee App Value</td>
            <td>{props.specs.appLeftKnee}°</td>
          </tr>
          <tr>
            <td>Right Knee App Value</td>
            <td>{props.specs.appRightKnee}°</td>
          </tr>
        </table>
        <IonButton style={{marginRight: "24px", float: "right"}} onClick={submit}>Submit Results</IonButton>
      </IonContent>
      <IonAlert
        isOpen={finishAlert}
        onDidDismiss={() => {
          setFinishAlert(false);
          props.onFinish();
        }}
        header={'Submitted!'}
        message={`Your information has been submitted to our servers. Thank you.`}
        buttons={[
          {
            text: 'Finish',
            handler: () => {},
          }
        ]}
      />
      <IonAlert
        isOpen={errorAlert}
        onDidDismiss={() => {
          setErrorAlert(false);
        }}
        header={'Error.'}
        message={`Sorry, there must have been some sort of issue. Please try again.`}
        buttons={[
          {
            text: 'Dismiss',
            handler: () => {},
          }
        ]}
      />
    </IonPage>
  );
};

export default Tab3;
