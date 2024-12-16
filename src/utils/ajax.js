import axios from 'axios';

import {
	ElMessageBox,
} from 'element-plus';

export const ajax = axios.create({
	headers: {
		'Content-Type': 'multipart/form-data',
	},
});

export default ajax;

export function ajaxPost(url, params = {}, errorCodeMessages = {}) {
	let formData = new FormData();
	for (let key in params) {
		if (params[key] !== null && params[key] !== undefined) {
			formData.append(key, params[key]);
		}
	}
	return ajax.post(url, formData).then(response => {
		return response.data;
	}).catch(error => {
		alertError(error, errorCodeMessages);
		throw error;
	});
}

export function ajaxGet(url, params = {}, errorCodeMessages = {}) {
	return ajax.get(url, {params}).then(response => {
		return response.data;
	}).catch(error => {
		alertError(error, errorCodeMessages);
		throw error;
	});
}

export function alertError(error, errorCodeMessages = {}) {
	console.error("AJAX error", error);

	let message = null;
	if (error) {
		if (error.response && error.response.status) {

			let customResponse = null;
			if (
				error.response.data &&
				error.response.data.errorCode &&
				errorCodeMessages[error.response.data.errorCode]
			) {
				customResponse = errorCodeMessages[error.response.data.errorCode];
			} else if (errorCodeMessages[error.response.status]) {
				customResponse = errorCodeMessages[error.response.status];
			}

			if (customResponse) {
				if (typeof customResponse === "string") {
					message = customResponse;
				} if (typeof customResponse == "object") {
					if (customResponse.message) {
						message = customResponse.message;
					}
					if (customResponse.callback) {
						customResponse.callback();
					}
				}
			}

			if (!message) {
				switch (error.response.status) {
					case 400: // Bad Request
						message = "Invalid request. Check the inputs you provided.";
						break;
					case 401: // Unauthorized
						message = "You are not authorized to perform this action. (401)";
						break;
					case 403: // Forbidden
						message = "You are not authorized to perform this action. (403)";
						break;
					case 404: // Not Found
						message = "The requested resource was not found. (404)";
						break;
					case 409: // Conflict
						message = "There is a conflict with existing data. (409)";
						break;
					case 422: // Unprocessable Entity
						message = "Unable to process request. (422)";
						break;
					case 429: // Too Many Requests
						message = "Take a break for a minute!";
						break
					case 500: // Internal Server Error
						message = "An internal server error occurred.";
						break;
					case 503: // Service Unavailable
						message = "The server is currently unavailable. It may be down for database upgrades. Please try again.";
						break;
					default:
						message = "Request failed with HTTP error code " + error.response.status + ".";
				}
			}

		} else if (error.message) {
			message = error.message;
		}
	}

	if (!message) {
		message = "An error occurred.";
	}

	ElMessageBox.alert(message, 'Error', {
		confirmButtonText: "Ok",
		type: "error",
	});
}
