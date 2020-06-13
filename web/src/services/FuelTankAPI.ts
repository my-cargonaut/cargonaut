import Client from "./Client";
import { FuelTank } from "@/types";

class FuelTankAPI {
  async list(): Promise<FuelTank[]> {
    const response = await Client().get("/tanks");
    return response.data;
  }

  create(fuelTank: FuelTank) {
    return Client().post("/tanks", fuelTank);
  }

  update(fuelTank: FuelTank) {
    return Client().put("/tanks/" + fuelTank.id, fuelTank);
  }

  delete(id: string) {
    return Client().delete("/tanks/" + id);
  }
}

export default new FuelTankAPI();
