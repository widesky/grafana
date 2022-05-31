import Map from 'ol/Map';
import TileLayer from 'ol/layer/Tile';
import OSM from 'ol/source/OSM';

import { MapLayerRegistryItem, MapLayerOptions } from '@grafana/data';

export const standard: MapLayerRegistryItem = {
  id: 'osm-standard',
  name: 'Open Street Map',
  isBaseMap: true,
  usesDataFrame: false,

  /**
   * Function that configures transformation and returns a transformer
   * @param options
   */
  create: async (map: Map, options: MapLayerOptions) => ({
    init: () => {
      return new TileLayer({
        source: new OSM(),
      });
    },
  }),
};

export const osmLayers = [standard];
