export interface Alert {
  type: string;
  message: string;
}

export interface FuelTank {
  id: string;
  truckId: string;
  name: string;
  length: number;
  width: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface Truck {
  id: string;
  manufacturer: string;
  model: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface RootState {
  fuelTanks: FuelTankState;
  trucks: TruckState;
}

export interface AlertState {
  alert: Alert;
}

export interface FuelTankState {
  fuelTanks: FuelTank[];
  loading: boolean;
}

export interface TruckState {
  trucks: Truck[];
  loading: boolean;
}
