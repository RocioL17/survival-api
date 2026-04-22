<script lang="ts">
	let {
		title,
		borderColor = '#2d3436',
		tagWidth = 180,
		titleColor = '#ffffff',
		className = '',
		children
	} = $props<{
		title: string;
		borderColor?: string;
		tagWidth?: number;
		titleColor?: string;
		className?: string;
		children?: import('svelte').Snippet;
	}>();
</script>

<section
	class={`retro-panel ${className}`}
	style={`--panel-border: ${borderColor}; --tag-width: ${tagWidth}px; --tag-color: ${titleColor};`}
>
	<div class="panel-tag">{title}</div>
	<div class="panel-content">
		{@render children?.()}
	</div>
</section>

<style>
	.retro-panel {
		position: relative;
		border: 4px solid var(--panel-border);
		background: #fdfbf7;
		box-shadow: 8px 8px 0 rgba(45, 52, 54, 0.15);
	}

	.panel-tag {
		position: absolute;
		top: -16px;
		left: 16px;
		width: min(var(--tag-width), calc(100% - 32px));
		height: 24px;
		display: grid;
		place-items: center;
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(8px, 1vw, 10px);
		letter-spacing: 0.5px;
		text-transform: uppercase;
		background: #2d3436;
		color: var(--tag-color);
		text-align: center;
		padding-inline: 8px;
	}

	.panel-content {
		padding: clamp(16px, 2vw, 20px);
	}
</style>
