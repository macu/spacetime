export function formatTimestamp(ms) {
 	let seconds = Math.floor((ms / 1000) % 60);
	let minutes = Math.floor((ms / (1000 * 60)) % 60);
	let hours = Math.floor((ms / (1000 * 60 * 60)) % 24);

	let timeString = '';
	if (hours > 0) {
		timeString += hours + 'h ';
	}
	if (minutes > 0) {
		timeString += minutes + 'm ';
	}
	if (seconds > 0) {
		timeString += seconds + 's';
	}

	return timeString;
}
