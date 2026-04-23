<script lang="ts">
	import OptionButton from "$lib/components/OptionButton.svelte";
	import RetroPanel from "$lib/components/RetroPanel.svelte";
	import SuccessPanel from "$lib/components/SuccessPanel.svelte";
	import ErrorPanel from "$lib/components/ErrorPanel.svelte";
	import Map from "$lib/components/Map.svelte";
	import { onMount } from "svelte";

	let datos = $state([]);

	onMount(async () => {
		const res = await fetch("http://localhost:8080/case", {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		});
		console.log(res);
		datos = await res.json();
	});

	$effect(() => {
		console.log(datos);
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

	const options: Option[] = [
		{
			key: "A",
			text: "Huir! Correr hacia el callejon antes de que dispare.",
			dead: true,
		},
		{
			key: "B",
			text: "Hackear el sistema de sensores del dron rapidamente.",
			dead: false,
		},
		{
			key: "C",
			text: "Aceptar la inyeccion pacificamente y esperar lo mejor.",
			dead: false,
		},
	];

	let selectedOptionKey = $state("");
	let showSuccessPanel = $state(false);
	let showErrorPanel = $state(false);

	function handleOptionSelect(option: Option) {
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

				<h2>SUJETO 404</h2>

				<div class="identity-grid">
					<span>Edad:</span><span>32</span>
					<span>Sexo:</span><span>M</span>
					<span>Nac:</span><span>ARG</span>
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
					<Map lat={-34.3917} lng={-58.8731} />
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
				<p>
					El sol de la manana ilumina suavemente las ruinas del Sector
					7. Un dron de vigilancia se detiene bruscamente frente a ti,
					escaneando tus signos vitales.
				</p>
				<p>
					'CIUDADANO 404. ANOMALIA CARDIACA DETECTADA. REQUIERE
					ASISTENCIA INMEDIATA.'
				</p>
				<p>
					El dron prepara una inyeccion de sedante de alto impacto.
					Tus pulsaciones se disparan al ver la aguja.
				</p>
			</div>

			<div class="question">COMO PROCEDES?</div>

			<div class="options-wrap">
				{#each options as option}
					<OptionButton
						keyLabel={option.key}
						text={option.text}
						dead={option.dead}
						selected={selectedOptionKey === option.key}
						onSelect={() => handleOptionSelect(option)}
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
