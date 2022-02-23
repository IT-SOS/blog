import { Plugin } from 'vite';
import legacy from '@vitejs/plugin-legacy';
import vue from '@vitejs/plugin-vue'
import importToCDN, {autoComplete} from 'vite-plugin-cdn-import'
import { viteVConsole } from 'vite-plugin-vconsole';
import { VITE_APP_VISUALIZER } from '../index';
import configVisualizerConfig from './visualizer';
import * as path from 'path'
export default function createVitePlugins() {
    const vitePlugins: (Plugin | Plugin[])[] = [
        vue(),
        viteVConsole({
            entry: path.resolve('src/main.ts'),
            localEnabled: true,
            enabled: false,
            config: {
                maxLogNumber: 1000,
                theme: 'dark'
            }
        }),
        legacy({
            // 不考虑ie的兼容性和老的vivo/锤子/荣耀等国产机型，则可以使用下面
            // targets: ['defaults', 'not IE 11', '> 0.25%, not dead'],
            // 如果考虑上面说的兼容性问题，使用下面配置
            targets: ['ie >= 11'],
            additionalLegacyPolyfills: ['regenerator-runtime/runtime'],

            // 下面可以根据情况选用
            renderLegacyChunks: true,
            polyfills: [
                'es.symbol',
                'es.array.filter',
                'es.promise',
                'es.promise.finally',
                'es/map',
                'es/set',
                'es.array.for-each',
                'es.object.define-properties',
                'es.object.define-property',
                'es.object.get-own-property-descriptor',
                'es.object.get-own-property-descriptors',
                'es.object.keys',
                'es.object.to-string',
                'web.dom-collections.for-each',
                'esnext.global-this',
                'esnext.string.match-all'
            ],
            modernPolyfills: ['es.promise.finally']
        }),
        importToCDN({
            modules:[
                //https://unpkg.com/vue@3.2.29/dist/vue.global.js
                autoComplete('vue'),
                autoComplete('axios'),
                autoComplete('crypto-js'),
                // autoComplete('@vueuse/shared'),
                // autoComplete('@vueuse/core'),
                {
                    name: 'element-plus',
                    var: 'ElementPlus',
                    path: `https://unpkg.com/element-plus@1.2.0-beta.6/dist/index.full.min.js`,
                    css: `https://unpkg.com/element-plus@1.2.0-beta.6/dist/index.css`
                },
                {
                    name: 'vue-router',
                    var: 'VueRouter',
                    path: `https://unpkg.com/vue-router@4.0.12`,
                },
                {
                    name: 'vuex',
                    var: 'createStore',
                    path: `https://unpkg.com/vuex@4.0.0/dist/vuex.global.js`,
                },
            ]
        }),
    ];
    VITE_APP_VISUALIZER && vitePlugins.push(configVisualizerConfig());
    return vitePlugins;
}