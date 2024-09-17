export class szenario {
    AvailabilityAvg?: string
    AvailabilityCur?: string
    AvgTime?: string
    Img?: string
    IncidentCount?: Number
    IncidentList?: string
    LastTime?: string
    LastUpdate?: string
    Name?: string
    PromName?: string

    get availabilityCur(): Number {
        if (this.AvailabilityCur === undefined) {
            return -1
        }
        return parseFloat(this.AvailabilityCur)
    }
}

export interface overview {
    Version?: string

    PromURL?: string

    Szenarios?: szenario[]

}