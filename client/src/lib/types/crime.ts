export interface Crime {
    address: string;
    neighborhood: string;
    latitude: number;
    longitude: number;
    crime_type: string;
    date: string;
}

export interface CrimesResponse {
    crimes: Crime[];
    count: number;
}

export type Neighborhood =
    | 'Eastside'
    | 'New Tacoma'
    | 'South Tacoma'
    | 'Central Tacoma'
    | 'West End'
    | 'Northeast Tacoma'
    | 'South End'
    | 'North End'
    | 'Hilltop'
    | 'Stadium District';

export type CrimeType =
    | 'Traffic - DUI (Liquor)'
    | 'Destruction/Damage/Vandalism'
    | 'Larceny/Theft Offenses'
    | 'Assault Offenses'
    | 'Traffic Accident/Collision - Non Fatal - Injury'
    | 'Traffic Accident/Collision - Non Fatal - Non Injury'
    | 'Burglary/Breaking & Entering'
    | 'Robbery'
    | 'Motor Vehicle Theft'
    | 'Stolen Property Offenses'
    | 'Fraud Offenses'
    | 'Animal Cruelty'
    | 'Drug/Narcotics Violations';

export interface StrictCrime {
    address: string;
    neighborhood: Neighborhood;
    latitude: number;
    longitude: number;
    crime_type: CrimeType;
    date: string;
    time: string;
}

export interface StrictCrimesResponse {
    crimes: StrictCrime[];
    count: number;
}