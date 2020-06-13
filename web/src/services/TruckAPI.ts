import Client from "./Client";
import { Truck } from "@/types";

class TruckAPI {
  async list(): Promise<Truck[]> {
    const response = await Client().get("/trucks");
    return response.data;
  }

  create(truck: Truck) {
    return Client().post("/trucks", truck);
  }

  update(truck: Truck) {
    return Client().put("/trucks/" + truck.id, truck);
  }

  delete(id: string) {
    return Client().delete("/trucks/" + id);
  }
}

export default new TruckAPI();
