<script lang="ts">
    import { onMount, onDestroy } from "svelte";

    let { lat, lng } = $props<{ lat: number; lng: number }>();
    console.log("Latitud:", lat, "Longitud:", lng);
    /** @type {HTMLDivElement | null} */
    let mapContainer: HTMLDivElement | null = null;

    /** @type {import("leaflet").Map | null} */
    let map: import("leaflet").Map | null = null;

    onMount(async () => {
        const L = (await import("leaflet")).default;

        if (!mapContainer) return;

        map = L.map(mapContainer).setView([lat, lng], 13);

        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: "© OpenStreetMap contributors",
        }).addTo(map);

        L.marker([lat, lng]).addTo(map).bindPopup("¡Estás aquí!").openPopup();
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
