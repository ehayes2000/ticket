import { type TicketableEvent, 
         type Game,
         EventKind} from "@/models/Event";
import { reactive } from "vue";

export async function getEvents(): Promise<TicketableEvent[]> {
  return [
    {
      name: "Invasion",
      description: "Watch this saturday as godzilla rises from his watery grave to put an end to the capitalist greed of New York",
      venue: "Upper Bay",
      date: new Date("July 6, 2024"),
      kind: EventKind.Game,
      team1: "Godzilla",
      team2: "NYC"
    } as Game,
    {
      name: "Big wave",
      description: "large wave",
      venue: "MSG",
      date: new Date("July 6, 2024"),
      kind: EventKind.Concert,
      artist: "Godzilla"
    } as Concert
  ]
}

const eventStore = reactive({
  events: [] as TicketableEvent[]
  }
);

export function useEventStore(){
  return eventStore;
}