<script lang="ts">
	import OptionButton from "$lib/components/OptionButton.svelte";
	import RetroPanel from "$lib/components/RetroPanel.svelte";
	import SuccessPanel from "$lib/components/SuccessPanel.svelte";
	import ErrorPanel from "$lib/components/ErrorPanel.svelte";
	import Map from "$lib/components/Map.svelte";
	import { onMount } from "svelte";

	type CaseData = {
		latitud?: number;
		longitud?: number;
		story?: {
			historia: string;
		};
		[key: string]: any;
		age?: number;
		gender?: string;
		zone?: string;
		accident?: string;
		name?: string;
	};

	let caseData: CaseData = $state({});
	let selectedIndex = $state<number | null>(null);
	let deadIndex = $state<number | null>(null);
	let loading = $state(false);
	let dataLoading = $state(true);

	onMount(async () => {
		try {
			const res = await fetch("http://localhost:8080/case", {
				method: "GET",
				headers: {
					"Content-Type": "application/json",
				},
			});
			caseData = await res.json();
			console.log("Datos recibidos:", caseData);
		} finally {
			dataLoading = false;
		}
	});

	$effect(() => {
		console.log("caseData actualizado:", caseData);
	});
	type Option = {
		key: string;
		text: string;
		dead: boolean;
	};

	type NotificationState = {
		type: "success" | "danger";
		message: string;
	};

	const options = $derived(
		caseData.choices?.map((text: string, i: number) => ({
			keyLabel: String.fromCharCode(65 + i),
			text,
			dead: deadIndex === i,
			selected: selectedIndex === i,
		})) ?? [],
	);

	async function selectOption(index: number) {
		if (loading || selectedIndex === index) return;
		selectedIndex = index;
		loading = true;
		// Simula una llamada a la API para verificar la opción seleccionada
		const res = await fetch("http://localhost:8080/options", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({choice: index}),
		});

		if (!res.ok) {
			console.error("Error del servidor:", res.status);
			loading = false;
			return;
		}

		const {value}  = await res.json();
		console.log(value);

		// if (!response.survived) {
		// 	deadIndex = index;
		// }

		loading = false;
	}

	let selectedOptionKey = $state("");
	let showSuccessPanel = $state(false);
	let showErrorPanel = $state(false);

	function handleOptionSelect(index: number) {
		const option = options[index];
		selectedOptionKey = option.key;
		showSuccessPanel = !option.dead;
		showErrorPanel = option.dead;
	}

	function handleContinue() {
		showSuccessPanel = false;
	}

	function handleRetry() {
		showErrorPanel = false;
		selectedOptionKey = "";
	}
</script>

{#if dataLoading}
	<div class="loading-screen">
		<div class="loading-box">
			<span class="loading-heart">♡</span>
			<p class="loading-text">CARGANDO...</p>
			<div class="loading-bar">
				<div class="loading-fill"></div>
			</div>
		</div>
	</div>
{/if}

{#if showSuccessPanel}
	<SuccessPanel onContinue={handleContinue} />
{/if}

{#if showErrorPanel}
	<ErrorPanel onContinue={handleRetry} />
{/if}

<main class="retro-app">
	<header class="game-title">
		<span class="heart">♡</span>
		<h1>VIVE O MUERE</h1>
		<span class="heart">♡</span>
	</header>

	<div class="game-layout">
		<section class="left-column">
			<RetroPanel
				title="Estado de Personaje"
				borderColor="#0984e3"
				tagWidth={224}
				className="character-panel"
			>
				<div class="avatar-wrap">
					<div class="pixel-avatar" aria-hidden="true"></div>
				</div>

				<h2>{caseData.name ?? "SUJETO 404"}</h2>

				<div class="identity-grid">
					<span>Edad:</span><span>{caseData.age ?? "N/A"}</span>
					<span>Sexo:</span><span>{caseData.gender ?? "N/A"}</span>
					<span>Zona:</span><span>{caseData.zone ?? "N/A"}</span>
				</div>
			</RetroPanel>

			<RetroPanel
				title="Mapa del Mundo"
				borderColor="#00b894"
				tagWidth={172}
				className="map-panel"
			>
				<div class="map-strip" aria-hidden="true"></div>
				<div class="location-box">
					<div class="location-label">UBICACION ACTUAL:</div>
					<Map
						lat={caseData.latitud ?? 0}
						lng={caseData.longitud ?? 0}
					/>
				</div>
			</RetroPanel>
		</section>

		<RetroPanel
			title="Historia"
			borderColor="#2d3436"
			tagWidth={108}
			className="story-panel"
		>
			<div class="story-text">
				<p>{caseData.story?.historia ?? "Cargando historia..."}</p>
				<!-- Aquí se mostraría la historia del caso -->
			</div>

			<div class="question">COMO PROCEDES?</div>

			<div class="options-wrap">
				{#each options as option, i}
					<OptionButton
						keyLabel={option.key}
						text={option.text}
						dead={option.dead}
						selected={option.selected}
						onSelect={() => selectOption(i)}
					/>
				{/each}
			</div>
		</RetroPanel>
	</div>
</main>

<style>
	.retro-app {
		width: min(100%, 1240px);
		height: 100dvh;
		margin: 0 auto;
		padding: clamp(8px, 1.2vw, 16px) clamp(10px, 1.6vw, 20px)
			clamp(10px, 1.6vw, 20px);
		display: grid;
		grid-template-rows: auto minmax(0, 1fr);
		gap: clamp(14px, 2vw, 24px);
		overflow: hidden;
	}

	.game-title {
		display: grid;
		grid-template-columns: auto 1fr auto;
		align-items: center;
		gap: clamp(6px, 0.8vw, 12px);
		width: min(100%, 470px);
		min-height: 56px;
		margin: 0 auto;
		padding: 0 clamp(12px, 1.6vw, 20px);
		border: 4px solid #2d3436;
		background: #fff;
		box-shadow: 6px 6px 0 #2d3436;
	}

	.game-title h1 {
		margin: 0;
		font-family: "Press Start 2P", monospace;
		font-size: clamp(14px, 1.5vw, 20px);
		letter-spacing: 2px;
		text-align: center;
		color: #2d3436;
	}

	.heart {
		font-family: "Press Start 2P", monospace;
		font-size: 14px;
		color: #ff7675;
	}

	.game-layout {
		display: grid;
		grid-template-columns: minmax(0, 360px) minmax(0, 1fr);
		gap: clamp(14px, 2vw, 24px);
		align-items: start;
		min-height: 0;
		padding-top: clamp(2px, 0.6vw, 8px);
	}

	.left-column {
		display: grid;
		gap: clamp(10px, 1.4vw, 18px);
		min-height: 0;
	}

	:global(.character-panel .panel-content) {
		display: grid;
		gap: clamp(8px, 1vw, 12px);
	}

	.avatar-wrap {
		display: grid;
		place-items: center;
	}

	.pixel-avatar {
		width: 84px;
		height: 84px;
		border: 4px solid #2d3436;
		background: radial-gradient(
				circle at 50% 16%,
				#fdcb6e 0 18%,
				transparent 19%
			),
			radial-gradient(circle at 34% 38%, #fdcb6e 0 14%, transparent 15%),
			radial-gradient(circle at 66% 38%, #fdcb6e 0 14%, transparent 15%),
			radial-gradient(circle at 50% 58%, #0984e3 0 11%, transparent 12%),
			#d6edff;
	}

	h2 {
		margin: 0;
		font-family: "Press Start 2P", monospace;
		font-size: clamp(14px, 1.5vw, 20px);
		line-height: 1.2;
		text-align: center;
		color: #0984e3;
	}

	.identity-grid {
		display: grid;
		grid-template-columns: 1fr auto;
		row-gap: 2px;
		font-family: "VT323", monospace;
		font-size: clamp(18px, 1.6vw, 22px);
		letter-spacing: 1px;
		text-transform: uppercase;
		color: #2d3436;
		padding: 0 8px;
	}

	:global(.map-panel .panel-content) {
		display: grid;
		gap: 6px;
	}

	.map-strip {
		height: 20px;
		border: 2px solid #2d3436;
		border-radius: 6px;
		background: repeating-linear-gradient(
				90deg,
				#74b9ff 0 15px,
				#55efc4 15px 30px,
				#00b894 30px 45px
			),
			repeating-linear-gradient(
				180deg,
				rgba(0, 0, 0, 0.1) 0 1px,
				transparent 1px 4px
			);
	}

	.location-box {
		border: 2px solid #2d3436;
		border-radius: 6px;
		background: #fff;
		padding: 10px 12px;
		text-align: center;
	}

	.location-label {
		font-family: "Press Start 2P", monospace;
		font-size: clamp(8px, 0.8vw, 10px);
		letter-spacing: 1px;
		color: #00b894;
		margin-bottom: 6px;
	}

	:global(.story-panel .panel-content) {
		display: grid;
		gap: clamp(10px, 1.4vw, 16px);
		min-height: 0;
	}

	.story-text {
		display: grid;
		gap: clamp(8px, 1vw, 12px);
		max-height: none;
		overflow: clip;
		padding-right: 8px;
	}

	.story-text p {
		margin: 0;
		font-family: "VT323", monospace;
		font-size: clamp(18px, 1.6vw, 28px);
		line-height: 1.1;
		color: #2d3436;
	}

	.question {
		padding: 10px 12px;
		border-left: 4px solid #0984e3;
		background: rgba(116, 185, 255, 0.1);
		font-family: "Press Start 2P", monospace;
		font-size: clamp(12px, 1.3vw, 16px);
		line-height: 1.35;
		color: #0984e3;
	}

	.options-wrap {
		display: grid;
		gap: 8px;
		padding-top: 10px;
		border-top: 4px dashed #dfe6e9;
	}

	@media (max-width: 1100px) {
		.game-layout {
			grid-template-columns: 1fr;
			gap: 14px;
		}

		.left-column {
			grid-template-columns: repeat(2, minmax(0, 1fr));
			align-items: start;
		}

		:global(.story-panel) {
			grid-column: 1 / -1;
		}

		:global(.story-panel .panel-content) {
			gap: 10px;
		}
	}

	@media (max-width: 820px) {
		.retro-app {
			padding: 8px 10px 10px;
		}

		.left-column {
			grid-template-columns: 1fr;
		}

		.game-title {
			max-width: 100%;
			box-shadow: 4px 4px 0 #2d3436;
		}

		:global(.panel-tag) {
			left: 12px;
			width: min(56vw, var(--tag-width));
		}

		:global(.retro-panel) {
			box-shadow: 4px 4px 0 rgba(45, 52, 54, 0.15);
		}

		.story-text p {
			font-size: clamp(16px, 4.2vw, 22px);
		}

		.question {
			font-size: clamp(11px, 3vw, 14px);
		}
	}

	.loading-screen {
		position: fixed;
		inset: 0;
		z-index: 100;
		background: #fff;
		display: grid;
		place-items: center;
	}

	.loading-box {
		display: grid;
		gap: 16px;
		place-items: center;
		padding: 32px 40px;
		border: 4px solid #2d3436;
		box-shadow: 8px 8px 0 #2d3436;
		background: #fff;
	}

	.loading-heart {
		font-family: "Press Start 2P", monospace;
		font-size: 28px;
		color: #ff7675;
		animation: heartbeat 0.8s ease-in-out infinite;
	}

	.loading-text {
		margin: 0;
		font-family: "Press Start 2P", monospace;
		font-size: clamp(12px, 1.5vw, 16px);
		letter-spacing: 3px;
		color: #2d3436;
	}

	.loading-bar {
		width: 200px;
		height: 16px;
		border: 3px solid #2d3436;
		background: #dfe6e9;
		overflow: hidden;
	}

	.loading-fill {
		height: 100%;
		background: #0984e3;
		animation: loading-progress 1.4s ease-in-out infinite;
	}

	@keyframes heartbeat {
		0%, 100% { transform: scale(1); }
		50% { transform: scale(1.3); }
	}

	@keyframes loading-progress {
		0% { width: 0%; }
		60% { width: 100%; }
		100% { width: 100%; }
	}

	@media (max-width: 600px) {
		.game-title {
			grid-template-columns: 1fr;
			min-height: 48px;
			padding-block: 10px;
			text-align: center;
		}

		.game-title h1 {
			font-size: 12px;
			line-height: 1.6;
		}

		.heart {
			display: none;
		}

		.avatar-wrap {
			padding-top: 0;
		}

		.pixel-avatar {
			width: 72px;
			height: 72px;
		}

		.location-box {
			padding: 8px 10px;
		}

		.options-wrap {
			padding-top: 8px;
		}

		:global(.panel-content) {
			padding: 12px;
		}
	}
</style>
