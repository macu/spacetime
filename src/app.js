import '@/styles/app.scss';
import '@/styles/layouts.scss';

import {createApp} from "vue";

import App from "./app.vue";
import router from "./router.js";
import store from "./store.js";

import ElementPlus from 'element-plus';

const app = createApp(App);

app.use(router);
app.use(store);

app.use(ElementPlus, {
	locale: window.ElementPlusLocaleEn,
});

app.mount("#app");
