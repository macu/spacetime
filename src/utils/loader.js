export const SCRIPTS = {
	DRAGULA: 'dragula',
};

let loadedScripts = {};

function loadAssets(assets) {
	let promises = [];
	for (var i = 0; i < assets.length; i++) {
		let a = assets[i];
		if (a.type === 'js' || /\.js$/.test(a.src)) {
			let script = document.createElement('script');
			script.type = 'text/javascript';
			script.crossorigin = 'anonymous';

			if (a.integrity) {
				script.integrity = a.integrity;
				script.crossOrigin = 'anonymous';
				script.referrerPolicy = 'no-referrer';
			}
			if (a.promise) {
				promises.push(a.promise);
			}
			promises.push(new Promise((resolve, reject) => {
				script.addEventListener('load', function () {
					resolve();
				});
				script.addEventListener('error', function () {
					reject();
				});
			}));
			script.src = a.src;
			document.body.appendChild(script);
		} else if (a.type === 'css' || /\.css$/.test(a.src)) {
			let link = document.createElement('link');
			link.rel = 'stylesheet';

			if (a.integrity) {
				link.integrity = a.integrity;
				link.crossOrigin = 'anonymous';
				link.referrerPolicy = 'no-referrer';
			}
			promises.push(new Promise((resolve, reject) => {
				link.addEventListener('load', function () {
					resolve();
				});
				link.addEventListener('error', function () {
					reject();
				});
			}));
			link.href = a.src;
			document.body.appendChild(link);
		}
	}
	// Return promise pending all assets
	return Promise.all(promises);
}

export function loadScript(scriptKey) {
	if (loadedScripts[scriptKey]) {
		return loadedScripts[scriptKey];
	}

	switch (scriptKey) {

		/*
		<script src="https://cdnjs.cloudflare.com/ajax/libs/dragula/3.6.6/dragula.min.js" integrity="sha512-MrA7WH8h42LMq8GWxQGmWjrtalBjrfIzCQ+i2EZA26cZ7OBiBd/Uct5S3NP9IBqKx5b+MMNH1PhzTsk6J9nPQQ==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/dragula/3.6.6/dragula.min.css" integrity="sha512-49xW99xceMN8dDoWaoCaXvuVMjnUctHv/jOlZxzFSMJYhqDZmSF/UnM6pLJjQu0YEBLSdO1DP0er6rUdm8/VqA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
		*/
		case SCRIPTS.DRAGULA:
			return loadedScripts[scriptKey] = loadAssets([
				{src: 'https://cdnjs.cloudflare.com/ajax/libs/dragula/3.6.6/dragula.min.js', type: 'js', integrity: 'sha512-MrA7WH8h42LMq8GWxQGmWjrtalBjrfIzCQ+i2EZA26cZ7OBiBd/Uct5S3NP9IBqKx5b+MMNH1PhzTsk6J9nPQQ=='},
				{src: 'https://cdnjs.cloudflare.com/ajax/libs/dragula/3.6.6/dragula.min.css', type: 'css', integrity: 'sha512-49xW99xceMN8dDoWaoCaXvuVMjnUctHv/jOlZxzFSMJYhqDZmSF/UnM6pLJjQu0YEBLSdO1DP0er6rUdm8/VqA=='},
			]).then(function() {
				console.debug('Dragula loaded');
				return window.dragula;
			});

	}
}
