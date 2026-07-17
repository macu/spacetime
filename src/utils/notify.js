import {ElMessage} from 'element-plus';

export function alertSuccess(message) {
	ElMessage({
		message,
		type: 'success',
	});
}

export function showLoading(message) {
	return ElMessage({
		message,
		type: 'info',
		duration: 0,
		closable: false,
	});
}

export function showError(message) {
	ElMessage({
		message,
		type: 'error',
		showClose: true,
		duration: 10 * 1000,
	});
}
