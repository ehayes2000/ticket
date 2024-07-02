type Venue = string; // maybe more complex later

export enum EventKind {
  Concert = "Concert",
  Game = "Game"
}

export interface TicketableEvent { 
  name: string
  description: string
  venue: Venue
  date: Date
  thumbnail?: HTMLImageElement,
  kind: EventKind
}

export interface Concert extends TicketableEvent { 
  artist: string
}

export interface Game extends TicketableEvent {  
  team1: string
  team2: string
}
