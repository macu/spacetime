import '@/styles/app.scss';
import '@/styles/layouts.scss';

import {createApp} from "vue";

import App from "./app.vue";
import router from "./router.js";
import store from "./store.js";

import ElementPlus from 'element-plus';

import MaterialIcon from '@/widgets/material-icon.vue';
import LoadingMessage from '@/widgets/loading-message.vue';
import HorizontalControls from '@/widgets/horizontal-controls.vue';
import FormLayout from '@/widgets/form-layout.vue';
import FormField from '@/widgets/form-field.vue';
import FormActions from '@/widgets/form-actions.vue';
import Moment from '@/widgets/moment.vue';
import ReturnToTop from '@/widgets/return-to-top.vue';

const app = createApp(App);

app.use(router);
app.use(store);

app.use(ElementPlus, {
	locale: window.ElementPlusLocaleEn,
});

app.component('material-icon', MaterialIcon);
app.component('loading-message', LoadingMessage);
app.component('horizontal-controls', HorizontalControls);
app.component('form-layout', FormLayout);
app.component('form-field', FormField);
app.component('form-actions', FormActions);
app.component('moment', Moment);
app.component('return-to-top', ReturnToTop);

app.mount("#app");
