// Runes-based stores. Module runs once per app load.

import { api } from "./api";

class VersionStore {
  version = $state("latest");
  build_sha = $state("—");
  build_time = $state("—");

  async load(): Promise<void> {
    try {
      const res = await api.getVersion();
      if (!res.ok) return;
      const data = await res.json();
      if (data?.version) this.version = data.version;
      if (data?.build_sha) this.build_sha = data.build_sha;
      if (data?.build_time) this.build_time = data.build_time;
    } catch {
      /* keep fallback */
    }
  }
}

export const versionStore = new VersionStore();
versionStore.load();
