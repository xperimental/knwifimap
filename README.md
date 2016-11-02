# knwifimap

This is a work-in-progress project for mapping the Wifis present in Konstanz, Germany. The database format used is compatible with the Wigle App which was used to generate the dataset.

For now this is just a simple viewer.

## Usage

You will need a SQLite database containing network data. Then you can run the following command to get a local viewer running:

```bash
knwifimap -a 127.0.0.1:8080 -f database.sqlite
```

## Database schema

The schema expected is currently:

```sql
CREATE TABLE network (
  bssid text primary key not null,
  ssid text not null,
  frequency int not null,
  capabilities text not null,
  lasttime long not null,
  lastlat double not null,
  lastlon double not null,
  type text not null default 'W',
  bestlevel integer not null default 0,
  bestlat double not null default 0,
  bestlon double not null default 0
);
```

## Libraries used

 - [Bootstrap](http://getbootstrap.com/)
 - [jQuery](http://jquery.com/)
 - [leaflet.js](http://leafletjs.com/)
 - [Leaflet.markercluster](https://github.com/Leaflet/Leaflet.markercluster)
 - [Leaflet.FeatureGroup.SubGroup](https://github.com/ghybs/Leaflet.FeatureGroup.SubGroup)
 - [leaflet-maskcanvas](https://github.com/domoritz/leaflet-maskcanvas)
