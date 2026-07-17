<template>
<div class="user-page">

	<loading-message message="Locating user..."/>

</div>
</template>

<script>
import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	mounted() {
		ajaxGet('/ajax/user-space-id', {
			id: this.$route.params.id,
		}).then(response => {
			if (response && response.spaceId) {
				this.$router.push({
					name: 'space',
					params: {
						spaceId: response.spaceId,
					},
				});
			} else {
				// Go back
				this.$router.back();
			}
		}).catch(() => {
			this.$message.error('Failed to locate user.');
			// Go back
			this.$router.back();
		});
	},
};
</script>
