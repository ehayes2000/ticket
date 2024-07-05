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
    const typedEvents: (Concert | Game)[] = events.map(e => {
      const et: Concert | Game = e as Concert | Game;
      et.date = new Date(et.date)
      return et
    })
    return typedEvents;
  })
  .catch(e => {
    console.error(e)
  })
  // return [ 
  //   {
  //     name: "Invasion",
  //     description: "Watch this saturday as godzilla rises from his watery grave to put an end to the capitalist greed of New York",
  //     venue: "Upper Bay",
  //     date: new Date("July 6, 2024"),
  //     kind: EventKind.Game,
  //     team1: "Godzilla",
  //     team2: "NYC"
  //   } as Game,
  //   {
  //     name: "Big wave",
  //     description: "large wave",
  //     venue: "MSG",
  //     date: new Date("July 6, 2024"),
  //     kind: EventKind.Concert,
  //     artist: "Godzilla"
  //   } as Concert
  // ]
}

const eventStore = reactive({
  events: [] as BaseEvent[]
  }
);

export function useEventStore(){
  return eventStore;
}