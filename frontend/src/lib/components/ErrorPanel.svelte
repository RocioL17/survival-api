<script lang="ts">
	let {
		title = 'HAS MUERTO',
		ctaLabel = '¿CONTINUAR? ↺',
		onContinue
	} = $props<{
		title?: string;
		ctaLabel?: string;
		onContinue?: () => void;
	}>();
</script>

<section class="error-shell" aria-live="assertive" aria-label="Mensaje de error">
	<div class="pixel pixel-1"></div>
	<div class="pixel pixel-2"></div>
	<div class="pixel pixel-3"></div>
	<div class="pixel pixel-4"></div>
	<div class="pixel pixel-5"></div>
	<div class="pixel pixel-6"></div>
	<div class="pixel pixel-7"></div>
	<div class="pixel pixel-8"></div>

	<div class="card">
		<div class="skull" aria-hidden="true">☠</div>
		<h2>{title}</h2>
		<button class="cta" type="button" onclick={onContinue}>{ctaLabel}</button>
	</div>
</section>

<style>
	.error-shell {
		position: fixed;
		inset: 0;
		display: grid;
		place-items: center;
		padding: 16px;
		background: #ef7777;
		overflow: hidden;
		z-index: 5001;
		animation: shellFade 180ms ease-out;
	}

	.card {
		position: relative;
		width: min(100%, 920px);
		background: #e7e7e7;
		border: 4px solid #3f4448;
		box-shadow: 8px 8px 0 rgba(63, 68, 72, 0.35);
		padding: clamp(24px, 3.5vw, 46px) clamp(20px, 3vw, 44px) clamp(22px, 2.8vw, 36px);
		text-align: center;
		animation: cardPop 260ms cubic-bezier(0.2, 0.8, 0.2, 1);
	}

	.skull {
		font-family: 'VT323', monospace;
		font-size: clamp(36px, 4vw, 58px);
		line-height: 1;
		color: #d83131;
		margin-bottom: 14px;
		animation: pulse 1.3s ease-in-out infinite;
	}

	h2 {
		margin: 0;
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(26px, 5vw, 64px);
		line-height: 1.15;
		letter-spacing: 1px;
		color: #2f363b;
	}

	.cta {
		margin-top: clamp(16px, 2.4vw, 28px);
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: min(100%, 420px);
		padding: 14px 24px;
		border: 0;
		background: #2f363b;
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(12px, 1.6vw, 20px);
		letter-spacing: 0.8px;
		color: #ffffff;
		cursor: pointer;
	}

	.pixel {
		position: absolute;
		width: 12px;
		height: 12px;
		background: rgba(206, 36, 36, 0.92);
		animation: drift 3.1s ease-in-out infinite;
	}

	.pixel-1 { top: 7%; left: 18%; }
	.pixel-2 { top: 11%; right: 16%; animation-delay: 260ms; }
	.pixel-3 { top: 24%; left: 10%; animation-delay: 120ms; }
	.pixel-4 { top: 31%; right: 7%; animation-delay: 390ms; }
	.pixel-5 { bottom: 16%; left: 12%; animation-delay: 520ms; }
	.pixel-6 { bottom: 23%; right: 22%; animation-delay: 200ms; }
	.pixel-7 { bottom: 9%; left: 34%; animation-delay: 480ms; }
	.pixel-8 { top: 15%; left: 62%; animation-delay: 320ms; }

	@keyframes shellFade {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	@keyframes cardPop {
		from {
			opacity: 0;
			transform: translateY(14px) scale(0.97);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	@keyframes pulse {
		0%, 100% { transform: scale(1); }
		50% { transform: scale(1.08); }
	}

	@keyframes drift {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-6px); }
	}

	@media (prefers-reduced-motion: reduce) {
		.error-shell,
		.card,
		.skull,
		.pixel {
			animation: none;
		}
	}

	@media (max-width: 640px) {
		.card {
			padding: 22px 14px 18px;
			box-shadow: 5px 5px 0 rgba(63, 68, 72, 0.25);
		}

		h2 {
			font-size: clamp(20px, 10vw, 34px);
		}

		.cta {
			min-width: 0;
			width: 100%;
			padding: 12px 14px;
		}
	}
</style>