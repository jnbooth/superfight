import {
  defineConfigWithVueTs,
  vueTsConfigs,
} from '@vue/eslint-config-typescript';
import pluginPrettier from 'eslint-plugin-prettier/recommended';
import pluginVue from 'eslint-plugin-vue';
import { globalIgnores } from 'eslint/config';

export default defineConfigWithVueTs(
  globalIgnores(['dist/']),
  pluginVue.configs['flat/recommended'],
  vueTsConfigs.recommended,
  pluginPrettier,
);
