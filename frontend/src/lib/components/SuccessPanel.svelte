<script lang="ts">
	let {
		title = 'HAS SOBREVIVIDO!',
		ctaLabel = 'SIGUIENTE NIVEL →',
		onContinue
	} = $props<{
		title?: string;
		ctaLabel?: string;
		onContinue?: () => void;
	}>();
</script>

<section class="success-shell" aria-live="polite" aria-label="Mensaje de exito">
	<div class="sparkle sparkle-1"></div>
	<div class="sparkle sparkle-2"></div>
	<div class="sparkle sparkle-3"></div>
	<div class="sparkle sparkle-4"></div>
	<div class="sparkle sparkle-5"></div>
	<div class="sparkle sparkle-6"></div>
	<div class="card">
		<div class="star" aria-hidden="true">☆</div>
		<h2>{title}</h2>
		<button class="cta" type="button" onclick={onContinue}>{ctaLabel}</button>
	</div>
</section>

<style>
	.success-shell {
		position: fixed;
		inset: 0;
		display: grid;
		place-items: center;
		padding: 16px;
		background: #f8e79a;
		overflow: hidden;
		z-index: 5000;
		animation: shellFade 220ms ease-out;
	}

	.card {
		position: relative;
		width: min(100%, 920px);
		padding: clamp(28px, 4vw, 52px) clamp(20px, 3vw, 44px) clamp(24px, 3vw, 42px);
		background: #ffffff;
		border: 4px solid #3f4448;
		box-shadow: 8px 8px 0 rgba(63, 68, 72, 0.45);
		text-align: center;
		animation: cardPop 320ms cubic-bezier(0.2, 0.8, 0.2, 1);
	}

	.star {
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(26px, 4vw, 42px);
		line-height: 1;
		color: #f6c25b;
		margin-bottom: clamp(18px, 2vw, 28px);
		animation: starFloat 1.6s ease-in-out infinite;
		transform-origin: center;
	}

	h2 {
		margin: 0;
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(24px, 4.5vw, 48px);
		line-height: 1.25;
		letter-spacing: 1px;
		color: #3f4448;
	}

	.cta {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		margin-top: clamp(22px, 3vw, 36px);
		min-width: min(100%, 360px);
		padding: 14px 26px;
		border: 0;
		background: #3f4448;
		font-family: 'Press Start 2P', monospace;
		font-size: clamp(12px, 1.5vw, 18px);
		line-height: 1;
		letter-spacing: 1px;
		color: #ffffff;
		cursor: pointer;
	}

	.sparkle {
		position: absolute;
		width: 12px;
		height: 12px;
		background: rgba(255, 255, 255, 0.95);
		animation: sparkleDrift 2.8s ease-in-out infinite;
	}

	.sparkle-1 {
		top: 5%;
		left: 12%;
	}

	.sparkle-2 {
		top: 9%;
		right: 18%;
		animation-delay: 180ms;
	}

	.sparkle-3 {
		top: 18%;
		left: 34%;
		animation-delay: 320ms;
	}

	.sparkle-4 {
		top: 2%;
		left: 50%;
		width: 14px;
		height: 14px;
		animation-delay: 90ms;
	}

	.sparkle-5 {
		bottom: 12%;
		left: 8%;
		animation-delay: 260ms;
	}

	.sparkle-6 {
		bottom: 14%;
		right: 24%;
		width: 8px;
		height: 8px;
		animation-delay: 420ms;
	}

	@keyframes shellFade {
		from {
			opacity: 0;
		}

		to {
			opacity: 1;
		}
	}

	@keyframes cardPop {
		from {
			opacity: 0;
			transform: translateY(18px) scale(0.96);
		}

		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	@keyframes starFloat {
		0%,
		100% {
			transform: translateY(0) rotate(0deg) scale(1);
		}

		50% {
			transform: translateY(-4px) rotate(-6deg) scale(1.05);
		}
	}

	@keyframes sparkleDrift {
		0%,
		100% {
			transform: translateY(0) rotate(0deg) scale(1);
			opacity: 0.85;
		}

		50% {
			transform: translateY(-6px) rotate(12deg) scale(1.15);
			opacity: 1;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.success-shell,
		.card,
		.star,
		.sparkle {
			animation: none;
		}
	}

	@media (max-width: 640px) {
		.success-shell {
			padding: 12px;
		}

		.card {
			padding: 24px 16px 20px;
			box-shadow: 5px 5px 0 rgba(63, 68, 72, 0.35);
		}

		h2 {
			font-size: clamp(18px, 8vw, 32px);
		}

		.cta {
			min-width: 0;
			width: 100%;
			padding: 12px 18px;
		}
	}
</style>