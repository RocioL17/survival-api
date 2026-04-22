<script lang="ts">
	let { label, filled = 8, total = 10, color = '#d63031', icon = '+' } = $props<{
		label: string;
		filled?: number;
		total?: number;
		color?: string;
		icon?: string;
	}>();
</script>

<div class="stat-row">
	<div class="icon-box">{icon}</div>
	<div class="label">{label}</div>
	<div class="bars" role="img" aria-label={`${label} ${filled} de ${total}`}>
		{#each Array(total) as _, i}
			<span class="bar" style={`--bar-color: ${i < filled ? color : '#ffffff'};`}></span>
		{/each}
	</div>
</div>

<style>
	.stat-row {
		display: grid;
		grid-template-columns: 32px minmax(64px, 80px) minmax(0, 1fr);
		align-items: center;
		gap: 8px;
	}

	.icon-box {
		display: grid;
		place-items: center;
		width: 32px;
		height: 32px;
		border: 2px solid #dfe6e9;
		border-radius: 6px;
		background: #fff;
		font-family: 'Press Start 2P', monospace;
		font-size: 9px;
		color: #2d3436;
	}

	.label {
		font-family: 'VT323', monospace;
		font-size: clamp(18px, 2vw, 20px);
		letter-spacing: 1px;
		text-transform: uppercase;
		color: #2d3436;
	}

	.bars {
		display: grid;
		grid-template-columns: repeat(10, minmax(0, 1fr));
		gap: 4px;
	}

	.bar {
		height: 16px;
		min-width: 0;
		border: 2px solid #2d3436;
		background: var(--bar-color);
	}

	@media (max-width: 600px) {
		.stat-row {
			grid-template-columns: 28px minmax(56px, 72px) minmax(0, 1fr);
			gap: 6px;
		}

		.icon-box {
			width: 28px;
			height: 28px;
		}

		.bars {
			gap: 3px;
		}

		.bar {
			height: 14px;
		}
	}
</style>
