$(document).ready(function() {
    var mymap = L.map('mapid').setView([47.68, 9.178287], 13);

    // Add OSM tile layer
    var osm = new L.TileLayer(
        'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        {
            minZoom: 8,
            maxZoom: 18,
            attribution: 'Map data © <a href="http://openstreetmap.org">OpenStreetMap</a> contributors'
        }
    );
    mymap.addLayer(osm);

    var networks = L.markerClusterGroup();
    mymap.addLayer(networks);

    var baseLayers = {
        "OpenStreetMap": osm
    };

    var overlays = {
        "Networks": networks
    };

    L.control.layers(baseLayers, overlays).addTo(mymap);

    $.getJSON('query', null, function(data) {
        $.each(data, function(index, network) {
            var marker = L.marker([network.lat, network.lon], {
                title: network.ssid
            });
            marker.bindPopup(network.text);
            networks.addLayer(marker);
        });
    });
})
