<template>
<material-icon :icon="icon"/>
</template>

<script>
import {
	NODE_CLASS,
} from '@/const.js';

export default {
	props: {
		node: {
			type: Object,
			required: true,
		},
	},
	computed: {
		icon() {
			if (this.node.id === 1) {
				return 'park'; // TreeTime info category
			}
			switch (this.node.class) {
				case NODE_CLASS.TAG:
					return 'label';
				case NODE_CLASS.LANG:
					return 'translate';
				case NODE_CLASS.TYPE:
					return 'category';
				case NODE_CLASS.FIELD:
					return 'variables';
				case NODE_CLASS.CATEGORY:
					return 'folder';
				case NODE_CLASS.POST:
					return 'description';
				case NODE_CLASS.COMMENT:
					return 'comment';
			}
			return 'network_node';
		},
	},
};
</script>
