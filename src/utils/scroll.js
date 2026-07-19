// Call when dragging starts, returns handler.
// Call handler.stop() when dragging stops.
export function startAutoscroll(maxScrollSpeed = 8, container = window) {
	const GUTTER_SIZE = 70; // distance from edge of viewport where scrolling starts

	const isWindow = container === window;

	let requestId = null;
	let clientY = null; // cursor position within viewport

	function handleMouseMove(e) {
		clientY = e.clientY;
	}

	window.addEventListener('mousemove', handleMouseMove);

	function handleScroll() {
		if (clientY !== null) {
			let viewportHeight = isWindow ? window.innerHeight : container.clientHeight, delta = 0;
			let clientYOffset = isWindow ? 0 : container.getBoundingClientRect().top;
			let clientYFromContainerTop = clientY - clientYOffset;
			if (clientYFromContainerTop < GUTTER_SIZE) { // Scroll up
				let factor = (GUTTER_SIZE - clientYFromContainerTop) / GUTTER_SIZE;
				delta = -((factor * maxScrollSpeed) + 1);
			} else if (clientYFromContainerTop > (viewportHeight - GUTTER_SIZE)) { // Scroll down
				let factor = (clientYFromContainerTop - (viewportHeight - GUTTER_SIZE)) / GUTTER_SIZE;
				delta = (factor * maxScrollSpeed) + 1;
			}
			if (delta !== 0) {
				if (isWindow) {
					window.scrollBy(0, delta);
				} else {
					container.scrollTop = container.scrollTop + delta;
				}
			}
		}
		requestId = window.requestAnimationFrame(handleScroll);
	}

	requestId = window.requestAnimationFrame(handleScroll);

	return {
		stop: function() {
			window.removeEventListener('mousemove', handleMouseMove);
			if (requestId) {
				window.cancelAnimationFrame(requestId);
				requestId = null;
			}
		},
	};
}
