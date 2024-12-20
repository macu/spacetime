let callbacks = {};

const bus = {
	on(event, callback) {
		callbacks[event] = callbacks[event] || [];
		callbacks[event].push(callback);
	},
	off(event, callback) {
		if (callbacks[event]) {
			callbacks[event] = callbacks[event].filter(cb => cb !== callback);
		}
	},
	emit(event, data) {
		if (callbacks[event]) {
			callbacks[event].forEach(cb => cb(data));
		}
	},
};

export default bus;
