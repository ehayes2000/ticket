import { type BaseEvent, 
         type Game,
         type Concert,
         EventKind} from "@/models/Event";
import { reactive } from "vue";

export async function getEvents(): Promise<(Concert | Game)[]> {
  return fetch("/api/getEvents", {
    method: "GET"
  }).then(response => { 
    if (!response.ok){ 
      throw Error(`Bad response ${response.status}`)
    }
    return response.json()
  }).then(events => {
    if (!events){
      return [];
    }
    const typedEvents: (Concert | Game)[] = events.map(e => {
      const et: Concert | Game = e as Concert | Game;
      et.date = new Date(et.date)
      return et
    })
    return typedEvents;
  })
  .catch(e => {
    console.error("error getting events:", e)
  })
}
