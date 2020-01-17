<?php

$url = "https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson";

$today = date('Y-m-d');
$yesterday = date('Y-m-d', time() - 24 * 60 * 60);

$fullUrl = "$url&starttime=$yesterday&endtime=$today";

$data = file_get_contents($fullUrl);

if ($data === false) {
    echo "Error retrieving json\n";
    exit(1);
}

$json = json_decode($data);

foreach ($json->features as $record) {
    $prop = $record->properties;
    echo "{$prop->place} {$prop->mag}\n";
}
