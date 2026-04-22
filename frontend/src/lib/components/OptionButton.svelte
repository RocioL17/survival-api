<script lang="ts">
	let {
		keyLabel,
		text,
		dead,
		selected = false,
		onSelect
	} = $props<{
		keyLabel: string;
		text: string;
		dead: boolean;
		selected?: boolean;
		onSelect?: () => void;
	}>();
</script>

<button
	class={`option ${selected ? (dead ? 'dead' : 'selected') : ''}`}
	type="button"
	aria-label={`Opcion ${keyLabel}: ${text}`}
	onclick={onSelect}
>
	<span class="key">{keyLabel}</span>
	<span class="option-text">{text}</span>
	<span class="caret">&lt;</span>
</button>

<style>
	.option {
		display: grid;
		grid-template-columns: 32px minmax(0, 1fr) auto;
		align-items: center;
		gap: 16px;
		width: 100%;
		min-height: 60px;
		padding: 12px;
		background: #fff;
		border: 2px solid #dfe6e9;
		border-radius: 6px;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
		font-family: 'VT323', monospace;
		font-size: clamp(18px, 2vw, 20px);
		line-height: 1;
		text-align: left;
		color: #2d3436;
		cursor: pointer;
		transition: transform 0.15s ease, border-color 0.15s ease;
	}

	.option:hover {
		border-color: #0984e3;
		transform: translateY(-1px);
	}

	.option.selected {
		border-color: #00b894;
		box-shadow: 0 0 0 2px rgba(0, 184, 148, 0.12), 0 1px 2px rgba(0, 0, 0, 0.1);
	}

	.option.dead {
		border-color: #d63031;
	}

	.option.dead .key {
		background: #d63031;
	}

	.key {
		display: grid;
		place-items: center;
		width: 32px;
		height: 32px;
		background: #0984e3;
		font-family: 'Press Start 2P', monospace;
		font-size: 10px;
		color: #fff;
	}

	.option-text {
		min-width: 0;
		white-space: normal;
		overflow: hidden;
		word-break: break-word;
	}

	.caret {
		opacity: 0.35;
		color: #0984e3;
		font-size: 24px;
		justify-self: end;
	}

	.option.dead .caret {
		color: #d63031;
	}

	@media (max-width: 600px) {
		.option {
			grid-template-columns: 28px minmax(0, 1fr);
			gap: 10px;
			align-items: start;
			min-height: 0;
		}

		.key {
			width: 28px;
			height: 28px;
		}

		.caret {
			grid-column: 2;
			justify-self: start;
			font-size: 20px;
			margin-top: -4px;
		}
	}
</style>
