import '@/styles/app.scss';
import '@/styles/layouts.scss';

import {createApp} from "vue";

import App from "./app.vue";
import router from "./router.js";
import store from "./store.js";

import ElementPlus from 'element-plus';

import MaterialIcon from '@/widgets/material-icon.vue';
import LoadingMessage from '@/widgets/loading-message.vue';

const app = createApp(App);

app.use(router);
app.use(store);

app.use(ElementPlus, {
	locale: window.ElementPlusLocaleEn,
});

app.component('material-icon', MaterialIcon);
app.component('loading-message', LoadingMessage);

app.mount("#app");
