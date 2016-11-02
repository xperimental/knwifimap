var NetworkIcon = L.Icon.extend({
    options: {
        iconSize: [32, 29],
        iconAnchor: [16, 28],
        popupAnchor: [0, -29]
    }
});

var hfOpenNetwork = new NetworkIcon({
    iconUrl: "icons/network-5-open.png"
});

var hfSecureNetwork = new NetworkIcon({
    iconUrl: "icons/network-5-secure.png"
});

var openNetwork = new NetworkIcon({
    iconUrl: "icons/network-2-open.png"
});

var secureNetwork = new NetworkIcon({
    iconUrl: "icons/network-2-secure.png"
});

function icon(highFreq, secure) {
    if (highFreq) {
        if (secure) {
            return hfSecureNetwork;
        } else {
            return hfOpenNetwork;
        }
    } else {
        if (secure) {
            return secureNetwork;
        } else {
            return openNetwork;
        }
    }
}

function download(openNetworks, secureNetworks, coverage) {
    $.getJSON('query', null, function(data) {
        var coverPoints = [];
        $.each(data, function(index, network) {
            var marker = L.marker([network.lat, network.lon], {
                title: network.ssid,
                icon: icon(network.highFreq, network.secure)
            });
            marker.bindPopup(network.text);

            if (network.secure) {
                secureNetworks.addLayer(marker);
            } else {
                openNetworks.addLayer(marker);
            }

            coverPoints.push([network.lat, network.lon]);
        });

        coverage.setData(coverPoints);
    });
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

    var coverage = L.TileLayer.maskCanvas({
        radius: 30,  // radius in pixels or in meters (see useAbsoluteRadius)
        useAbsoluteRadius: true,  // true: r in meters, false: r in pixels
        color: '#0F0',  // the color of the layer
        opacity: 0.5,  // opacity of the not covered area
        noMask: true,  // true results in normal (filled) circled, instead masked circles
    });
    mymap.addLayer(coverage);

    var baseLayers = {
        "OpenStreetMap": osm
    };

    var overlays = {
        "Open networks": openNetworks,
        "Secure networks": secureNetworks,
        "Coverage": coverage
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

    download(openNetworks, secureNetworks, coverage);
})
