type Venue = string; // maybe more complex later

export enum EventKind {
  Concert = "CONCERT",
  Game = "GAME"
}

export interface Tickets  { 
  event: (Game | Concert)
  seats: Array<string>
}

export interface BaseEvent { 
  id: number
  name: string
  description: string
  venue: Venue
  date: Date
  thumbnail?: HTMLImageElement,
  kind: EventKind
}

export interface Concert extends BaseEvent { 
  artist: string
}

export interface Game extends BaseEvent {  
  team1: string
  team2: string
}
