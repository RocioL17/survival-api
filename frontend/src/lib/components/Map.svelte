<script lang="ts">
    import { onMount, onDestroy } from "svelte";

    let { lat, lng } = $props<{ lat: number; lng: number }>();

    let mapContainer: HTMLDivElement | null = null;
    let map = $state<import("leaflet").Map | null>(null);
    let marker = $state<import("leaflet").Marker | null>(null);

    onMount(async () => {
        const L = (await import("leaflet")).default;
        if (!mapContainer) return;

        map = L.map(mapContainer).setView([0, 0], 2);

        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: "© OpenStreetMap contributors",
        }).addTo(map);
    });

    $effect(() => {
        if (!map || !lat || !lng) return;

        map.setView([lat, lng], 13);

        if (marker) {
            marker.setLatLng([lat, lng]);
        } else {
            import("leaflet").then(({ default: L }) => {
                marker = L.marker([lat, lng]).addTo(map!).bindPopup("¡Estás aquí!").openPopup();
            });
        }
    });

    onDestroy(() => {
        if (map) map.remove();
    });
</script>

<svelte:head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet/dist/leaflet.css" />
</svelte:head>

<div bind:this={mapContainer} class="mapa"></div>

<style>
    .mapa {
        width: 100%;
        height: 400px;
    }
</style>
