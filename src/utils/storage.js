export function getStorage(key, defaultValue) {
	if (window.sessionStorage) {
		try {
			let val = window.sessionStorage.getItem(key);
			if (val) {
				val = JSON.parse(val);
				if (val) {
					return val;
				}
			}
		} catch (e) {
		}
		return defaultValue;
	}
}

export function setStorage(key, value) {
	if (window.sessionStorage) {
		window.sessionStorage.setItem(key, JSON.stringify(value));
	}
}
