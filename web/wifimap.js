var NetworkIcon = L.Icon.extend({
    options: {
        iconSize: [32, 29],
        iconAnchor: [16, 28],
        popupAnchor: [0, -29]
    }
});

var openNetwork = new NetworkIcon({
    iconUrl: "icons/wireless-open.png"
});

var secureNetwork = new NetworkIcon({
    iconUrl: "icons/wireless-secure.png"
});

function icon(secure) {
    if (secure) {
        return secureNetwork;
    } else {
        return openNetwork;
    }
}

$(document).ready(function() {
    var mymap = L.map('mapid', {
        zoomControl: false
    }).setView([47.68, 9.178287], 13);

    // Add OSM tile layer
    var osm = new L.TileLayer(
        'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        {
            minZoom: 8,
            maxZoom: 19,
            attribution: 'Map data © <a href="http://openstreetmap.org">OpenStreetMap</a> contributors'
        }
    );
    mymap.addLayer(osm);

    var clusterGroup = L.markerClusterGroup({
        disableClusteringAtZoom: 17,
        spiderfyOnMaxZoom: false
    });
    mymap.addLayer(clusterGroup);

    var openNetworks = L.featureGroup.subGroup(clusterGroup, []);
    mymap.addLayer(openNetworks);

    var secureNetworks = L.featureGroup.subGroup(clusterGroup, []);
    mymap.addLayer(secureNetworks);

    var baseLayers = {
        "OpenStreetMap": osm
    };

    var overlays = {
        "Open networks": openNetworks,
        "Secure networks": secureNetworks
    };

    L.control.layers(baseLayers, overlays, {
        collapsed: false,
        position: "bottomright"
    }).addTo(mymap);

    L.control.scale().addTo(mymap);
    L.control.zoom({
        position: "bottomleft"
    }).addTo(mymap);

    var networkIcon = L.icon({
        iconUrl: "wireless-open.png",
        iconSize: [32, 29],
        iconAnchor: [16, 28],
        popupAnchor: [0, -29]
    });

    $.getJSON('query', null, function(data) {
        $.each(data, function(index, network) {
            var marker = L.marker([network.lat, network.lon], {
                title: network.ssid,
                icon: icon(network.secure)
            });
            marker.bindPopup(network.text);

            if (network.secure) {
                secureNetworks.addLayer(marker);
            } else {
                openNetworks.addLayer(marker);
            }
        });
    });
})
