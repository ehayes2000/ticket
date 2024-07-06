import { 
  type BaseEvent, 
  type Game, 
  type Concert,
  EventKind 
} from "@/models/Event";

export async function getMyEvents(): Promise<(Concert | Game)[]> { 
  return fetch("/api/getSavedEvents", {
    method: "GET", 
    credentials: "same-origin"
  }).then(response => {
    if(!response.ok){ 
      throw Error(`Bad response ${response.status}`)
    }
    return response.json()
  }).then(data => {
    if (!data) { 
      return [];
    } 
    const typedEvents: (Concert | Game)[] = data.map((e: any) => {
      const et: Concert | Game = e as Concert | Game;
      et.date = new Date(et.date)
      return et;
    })
    return typedEvents
  }).catch(e => { 
    console.error(e)
  })
}
