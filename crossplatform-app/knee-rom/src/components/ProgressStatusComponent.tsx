import "./ProgressStatusComponent.css"
import {pencil, compass, checkmarkCircleOutline} from "ionicons/icons";
import {IonIcon} from "@ionic/react";

interface IProgressStatusComponentProps {
  id: number;
}

const progressItems = [
  {title: "Information", icon: pencil},
  {title: "Goniometer", icon: compass},
  {title: "Finish", icon: checkmarkCircleOutline}
]

export default (props: IProgressStatusComponentProps) => {
  return <div className={"c-status-bar"}>
    {progressItems.map((ev, i) => {
      return (<div className={"c-status-bar__element" + ((i == props.id) ? " c-status-bar__element--active" : "")}>
        <IonIcon icon={ev.icon} className={"c-status-bar__icon"}/>
        <span className={"c-status-bar__title"}>{ev.title}</span>
      </div>)
    })}
  </div>
}